package easyjson

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// safe_parsing.go - Safe parsing with helpful error messages

// ParseResult holds parsing results with helpful feedback
type ParseResult struct {
	Data        *JSONValue
	Error       error
	Suggestions []string
}

// ParseSafely never panics, always returns a valid JSONValue
// Usage: result := easyjson.ParseSafely(jsonString)
func ParseSafely(jsonStr string) *ParseResult {
	result := &ParseResult{}

	// Validate input size
	if err := validateStringLength(jsonStr); err != nil {
		result.Error = fmt.Errorf("JSON string too large: %w", err)
		result.Data = NewObject()
		result.Suggestions = []string{"Reduce JSON size or increase memory limits"}
		return result
	}

	// Validate JSON size estimate
	if err := validateJSONSize(int64(len(jsonStr))); err != nil {
		result.Error = fmt.Errorf("JSON exceeds size limits: %w", err)
		result.Data = NewObject()
		result.Suggestions = []string{"Reduce JSON size or increase memory limits"}
		return result
	}

	if data, err := Loads(jsonStr); err == nil {
		result.Data = data
		return result
	} else {
		result.Error = err
		result.Data = NewObject() // Always return valid JSONValue

		// Provide helpful suggestions for common errors
		suggestions := []string{}

		if strings.Contains(err.Error(), "unexpected end") {
			suggestions = append(suggestions, "JSON appears to be truncated - check if the string is complete")
		}

		if strings.Contains(err.Error(), "invalid character") {
			suggestions = append(suggestions, "Check for unescaped quotes or special characters")
			suggestions = append(suggestions, "Verify all strings are properly quoted")
		}

		if strings.Contains(err.Error(), "cannot unmarshal") {
			suggestions = append(suggestions, "Check data types - ensure numbers aren't quoted as strings")
		}

		// Check for common Python-style boolean mistakes
		if strings.Contains(strings.ToLower(jsonStr), "true") ||
			strings.Contains(strings.ToLower(jsonStr), "false") {
			if strings.Contains(jsonStr, "True") || strings.Contains(jsonStr, "False") {
				suggestions = append(suggestions, "Use lowercase 'true'/'false' instead of 'True'/'False'")
			}
		}

		// Check for Python None vs null
		if strings.Contains(jsonStr, "None") {
			suggestions = append(suggestions, "Use 'null' instead of 'None'")
		}

		// Check for single quotes (common mistake)
		if strings.Contains(jsonStr, "'") && !strings.Contains(jsonStr, "\"") {
			suggestions = append(suggestions, "Use double quotes (\") instead of single quotes (')")
		}

		result.Suggestions = suggestions
		return result
	}
}

// ParseSafelyFrom parses JSON from various sources with safety
// Usage: result := easyjson.ParseSafelyFrom(reader)
func ParseSafelyFrom(source interface{}) *ParseResult {
	// Validate source
	if source == nil {
		return &ParseResult{
			Data:        NewObject(),
			Error:       fmt.Errorf("source cannot be nil"),
			Suggestions: []string{"Provide a valid source (string, []byte, or io.Reader)"},
		}
	}

	switch s := source.(type) {
	case string:
		return ParseSafely(s)
	case []byte:
		// Validate byte slice size
		if err := validateJSONSize(int64(len(s))); err != nil {
			return &ParseResult{
				Data:        NewObject(),
				Error:       fmt.Errorf("byte slice too large: %w", err),
				Suggestions: []string{"Reduce data size or increase memory limits"},
			}
		}
		return ParseSafely(string(s))
	case io.Reader:
		// Limit reading to prevent memory exhaustion
		limitedReader := io.LimitReader(s, MaxJSONSize)
		data, err := io.ReadAll(limitedReader)
		if err != nil {
			return &ParseResult{
				Data:        NewObject(),
				Error:       fmt.Errorf("failed to read from source: %v", err),
				Suggestions: []string{"Check if the reader is valid and contains data"},
			}
		}
		// Check if we hit the size limit
		if int64(len(data)) >= MaxJSONSize {
			return &ParseResult{
				Data:        NewObject(),
				Error:       fmt.Errorf("data from reader exceeds maximum size limit of %d bytes", MaxJSONSize),
				Suggestions: []string{"Reduce data size or increase memory limits"},
			}
		}
		return ParseSafely(string(data))
	default:
		return &ParseResult{
			Data:        NewObject(),
			Error:       fmt.Errorf("unsupported source type: %T", source),
			Suggestions: []string{"Use string, []byte, or io.Reader as source"},
		}
	}
}

// MustParse panics in development, returns empty object in production
// Usage: data := easyjson.MustParse(jsonString)
func MustParse(jsonStr string) *JSONValue {
	if data, err := Loads(jsonStr); err == nil {
		return data
	} else {
		if isDevelopment() {
			panic(fmt.Sprintf("JSON parsing failed: %v\nJSON: %s", err, jsonStr))
		}
		return NewObject()
	}
}

