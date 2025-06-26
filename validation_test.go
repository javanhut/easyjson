package easyjson

import (
	"strings"
	"testing"
)

// Test path validation
func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errType string
	}{
		{"valid simple path", "user.name", false, ""},
		{"valid array index path", "users.0.name", false, ""},
		{"valid complex path", "data.items.0.details.name", false, ""},
		{"empty path", "", true, "Path"},
		{"too long path", strings.Repeat("a.", 501), true, "Path"},
		{"too deep path", strings.Repeat("a.", 51)[:101], true, "Path"},
		{"consecutive dots", "user..name", true, "Path"},
		{"parent reference", "user../name", true, "Path"},
		{"negative array index", "users.-1.name", true, "Path"},
		{"large array index", "users.1000000.name", true, "Path"},
		{"control characters", "user.\x00.name", true, "Path"},
		{"invalid UTF-8", "user.\xff.name", true, "Path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errType != "" {
				if valErr, ok := err.(*ValidationError); ok {
					if valErr.Type != tt.errType {
						t.Errorf("validatePath() error type = %v, want %v", valErr.Type, tt.errType)
					}
				}
			}
		})
	}
}

// Test array index validation
func TestValidateArrayIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		arrayLen int
		wantErr  bool
	}{
		{"valid index", 5, 10, false},
		{"boundary index", 9, 10, false},
		{"zero index", 0, 10, false},
		{"negative index", -1, 10, true},
		{"out of bounds", 10, 10, true},
		{"very large index", 1000000, -1, true},
		{"empty array valid", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArrayIndex(tt.index, tt.arrayLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArrayIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test key validation
func TestValidateKey(t *testing.T) {
	tests := []struct {
		name    string
		key     interface{}
		wantErr bool
	}{
		{"valid string key", "username", false},
		{"valid int key", 42, false},
		{"valid float64 int", float64(42), false},
		{"empty string key", "", true},
		{"too long string key", strings.Repeat("a", 501), true},
		{"negative int key", -1, true},
		{"float key", 3.14, true},
		{"nil key", nil, true},
		{"bool key", true, true},
		{"slice key", []int{1, 2, 3}, true},
		{"control char key", "user\x00name", true},
		{"invalid UTF-8 key", "\xff\xfe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test recursion depth validation
func TestValidateRecursionDepth(t *testing.T) {
	tests := []struct {
		name    string
		current int
		max     int
		wantErr bool
	}{
		{"valid depth", 5, 10, false},
		{"zero depth", 0, 10, false},
		{"boundary depth", 9, 10, false},
		{"exceeds max", 10, 10, true},
		{"negative current", -1, 10, true},
		{"zero max", 5, 0, true},
		{"negative max", 5, -1, true},
		{"exceeds system limit", 5, 1001, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecursionDepth(tt.current, tt.max)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRecursionDepth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test memory limit validation
func TestValidateMemoryLimit(t *testing.T) {
	tests := []struct {
		name    string
		size    int64
		limit   int64
		wantErr bool
	}{
		{"under limit", 1000, 2000, false},
		{"at limit", 2000, 2000, false},
		{"over limit", 2001, 2000, true},
		{"negative size", -1, 2000, true},
		{"zero limit", 1000, 0, true},
		{"negative limit", 1000, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMemoryLimit(tt.size, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMemoryLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test JSON value type validation
func TestValidateJSONValueType(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"nil value", nil, false},
		{"bool value", true, false},
		{"string value", "hello", false},
		{"int value", 42, false},
		{"float64 value", 3.14, false},
		{"array value", []interface{}{1, 2, 3}, false},
		{"object value", map[string]interface{}{"key": "value"}, false},
		{"JSONValue", &JSONValue{data: "test"}, false},
		{"invalid struct", struct{}{}, true},
		{"invalid func", func() {}, true},
		{"invalid channel", make(chan int), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONValueType(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateJSONValueType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test composite validation for Get operations
func TestValidateForGet(t *testing.T) {
	obj := New(map[string]interface{}{"key": "value"})
	arr := New([]interface{}{1, 2, 3})

	tests := []struct {
		name    string
		jv      *JSONValue
		key     interface{}
		wantErr bool
	}{
		{"valid object access", obj, "key", false},
		{"valid array access", arr, 1, false},
		{"nil JSONValue", nil, "key", true},
		{"wrong key type for object", obj, 123, true},
		{"wrong key type for array", arr, "key", true},
		{"out of bounds array", arr, 10, true},
		{"invalid key", obj, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateForGet(tt.jv, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateForGet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test composite validation for Set operations
func TestValidateForSet(t *testing.T) {
	obj := New(map[string]interface{}{"key": "value"})
	arr := New([]interface{}{1, 2, 3})

	tests := []struct {
		name    string
		jv      *JSONValue
		key     interface{}
		value   interface{}
		wantErr bool
	}{
		{"valid object set", obj, "newkey", "newvalue", false},
		{"valid array set", arr, 1, "newvalue", false},
		{"nil JSONValue", nil, "key", "value", true},
		{"wrong key type for object", obj, 123, "value", true},
		{"wrong key type for array", arr, "key", "value", true},
		{"invalid value type", obj, "key", func() {}, true},
		{"invalid key", obj, nil, "value", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateForSet(tt.jv, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateForSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test validation error structure
func TestValidationError(t *testing.T) {
	err := newValidationError("TestType", "test message", "test value")

	if err.Type != "TestType" {
		t.Errorf("Expected error type 'TestType', got '%s'", err.Type)
	}

	if err.Message != "test message" {
		t.Errorf("Expected error message 'test message', got '%s'", err.Message)
	}

	if err.Value != "test value" {
		t.Errorf("Expected error value 'test value', got '%v'", err.Value)
	}

	expectedErrorString := "TestType validation error: test message (value: test value)"
	if err.Error() != expectedErrorString {
		t.Errorf("Expected error string '%s', got '%s'", expectedErrorString, err.Error())
	}
}

// Test batch size validation
func TestValidateBatchSize(t *testing.T) {
	tests := []struct {
		name      string
		batchSize int
		wantErr   bool
	}{
		{"valid small batch", 10, false},
		{"valid large batch", 5000, false},
		{"zero batch", 0, false},
		{"max batch size", MaxBatchOperations, false},
		{"negative batch", -1, true},
		{"too large batch", MaxBatchOperations + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBatchSize(tt.batchSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateBatchSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test JSON size estimation
func TestEstimateJSONSize(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		minSize  int64
		maxSize  int64
	}{
		{"null", nil, 4, 4},
		{"true", true, 4, 4},
		{"false", false, 5, 5},
		{"string", "hello", 7, 7}, // "hello" = 5 + 2 quotes
		{"number", 42, 10, 30},
		{"empty array", []interface{}{}, 2, 2},
		{"simple array", []interface{}{1, 2, 3}, 10, 70},
		{"empty object", map[string]interface{}{}, 2, 2},
		{"simple object", map[string]interface{}{"key": "value"}, 15, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := estimateJSONSize(tt.data)
			if size < tt.minSize || size > tt.maxSize {
				t.Errorf("estimateJSONSize() = %v, want between %v and %v", size, tt.minSize, tt.maxSize)
			}
		})
	}
}

// Test default validation config
func TestDefaultValidationConfig(t *testing.T) {
	config := DefaultValidationConfig()

	if config.MaxPathLength != MaxPathLength {
		t.Errorf("Expected MaxPathLength %d, got %d", MaxPathLength, config.MaxPathLength)
	}

	if config.MaxPathDepth != MaxPathDepth {
		t.Errorf("Expected MaxPathDepth %d, got %d", MaxPathDepth, config.MaxPathDepth)
	}

	if config.MaxRecursionDepth != DefaultMaxRecursionDepth {
		t.Errorf("Expected MaxRecursionDepth %d, got %d", DefaultMaxRecursionDepth, config.MaxRecursionDepth)
	}
}

// Benchmark validation functions
func BenchmarkValidatePath(b *testing.B) {
	path := "data.items.0.details.name"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validatePath(path)
	}
}

func BenchmarkValidateKey(b *testing.B) {
	key := "username"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateKey(key)
	}
}

func BenchmarkValidateArrayIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateArrayIndex(5, 10)
	}
}

func BenchmarkEstimateJSONSize(b *testing.B) {
	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "item1"},
			map[string]interface{}{"id": 2, "name": "item2"},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		estimateJSONSize(data)
	}
}