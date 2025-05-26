package easyjson

import (
	"encoding/json"
	"fmt"
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

// Get retrieves a value by key (for objects) or index (for arrays)
func (jv *JSONValue) Get(key interface{}) *JSONValue {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			if val, exists := v[keyStr]; exists {
				return &JSONValue{data: val}
			}
		}
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(v) {
				return &JSONValue{data: v[keyInt]}
			}
		}
	}
	return &JSONValue{data: nil}
}

// Set sets a value by key (for objects) or index (for arrays)
func (jv *JSONValue) Set(key interface{}, value interface{}) error {
	switch v := jv.data.(type) {
	case map[string]interface{}:
		if keyStr, ok := key.(string); ok {
			v[keyStr] = value
			return nil
		}
		return fmt.Errorf("key must be string for object")
	case []interface{}:
		if keyInt, ok := key.(int); ok {
			if keyInt >= 0 && keyInt < len(v) {
				v[keyInt] = value
				return nil
			}
			return fmt.Errorf("index out of range")
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
			if keyInt >= 0 && keyInt < len(v) {
				// Remove element at index
				copy(v[keyInt:], v[keyInt+1:])
				v = v[:len(v)-1]
				jv.data = v
				return nil
			}
			return fmt.Errorf("index out of range")
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
	bytes, err := json.Marshal(jv.data)
	if err != nil {
		return &JSONValue{data: nil}
	}

	var cloned interface{}
	if err := json.Unmarshal(bytes, &cloned); err != nil {
		return &JSONValue{data: nil}
	}

	return &JSONValue{data: cloned}
}

// Path retrieves a nested value using a dot-separated path
func (jv *JSONValue) Path(path string) *JSONValue {
	parts := strings.Split(path, ".")
	current := jv

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Try as array index first
		if index, err := strconv.Atoi(part); err == nil {
			current = current.Get(index)
		} else {
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
				if _, err := strconv.Atoi(parts[i+1]); err == nil {
					// Next part is an array index
					newArray := make([]interface{}, 0)
					current.Set(part, newArray)
				} else {
					// Next part is an object key
					newObj := make(map[string]interface{})
					current.Set(part, newObj)
				}
			} else {
				newObj := make(map[string]interface{})
				current.Set(part, newObj)
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
