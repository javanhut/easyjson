package easyjson

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// JSONValue represents a flexible JSON value that can be any type
type JSONValue struct {
	data interface{}
}

// Q provides a fluent query interface for chaining access
// Usage: data.Q("name", 0, "hair_color").String()
func (jv *JSONValue) Q(keys ...interface{}) *JSONValue {
	current := jv
	for _, key := range keys {
		current = current.Get(key)
		if current.IsNull() {
			break
		}
	}
	return current
}

// New creates a new JSONValue from any Go value
func New(data interface{}) *JSONValue {
	return &JSONValue{data: data}
}

// Loads parses a JSON string and returns a JSONValue
func Loads(jsonStr string) (*JSONValue, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return &JSONValue{data: data}, nil
}

// Load parses JSON from a byte slice and returns a JSONValue
func Load(jsonBytes []byte) (*JSONValue, error) {
	var data interface{}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return nil, err
	}
	return &JSONValue{data: data}, nil
}

// LoadFile reads and parses JSON from a file
func LoadFile(filename string) (*JSONValue, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return Load(data)
}

// Dumps converts the JSONValue to a JSON string
func (jv *JSONValue) Dumps() (string, error) {
	bytes, err := json.Marshal(jv.data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// DumpsIndent converts the JSONValue to a pretty-printed JSON string
func (jv *JSONValue) DumpsIndent(indent string) (string, error) {
	bytes, err := json.MarshalIndent(jv.data, "", indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Dump converts the JSONValue to JSON bytes
func (jv *JSONValue) Dump() ([]byte, error) {
	return json.Marshal(jv.data)
}

// SaveFile writes the JSONValue to a file as JSON
func (jv *JSONValue) SaveFile(filename string) error {
	data, err := jv.Dump()
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

// SaveFileIndent writes the JSONValue to a file as pretty-printed JSON
func (jv *JSONValue) SaveFileIndent(filename string, indent string) error {
	data, err := json.MarshalIndent(jv.data, "", indent)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}

// Get retrieves a value by key (for objects) or index (for arrays)
func (jv *JSONValue) Get(key interface{}) *JSONValue {
	// Validate key
	if err := validateKey(key); err != nil {
		return &JSONValue{data: nil}
	}

	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			if val, exists := v[keyStr]; exists {
				return &JSONValue{data: val}
			}
		}
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			// Validate array index
			if err := validateArrayIndex(keyInt, len(v)); err != nil {
				return &JSONValue{data: nil}
			}
			return &JSONValue{data: v[keyInt]}
		}
	}
	return &JSONValue{data: nil}
}

// Set sets a value by key (for objects) or index (for arrays)
func (jv *JSONValue) Set(key interface{}, value interface{}) error {
	// Validate inputs
	if err := validateKey(key); err != nil {
		return fmt.Errorf("invalid key: %w", err)
	}
	
	if err := validateJSONValueType(value); err != nil {
		return fmt.Errorf("invalid value type: %w", err)
	}

	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			v[keyStr] = value
			return nil
		}
		return fmt.Errorf("key must be string for object")
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt < 0 {
				return fmt.Errorf("negative array index not allowed")
			}
			// Validate array index for insertion
			if err := validateArrayIndexForInsertion(keyInt); err != nil {
				return fmt.Errorf("invalid array index: %w", err)
			}
			// Resize array if necessary
			if keyInt >= len(v) {
				newSlice := make([]interface{}, keyInt+1)
				copy(newSlice, v)
				jv.data = newSlice
				v = newSlice
			}
			v[keyInt] = value
			return nil
		}
		return fmt.Errorf("key must be int for array")
	default:
		return fmt.Errorf("cannot set on non-object/array type")
	}
}

// Has checks if a key exists (for objects) or index is valid (for arrays)
func (jv *JSONValue) Has(key interface{}) bool {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			_, exists := v[keyStr]
			return exists
		}
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			return keyInt >= 0 && keyInt < len(v)
		}
	}
	return false
}

// Delete removes a key from an object or index from array
func (jv *JSONValue) Delete(key interface{}) error {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			delete(v, keyStr)
			return nil
		}
		return fmt.Errorf("key must be string for object")
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt < 0 || keyInt >= len(v) {
				return fmt.Errorf("index out of range: %d", keyInt)
			}
			// Create new slice to avoid slice behavior issues
			newSlice := make([]interface{}, 0, len(v)-1)
			// Copy elements before the deleted index
			newSlice = append(newSlice, v[:keyInt]...)
			// Copy elements after the deleted index
			if keyInt < len(v)-1 {
				newSlice = append(newSlice, v[keyInt+1:]...)
			}
			jv.data = newSlice
			return nil
		}
		return fmt.Errorf("key must be int for array")
	default:
		return fmt.Errorf("cannot delete from non-object/array type")
	}
}