// MustParseFrom is like MustParse but accepts various sources
func MustParseFrom(source interface{}) *JSONValue {
	result := ParseSafelyFrom(source)
	if result.Error != nil {
		if isDevelopment() {
			panic(fmt.Sprintf("JSON parsing failed: %v", result.Error))
		}
	}
	return result.Data
}

// TryParse attempts to parse, returns success boolean and data
// Usage: if data, ok := easyjson.TryParse(jsonString); ok { ... }
func TryParse(jsonStr string) (*JSONValue, bool) {
	if data, err := Loads(jsonStr); err == nil {
		return data, true
	}
	return NewObject(), false
}

// ParseOrDefault parses JSON or returns default on error
// Usage: data := easyjson.ParseOrDefault(jsonString, easyjson.NewObject())
func ParseOrDefault(jsonStr string, defaultValue *JSONValue) *JSONValue {
	if data, err := Loads(jsonStr); err == nil {
		return data
	}
	return defaultValue
}

// ValidateJSON checks if string is valid JSON without parsing
// Usage: if easyjson.ValidateJSON(jsonString) { ... }
func ValidateJSON(jsonStr string) bool {
	_, err := Loads(jsonStr)
	return err == nil
}

// ValidateJSONWithDetails provides detailed validation info
func ValidateJSONWithDetails(jsonStr string) (bool, error, []string) {
	result := ParseSafely(jsonStr)
	return result.Error == nil, result.Error, result.Suggestions
}

// FixCommonIssues attempts to fix common JSON formatting issues
// Usage: fixed := easyjson.FixCommonIssues(brokenJSON)
func FixCommonIssues(jsonStr string) string {
	// Validate input size
	if err := validateStringLength(jsonStr); err != nil {
		return "{}" // Return empty object for oversized input
	}

	fixed := jsonStr

	// Fix Python-style booleans
	fixed = strings.ReplaceAll(fixed, "True", "true")
	fixed = strings.ReplaceAll(fixed, "False", "false")

	// Fix Python None
	fixed = strings.ReplaceAll(fixed, "None", "null")

	// Attempt to fix single quotes (simple case)
	if !strings.Contains(fixed, "\"") && strings.Contains(fixed, "'") {
		// Only if no double quotes exist, replace single quotes
		fixed = strings.ReplaceAll(fixed, "'", "\"")
	}

	return fixed
}

// ParseWithFixes attempts to parse after applying common fixes
// Usage: data := easyjson.ParseWithFixes(messyJSONString)
func ParseWithFixes(jsonStr string) (*JSONValue, error) {
	// Try original first
	if data, err := Loads(jsonStr); err == nil {
		return data, nil
	}

	// Try with fixes
	fixed := FixCommonIssues(jsonStr)
	return Loads(fixed)
}

// isDevelopment checks if we're in development mode
func isDevelopment() bool {
	env := strings.ToLower(os.Getenv("GO_ENV"))
	return env == "" || env == "development" || env == "dev"
}

// ParseLenient is very forgiving - tries multiple strategies to parse JSON
// Usage: data := easyjson.ParseLenient(messyJSONString)
func ParseLenient(jsonStr string) *JSONValue {
	// Validate input size first
	if err := validateStringLength(jsonStr); err != nil {
		return NewObject() // Return empty object for oversized input
	}

	// Validate JSON size estimate
	if err := validateJSONSize(int64(len(jsonStr))); err != nil {
		return NewObject() // Return empty object if too large
	}

	// Strategy 1: Try as-is
	if data, err := Loads(jsonStr); err == nil {
		return data
	}

	// Strategy 2: Try with common fixes
	if data, err := ParseWithFixes(jsonStr); err == nil {
		return data
	}

	// Strategy 3: Try to extract JSON from a larger string
	trimmed := strings.TrimSpace(jsonStr)
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		// Validate trimmed string size
		if err := validateStringLength(trimmed); err != nil {
			return NewObject()
		}

		// Find the end of JSON with depth and length limits
		var depth int
		var inString bool
		var escaped bool
		var end int
		maxProcessed := 0

		for i, char := range trimmed {
			// Prevent infinite processing
			maxProcessed++
			if maxProcessed > MaxStringLength {
				break
			}
			switch char {
			case '\\':
				escaped = !escaped
				continue
			case '"':
				if !escaped {
					inString = !inString
				}
			case '{', '[':
				if !inString {
					depth++
					// Prevent excessive nesting
					if depth > DefaultMaxRecursionDepth {
						return NewObject()
					}
				}
			case '}', ']':
				if !inString {
					depth--
					if depth == 0 {
						end = i + 1
						break
					}
				}
			}
			escaped = false
		}

		if end > 0 {
			extracted := trimmed[:end]
			if data, err := Loads(extracted); err == nil {
				return data
			}
		}
	}

	// Strategy 4: Return empty object as fallback
	return NewObject()
}
