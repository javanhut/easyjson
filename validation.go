package easyjson

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// validation.go - Comprehensive input validation functions for EasyJSON library

// ValidationError represents a validation error with context
type ValidationError struct {
	Type    string
	Message string
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s validation error: %s (value: %v)", e.Type, e.Message, e.Value)
}

// newValidationError creates a new ValidationError
func newValidationError(errorType, message string, value interface{}) *ValidationError {
	return &ValidationError{
		Type:    errorType,
		Message: message,
		Value:   value,
	}
}

// Configuration constants for validation limits
const (
	// Path validation limits
	MaxPathLength       = 1000
	MaxPathDepth        = 50
	MaxPathSegmentLength = 200
	
	// Array validation limits
	MaxArraySize        = 1000000
	MaxArrayIndex       = 999999
	
	// Key validation limits
	MaxKeyLength        = 500
	MaxObjectKeys       = 10000
	
	// Recursion limits
	DefaultMaxRecursionDepth = 100
	MaxRecursionDepth        = 1000
	
	// Memory limits (in bytes)
	MaxStringLength          = 10 * 1024 * 1024  // 10MB
	MaxJSONSize             = 100 * 1024 * 1024  // 100MB
	
	// Operation limits
	MaxBatchOperations      = 10000
)

// Path validation