// Keys returns all keys for an object
func (jv *JSONValue) Keys() []string {
	if obj, ok := jv.data.(map[string]interface{}); ok {
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		return keys
	}
	return []string{}
}

// Values returns all values for an object or array
func (jv *JSONValue) Values() []*JSONValue {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		values := make([]*JSONValue, 0, len(v))
		for _, val := range v {
			values = append(values, &JSONValue{data: val})
		}
		return values
	case []interface{}:
		values := make([]*JSONValue, len(v))
		for i, val := range v {
			values[i] = &JSONValue{data: val}
		}
		return values
	}
	return []*JSONValue{}
}

// Items returns key-value pairs for an object
func (jv *JSONValue) Items() map[string]*JSONValue {
	if obj, ok := jv.data.(map[string]interface{}); ok {
		items := make(map[string]*JSONValue)
		for k, v := range obj {
			items[k] = &JSONValue{data: v}
		}
		return items
	}
	return map[string]*JSONValue{}
}

// Len returns the length of an array or object
func (jv *JSONValue) Len() int {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		return len(v)
	case []interface{}:
		return len(v)
	case string:
		return len(v)
	}
	return 0
}

// IsNull checks if the value is null
func (jv *JSONValue) IsNull() bool {
	return jv.data == nil
}

// IsObject checks if the value is an object
func (jv *JSONValue) IsObject() bool {
	_, ok := jv.data.(map[string]interface{})
	return ok
}

// IsArray checks if the value is an array
func (jv *JSONValue) IsArray() bool {
	_, ok := jv.data.([]interface{})
	return ok
}

// IsString checks if the value is a string
func (jv *JSONValue) IsString() bool {
	_, ok := jv.data.(string)
	return ok
}

// IsNumber checks if the value is a number
func (jv *JSONValue) IsNumber() bool {
	switch jv.data.(type) {
	case float64, int, int64, float32:
		return true
	}
	return false
}

// IsBool checks if the value is a boolean
func (jv *JSONValue) IsBool() bool {
	_, ok := jv.data.(bool)
	return ok
}

// AsString returns the value as a string
func (jv *JSONValue) AsString() string {
	if str, ok := jv.data.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", jv.data)
}

// AsInt returns the value as an integer
func (jv *JSONValue) AsInt() int {
	switch v := jv.data.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}

// AsFloat returns the value as a float64
func (jv *JSONValue) AsFloat() float64 {
	switch v := jv.data.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0.0
}

// AsBool returns the value as a boolean
func (jv *JSONValue) AsBool() bool {
	switch v := jv.data.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true"
	case float64:
		return v != 0
	case int:
		return v != 0
	}
	return false
}

// AsArray returns the value as a slice of JSONValues
func (jv *JSONValue) AsArray() []*JSONValue {
	if arr, ok := jv.data.([]interface{}); ok {
		result := make([]*JSONValue, len(arr))
		for i, v := range arr {
			result[i] = &JSONValue{data: v}
		}
		return result
	}
	return []*JSONValue{}
}

// AsObject returns the value as a map of JSONValues
func (jv *JSONValue) AsObject() map[string]*JSONValue {
	if obj, ok := jv.data.(map[string]interface{}); ok {
		result := make(map[string]*JSONValue)
		for k, v := range obj {
			result[k] = &JSONValue{data: v}
		}
		return result
	}
	return map[string]*JSONValue{}
}

// Raw returns the underlying Go value
func (jv *JSONValue) Raw() interface{} {
	return jv.data
}

// String returns a string representation of the JSONValue
func (jv *JSONValue) String() string {
	if str, err := jv.Dumps(); err == nil {
		return str
	}
	return fmt.Sprintf("%v", jv.data)
}

// Append adds a value to an array
func (jv *JSONValue) Append(value interface{}) error {
	if arr, ok := jv.data.([]interface{}); ok {
		jv.data = append(arr, value)
		return nil
	}
	return fmt.Errorf("cannot append to non-array type")
}

// Extend adds multiple values to an array
func (jv *JSONValue) Extend(values []interface{}) error {
	if arr, ok := jv.data.([]interface{}); ok {
		jv.data = append(arr, values...)
		return nil
	}
	return fmt.Errorf("cannot extend non-array type")
}

// Update merges another object into this one
func (jv *JSONValue) Update(other *JSONValue) error {
	if obj, ok := jv.data.(map[string]interface{}); ok {
		if otherObj, ok := other.data.(map[string]interface{}); ok {
			for k, v := range otherObj {
				obj[k] = v
			}
			return nil
		}
		return fmt.Errorf("can only update with another object")
	}
	return fmt.Errorf("cannot update non-object type")
}

// Clone creates a deep copy of the JSONValue
func (jv *JSONValue) Clone() *JSONValue {
	return &JSONValue{data: jv.deepClone(jv.data)}
}

