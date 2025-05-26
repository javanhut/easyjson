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
	switch s := source.(type) {
	case string:
		return ParseSafely(s)
	case []byte:
		return ParseSafely(string(s))
	case io.Reader:
		data, err := io.ReadAll(s)
		if err != nil {
			return &ParseResult{
				Data:        NewObject(),
				Error:       fmt.Errorf("failed to read from source: %v", err),
				Suggestions: []string{"Check if the reader is valid and contains data"},
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
		// Find the end of JSON
		var depth int
		var inString bool
		var escaped bool
		var end int

		for i, char := range trimmed {
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