// validatePath validates a dot-separated path string
func validatePath(path string) error {
	if path == "" {
		return newValidationError("Path", "path cannot be empty", path)
	}
	
	if len(path) > MaxPathLength {
		return newValidationError("Path", fmt.Sprintf("path exceeds maximum length of %d characters", MaxPathLength), len(path))
	}
	
	// Check for invalid characters
	if strings.Contains(path, "..") {
		return newValidationError("Path", "path cannot contain '..' (parent directory references)", path)
	}
	
	// Split and validate each segment
	segments := strings.Split(path, ".")
	if len(segments) > MaxPathDepth {
		return newValidationError("Path", fmt.Sprintf("path exceeds maximum depth of %d levels", MaxPathDepth), len(segments))
	}
	
	for i, segment := range segments {
		if segment == "" && i != 0 && i != len(segments)-1 {
			return newValidationError("Path", "path cannot contain empty segments (consecutive dots)", path)
		}
		
		if len(segment) > MaxPathSegmentLength {
			return newValidationError("Path", fmt.Sprintf("path segment exceeds maximum length of %d characters", MaxPathSegmentLength), segment)
		}
		
		// Check for control characters and invalid Unicode
		if !utf8.ValidString(segment) {
			return newValidationError("Path", "path segment contains invalid UTF-8", segment)
		}
		
		for _, r := range segment {
			if unicode.IsControl(r) && r != '\t' {
				return newValidationError("Path", "path segment contains control characters", segment)
			}
		}
		
		// Validate array index if numeric
		if _, err := strconv.Atoi(segment); err == nil {
			if err := validatePathArrayIndex(segment); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// validatePathArrayIndex validates an array index within a path
func validatePathArrayIndex(indexStr string) error {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return newValidationError("Path", "invalid array index format", indexStr)
	}
	
	if index < 0 {
		return newValidationError("Path", "negative array indices not allowed in paths", index)
	}
	
	if index > MaxArrayIndex {
		return newValidationError("Path", fmt.Sprintf("array index exceeds maximum value of %d", MaxArrayIndex), index)
	}
	
	return nil
}

// validatePathFormat validates path format with regex
func validatePathFormat(path string) error {
	// Allow alphanumeric, dots, underscores, hyphens, and array indices
	validPathRegex := regexp.MustCompile(`^[a-zA-Z0-9._\-\[\]]+$`)
	if !validPathRegex.MatchString(path) {
		return newValidationError("Path", "path contains invalid characters (only letters, numbers, dots, underscores, hyphens, and brackets allowed)", path)
	}
	
	return nil
}

// Array validation

// validateArrayIndex validates an array index against array length
func validateArrayIndex(index int, arrayLen int) error {
	if index < 0 {
		return newValidationError("ArrayIndex", "negative array indices not allowed", index)
	}
	
	if arrayLen >= 0 && index >= arrayLen {
		return newValidationError("ArrayIndex", fmt.Sprintf("index %d out of bounds for array of length %d", index, arrayLen), index)
	}
	
	if index > MaxArrayIndex {
		return newValidationError("ArrayIndex", fmt.Sprintf("index exceeds maximum allowed value of %d", MaxArrayIndex), index)
	}
	
	return nil
}

// validateArraySize validates array size constraints
func validateArraySize(size int) error {
	if size < 0 {
		return newValidationError("ArraySize", "array size cannot be negative", size)
	}
	
	if size > MaxArraySize {
		return newValidationError("ArraySize", fmt.Sprintf("array size exceeds maximum allowed size of %d", MaxArraySize), size)
	}
	
	return nil
}

// validateArrayInsertIndex validates index for array insertion
func validateArrayInsertIndex(index int, arrayLen int) error {
	if index < 0 {
		return newValidationError("ArrayInsert", "negative insertion index not allowed", index)
	}
	
	// For insertion, index can be equal to array length (append operation)
	if arrayLen >= 0 && index > arrayLen {
		return newValidationError("ArrayInsert", fmt.Sprintf("insertion index %d exceeds array length %d", index, arrayLen), index)
	}
	
	if index > MaxArrayIndex {
		return newValidationError("ArrayInsert", fmt.Sprintf("insertion index exceeds maximum allowed value of %d", MaxArrayIndex), index)
	}
	
	return nil
}

// validateArrayIndexForInsertion is an alias for validateArrayInsertIndex for compatibility
func validateArrayIndexForInsertion(index int) error {
	if index < 0 {
		return newValidationError("ArrayInsert", "negative insertion index not allowed", index)
	}
	
	if index > MaxArrayIndex {
		return newValidationError("ArrayInsert", fmt.Sprintf("insertion index exceeds maximum allowed value of %d", MaxArrayIndex), index)
	}
	
	return nil
}

// Key validation

// validateKey validates a key for object operations
func validateKey(key interface{}) error {
	switch k := key.(type) {
	case string:
		return validateStringKey(k)
	case int:
		return validateIntKey(k)
	case float64:
		// JSON numbers are float64, but for keys we only allow integers
		if k != float64(int(k)) {
			return newValidationError("Key", "floating-point numbers not allowed as keys", k)
		}
		return validateIntKey(int(k))
	default:
		return newValidationError("Key", fmt.Sprintf("invalid key type: %T (only string and int allowed)", key), key)
	}
}

// validateStringKey validates a string key
func validateStringKey(key string) error {
	if key == "" {
		return newValidationError("StringKey", "empty string keys not allowed", key)
	}
	
	if len(key) > MaxKeyLength {
		return newValidationError("StringKey", fmt.Sprintf("key exceeds maximum length of %d characters", MaxKeyLength), len(key))
	}
	
	// Check for valid UTF-8
	if !utf8.ValidString(key) {
		return newValidationError("StringKey", "key contains invalid UTF-8", key)
	}
	
	// Check for control characters (except tab)
	for _, r := range key {
		if unicode.IsControl(r) && r != '\t' {
			return newValidationError("StringKey", "key contains control characters", key)
		}
	}
	
	return nil
}

// validateIntKey validates an integer key (for arrays)
func validateIntKey(key int) error {
	if key < 0 {
		return newValidationError("IntKey", "negative integer keys not allowed", key)
	}
	
	if key > MaxArrayIndex {
		return newValidationError("IntKey", fmt.Sprintf("integer key exceeds maximum value of %d", MaxArrayIndex), key)
	}
	
	return nil
}

// validateObjectKeyCount validates the number of keys in an object
func validateObjectKeyCount(keyCount int) error {
	if keyCount < 0 {
		return newValidationError("ObjectKeys", "key count cannot be negative", keyCount)
	}
	
	if keyCount > MaxObjectKeys {
		return newValidationError("ObjectKeys", fmt.Sprintf("object exceeds maximum key count of %d", MaxObjectKeys), keyCount)
	}
	
	return nil
}

// Recursion validation

// validateRecursionDepth validates current recursion depth against maximum
func validateRecursionDepth(current, max int) error {
	if current < 0 {
		return newValidationError("Recursion", "recursion depth cannot be negative", current)
	}
	
	if max <= 0 {
		return newValidationError("Recursion", "maximum recursion depth must be positive", max)
	}
	
	if max > MaxRecursionDepth {
		return newValidationError("Recursion", fmt.Sprintf("maximum recursion depth exceeds system limit of %d", MaxRecursionDepth), max)
	}
	
	if current >= max {
		return newValidationError("Recursion", fmt.Sprintf("recursion depth %d exceeds maximum allowed depth of %d", current, max), current)
	}
	
	return nil
}

// Memory validation

// validateMemoryLimit validates memory usage constraints
func validateMemoryLimit(size int64, limit int64) error {
	if size < 0 {
		return newValidationError("Memory", "memory size cannot be negative", size)
	}
	
	if limit <= 0 {
		return newValidationError("Memory", "memory limit must be positive", limit)
	}
	
	if size > limit {
		return newValidationError("Memory", fmt.Sprintf("memory usage %d bytes exceeds limit of %d bytes", size, limit), size)
	}
	
	return nil
}

// validateStringLength validates string length constraints
func validateStringLength(str string) error {
	if len(str) > MaxStringLength {
		return newValidationError("StringLength", fmt.Sprintf("string length %d exceeds maximum of %d bytes", len(str), MaxStringLength), len(str))
	}
	
	return nil
}

// validateJSONSize validates total JSON data size
func validateJSONSize(size int64) error {
	if size > MaxJSONSize {
		return newValidationError("JSONSize", fmt.Sprintf("JSON size %d bytes exceeds maximum of %d bytes", size, MaxJSONSize), size)
	}
	
	return nil
}

// Batch operation validation

// validateBatchSize validates the size of batch operations
func validateBatchSize(batchSize int) error {
	if batchSize < 0 {
		return newValidationError("BatchSize", "batch size cannot be negative", batchSize)
	}
	
	if batchSize > MaxBatchOperations {
		return newValidationError("BatchSize", fmt.Sprintf("batch size %d exceeds maximum of %d operations", batchSize, MaxBatchOperations), batchSize)
	}
	
	return nil
}

// Type validation

// validateJSONValueType validates that a value is a supported JSON type
func validateJSONValueType(value interface{}) error {
	switch value.(type) {
	case nil, bool, string, float64, int, int64, float32:
		return nil
	case []interface{}, map[string]interface{}:
		return nil
	case *JSONValue:
		return nil
	default:
		return newValidationError("JSONType", fmt.Sprintf("unsupported type: %T (must be a valid JSON type)", value), value)
	}
}

// Composite validation functions

// validateForGet validates parameters for Get operations
func validateForGet(jv *JSONValue, key interface{}) error {
	if jv == nil {
		return newValidationError("Get", "JSONValue cannot be nil", nil)
	}
	
	if err := validateKey(key); err != nil {
		return err
	}
	
	// Additional validation based on JSONValue type
	switch jv.data.(type) {
	case map[string]interface{}:
		if _, ok := key.(string); !ok {
			return newValidationError("Get", "key must be string for object access", key)
		}
	case []interface{}:
		if _, ok := key.(int); !ok {
			return newValidationError("Get", "key must be int for array access", key)
		}
		if arr := jv.data.([]interface{}); len(arr) > 0 {
			if idx, ok := key.(int); ok {
				return validateArrayIndex(idx, len(arr))
			}
		}
	default:
		return newValidationError("Get", "cannot access key on non-object/array type", jv.data)
	}
	
	return nil
}

// validateForSet validates parameters for Set operations
func validateForSet(jv *JSONValue, key interface{}, value interface{}) error {
	if jv == nil {
		return newValidationError("Set", "JSONValue cannot be nil", nil)
	}
	
	if err := validateKey(key); err != nil {
		return err
	}
	
	if err := validateJSONValueType(value); err != nil {
		return err
	}
	
	// Additional validation based on JSONValue type
	switch v := jv.data.(type) {
	case map[string]interface{}:
		if _, ok := key.(string); !ok {
			return newValidationError("Set", "key must be string for object", key)
		}
		if err := validateObjectKeyCount(len(v) + 1); err != nil {
			return err
		}
	case []interface{}:
		if _, ok := key.(int); !ok {
			return newValidationError("Set", "key must be int for array", key)
		}
		if idx, ok := key.(int); ok {
			// For set operations, we might expand the array
			newSize := idx + 1
			if len(v) < newSize {
				return validateArraySize(newSize)
			}
		}
	default:
		return newValidationError("Set", "cannot set on non-object/array type", jv.data)
	}
	
	return nil
}

// validateForPath validates parameters for Path operations
func validateForPath(path string) error {
	if err := validatePath(path); err != nil {
		return err
	}
	
	// Additional path-specific validations can be added here
	return nil
}

// Utility functions

// isValidUTF8 checks if a string contains valid UTF-8
func isValidUTF8(s string) bool {
	return utf8.ValidString(s)
}

// containsControlChars checks if a string contains control characters
func containsControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}
	return false
}