// deepClone performs a deep copy of the data without JSON marshal/unmarshal
func (jv *JSONValue) deepClone(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		cloned := make(map[string]interface{}, len(v))
		for key, value := range v {
			cloned[key] = jv.deepClone(value)
		}
		return cloned
		
	case []interface{}:
		cloned := make([]interface{}, len(v))
		for i, value := range v {
			cloned[i] = jv.deepClone(value)
		}
		return cloned
		
	case map[interface{}]interface{}:
		// Handle generic maps
		cloned := make(map[interface{}]interface{}, len(v))
		for key, value := range v {
			cloned[key] = jv.deepClone(value)
		}
		return cloned
		
	case []string:
		// Handle string slices
		cloned := make([]string, len(v))
		copy(cloned, v)
		return cloned
		
	case []int:
		// Handle int slices
		cloned := make([]int, len(v))
		copy(cloned, v)
		return cloned
		
	case []float64:
		// Handle float64 slices
		cloned := make([]float64, len(v))
		copy(cloned, v)
		return cloned
		
	case []bool:
		// Handle bool slices
		cloned := make([]bool, len(v))
		copy(cloned, v)
		return cloned
		
	case string, int, int8, int16, int32, int64,
		 uint, uint8, uint16, uint32, uint64,
		 float32, float64, bool:
		// Primitive types are copied by value
		return v
		
	default:
		// For other types, try to use JSON as fallback
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil
		}
		
		var cloned interface{}
		if err := json.Unmarshal(bytes, &cloned); err != nil {
			return nil
		}
		
		return cloned
	}
}

// Path retrieves a nested value using a dot-separated path
func (jv *JSONValue) Path(path string) *JSONValue {
	// Validate path
	if err := validatePath(path); err != nil {
		return &JSONValue{data: nil}
	}

	parts := strings.Split(path, ".")
	current := jv

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Try as array index first
		if index, err := strconv.Atoi(part); err == nil {
			// Validate array index if we're dealing with an array
			if current.IsArray() {
				if err := validateArrayIndex(index, current.Len()); err != nil {
					return &JSONValue{data: nil}
				}
			}
			current = current.Get(index)
		} else {
			// Validate key
			if err := validateKey(part); err != nil {
				return &JSONValue{data: nil}
			}
			current = current.Get(part)
		}

		if current.IsNull() {
			break
		}
	}

	return current
}

// SetPath sets a nested value using a dot-separated path
func (jv *JSONValue) SetPath(path string, value interface{}) error {
	// Validate path
	if err := validatePath(path); err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Validate value type
	if err := validateJSONValueType(value); err != nil {
		return fmt.Errorf("invalid value type: %w", err)
	}

	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return fmt.Errorf("empty path")
	}

	current := jv
	for i, part := range parts[:len(parts)-1] {
		if part == "" {
			continue
		}

		var next *JSONValue
		if index, err := strconv.Atoi(part); err == nil {
			next = current.Get(index)
		} else {
			next = current.Get(part)
		}

		if next.IsNull() {
			// Create intermediate objects/arrays as needed
			if i+1 < len(parts)-1 {
				if nextIndex, err := strconv.Atoi(parts[i+1]); err == nil {
					// Next part is an array index - create array with proper size
					arraySize := nextIndex + 1
					newArray := make([]interface{}, arraySize)
					current.Set(part, newArray)
				} else {
					// Next part is an object key
					newObj := make(map[string]interface{})
					current.Set(part, newObj)
				}
			} else {
				// Last intermediate part - determine type based on final part
				finalPart := parts[len(parts)-1]
				if _, err := strconv.Atoi(finalPart); err == nil {
					// Final part is array index
					finalIndex, _ := strconv.Atoi(finalPart)
					arraySize := finalIndex + 1
					newArray := make([]interface{}, arraySize)
					current.Set(part, newArray)
				} else {
					// Final part is object key
					newObj := make(map[string]interface{})
					current.Set(part, newObj)
				}
			}

			if index, err := strconv.Atoi(part); err == nil {
				next = current.Get(index)
			} else {
				next = current.Get(part)
			}
		}

		current = next
	}

	lastPart := parts[len(parts)-1]
	if index, err := strconv.Atoi(lastPart); err == nil {
		return current.Set(index, value)
	} else {
		return current.Set(lastPart, value)
	}
}

// NewObject creates a new JSONValue representing an empty object
func NewObject() *JSONValue {
	return &JSONValue{data: make(map[string]interface{})}
}

// NewArray creates a new JSONValue representing an empty array
func NewArray() *JSONValue {
	return &JSONValue{data: make([]interface{}, 0)}
}

// NewArrayFrom creates a new JSONValue array from a slice
func NewArrayFrom(items []interface{}) *JSONValue {
	return &JSONValue{data: items}
}

// NewObjectFrom creates a new JSONValue object from a map
func NewObjectFrom(obj map[string]interface{}) *JSONValue {
	return &JSONValue{data: obj}
}