// estimateJSONSize estimates the memory size of a JSON structure
func estimateJSONSize(data interface{}) int64 {
	switch v := data.(type) {
	case nil:
		return 4 // "null"
	case bool:
		if v {
			return 4 // "true"
		}
		return 5 // "false"
	case string:
		return int64(len(v)) + 2 // quotes
	case float64, int, int64, float32:
		return 20 // rough estimate for numbers
	case []interface{}:
		size := int64(2) // brackets
		for i, item := range v {
			if i > 0 {
				size += 1 // comma
			}
			size += estimateJSONSize(item)
		}
		return size
	case map[string]interface{}:
		size := int64(2) // braces
		first := true
		for key, value := range v {
			if !first {
				size += 1 // comma
			}
			first = false
			size += int64(len(key)) + 3 // key + quotes + colon
			size += estimateJSONSize(value)
		}
		return size
	default:
		return 50 // fallback estimate
	}
}

// ValidationConfig holds configuration for validation limits
type ValidationConfig struct {
	MaxPathLength       int
	MaxPathDepth        int
	MaxArraySize        int
	MaxKeyLength        int
	MaxRecursionDepth   int
	MaxStringLength     int
	MaxJSONSize         int64
	MaxBatchOperations  int
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		MaxPathLength:      MaxPathLength,
		MaxPathDepth:       MaxPathDepth,
		MaxArraySize:       MaxArraySize,
		MaxKeyLength:       MaxKeyLength,
		MaxRecursionDepth:  DefaultMaxRecursionDepth,
		MaxStringLength:    MaxStringLength,
		MaxJSONSize:        MaxJSONSize,
		MaxBatchOperations: MaxBatchOperations,
	}
}

// Custom validation with configuration

// validateWithConfig validates using custom configuration
func validateWithConfig(config *ValidationConfig) func(interface{}) error {
	return func(value interface{}) error {
		// This can be extended to use custom config values
		return validateJSONValueType(value)
	}
}