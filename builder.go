package easyjson

import (
	"time"
)

// builder.go - Enhanced fluent JSON building

// JSONBuilder provides fluent interface for building JSON structures
type JSONBuilder struct {
	value *JSONValue
}

// NewBuilder creates a new object builder
// Usage: builder := easyjson.NewBuilder()
func NewBuilder() *JSONBuilder {
	return &JSONBuilder{value: NewObject()}
}

// NewArrayBuilder creates a new array builder
// Usage: builder := easyjson.NewArrayBuilder()
func NewArrayBuilder() *JSONBuilder {
	return &JSONBuilder{value: NewArray()}
}

// AddField adds a field to the JSON object
// Usage: builder.AddField("name", "John")
func (jb *JSONBuilder) AddField(key string, value interface{}) *JSONBuilder {
	jb.value.Set(key, value)
	return jb
}

// AddFields adds multiple fields at once from a map
// Usage: builder.AddFields(map[string]interface{}{"name": "John", "age": 30})
func (jb *JSONBuilder) AddFields(fields map[string]interface{}) *JSONBuilder {
	for key, value := range fields {
		jb.value.Set(key, value)
	}
	return jb
}

// AddItem adds an item to the JSON array
// Usage: builder.AddItem("value")
func (jb *JSONBuilder) AddItem(value interface{}) *JSONBuilder {
	jb.value.Append(value)
	return jb
}

// AddItems adds multiple items at once
// Usage: builder.AddItems("item1", "item2", "item3")
func (jb *JSONBuilder) AddItems(values ...interface{}) *JSONBuilder {
	for _, value := range values {
		jb.value.Append(value)
	}
	return jb
}

// AddObject adds a nested object using a builder function
// Usage: builder.AddObject("user", func(user *JSONBuilder) { user.AddField("name", "John") })
func (jb *JSONBuilder) AddObject(key string, builderFn func(*JSONBuilder)) *JSONBuilder {
	nested := NewBuilder()
	builderFn(nested)
	jb.value.Set(key, nested.value.Raw())
	return jb
}

// AddArray adds a nested array using a builder function
// Usage: builder.AddArray("items", func(arr *JSONBuilder) { arr.AddItem("item1") })
func (jb *JSONBuilder) AddArray(key string, builderFn func(*JSONBuilder)) *JSONBuilder {
	nested := NewArrayBuilder()
	builderFn(nested)
	jb.value.Set(key, nested.value.Raw())
	return jb
}

// AddIf conditionally adds a field
// Usage: builder.AddIf(user.IsAdmin, "admin_data", adminInfo)
func (jb *JSONBuilder) AddIf(condition bool, key string, value interface{}) *JSONBuilder {
	if condition {
		jb.value.Set(key, value)
	}
	return jb
}

// AddIfNotEmpty adds field only if value is not empty
// Usage: builder.AddIfNotEmpty("description", description)
func (jb *JSONBuilder) AddIfNotEmpty(key string, value interface{}) *JSONBuilder {
	switch v := value.(type) {
	case string:
		if v != "" {
			jb.value.Set(key, value)
		}
	case []interface{}:
		if len(v) > 0 {
			jb.value.Set(key, value)
		}
	case map[string]interface{}:
		if len(v) > 0 {
			jb.value.Set(key, value)
		}
	default:
		if value != nil {
			jb.value.Set(key, value)
		}
	}
	return jb
}

// AddTimestamp adds current timestamp
// Usage: builder.AddTimestamp("created_at")
func (jb *JSONBuilder) AddTimestamp(key string) *JSONBuilder {
	jb.value.Set(key, time.Now().Unix())
	return jb
}

// AddISO8601Timestamp adds current timestamp in ISO8601 format
// Usage: builder.AddISO8601Timestamp("created_at")
func (jb *JSONBuilder) AddISO8601Timestamp(key string) *JSONBuilder {
	jb.value.Set(key, time.Now().UTC().Format(time.RFC3339))
	return jb
}

// AddUserInfo adds common user information fields
// Usage: builder.AddUserInfo(user)
func (jb *JSONBuilder) AddUserInfo(user interface{}) *JSONBuilder {
	// This is a placeholder - in real implementation, you'd use reflection
	// or type assertion to extract user fields
	if userMap, ok := user.(map[string]interface{}); ok {
		commonFields := []string{"id", "name", "email", "username", "role"}
		for _, field := range commonFields {
			if value, exists := userMap[field]; exists {
				jb.value.Set(field, value)
			}
		}
	}
	return jb
}

// AddPaginationInfo adds common pagination fields
// Usage: builder.AddPaginationInfo(page, total, limit)
func (jb *JSONBuilder) AddPaginationInfo(page, total, limit int) *JSONBuilder {
	return jb.AddObject("pagination", func(pag *JSONBuilder) {
		pag.AddField("page", page).
			AddField("total", total).
			AddField("limit", limit).
			AddField("has_next", (page*limit) < total).
			AddField("has_prev", page > 1)
	})
}

// AddAPIStatus adds standard API response status
// Usage: builder.AddAPIStatus("success", "Operation completed")
func (jb *JSONBuilder) AddAPIStatus(status, message string) *JSONBuilder {
	return jb.AddField("status", status).
		AddField("message", message).
		AddTimestamp("timestamp")
}

// AddError adds error information
// Usage: builder.AddError("VALIDATION_ERROR", "Invalid input", errors)
func (jb *JSONBuilder) AddError(code, message string, details interface{}) *JSONBuilder {
	return jb.AddObject("error", func(err *JSONBuilder) {
		err.AddField("code", code).
			AddField("message", message).
			AddIfNotEmpty("details", details)
	})
}

// Merge merges another JSONValue into this builder
// Usage: builder.Merge(otherJSONValue)
func (jb *JSONBuilder) Merge(other *JSONValue) *JSONBuilder {
	if other.IsObject() {
		for key, value := range other.AsObject() {
			jb.value.Set(key, value.Raw())
		}
	}
	return jb
}

// When allows conditional building
// Usage: builder.When(condition, func(b *JSONBuilder) { b.AddField("extra", "data") })
func (jb *JSONBuilder) When(condition bool, fn func(*JSONBuilder)) *JSONBuilder {
	if condition {
		fn(jb)
	}
	return jb
}

// Unless is the opposite of When
// Usage: builder.Unless(isGuest, func(b *JSONBuilder) { b.AddField("admin_data", data) })
func (jb *JSONBuilder) Unless(condition bool, fn func(*JSONBuilder)) *JSONBuilder {
	if !condition {
		fn(jb)
	}
	return jb
}

// ToJSON returns the built JSONValue
// Usage: jsonValue := builder.ToJSON()
func (jb *JSONBuilder) ToJSON() *JSONValue {
	return jb.value
}

// ToJSONString returns the JSON as a string
// Usage: jsonString := builder.ToJSONString()
func (jb *JSONBuilder) ToJSONString() string {
	result, _ := jb.value.Dumps()
	return result
}

// ToPrettyString returns pretty-printed JSON string
// Usage: prettyJSON := builder.ToPrettyString()
func (jb *JSONBuilder) ToPrettyString() string {
	result, _ := jb.value.DumpsIndent("  ")
	return result
}

// ToBytes returns JSON as byte slice
// Usage: jsonBytes := builder.ToBytes()
func (jb *JSONBuilder) ToBytes() []byte {
	result, _ := jb.value.Dump()
	return result
}

// Size returns the approximate size of the JSON in bytes
// Usage: size := builder.Size()
func (jb *JSONBuilder) Size() int {
	bytes, _ := jb.value.Dump()
	return len(bytes)
}

// Validate checks if the built JSON is valid
// Usage: isValid := builder.Validate()
func (jb *JSONBuilder) Validate() bool {
	_, err := jb.value.Dumps()
	return err == nil
}

// Clone creates a deep copy of the builder
// Usage: newBuilder := builder.Clone()
func (jb *JSONBuilder) Clone() *JSONBuilder {
	return &JSONBuilder{value: jb.value.Clone()}
}

// Reset clears the builder and starts fresh
// Usage: builder.Reset()
func (jb *JSONBuilder) Reset() *JSONBuilder {
	if jb.value.IsArray() {
		jb.value = NewArray()
	} else {
		jb.value = NewObject()
	}
	return jb
}

// Quick builder functions for common patterns

// QuickObject creates object with key-value pairs
// Usage: obj := easyjson.QuickObject("name", "John", "age", 30)
func QuickObject(pairs ...interface{}) *JSONValue {
	if len(pairs)%2 != 0 {
		return NewObject() // Return empty on odd number
	}

	builder := NewBuilder()
	for i := 0; i < len(pairs); i += 2 {
		if key, ok := pairs[i].(string); ok {
			builder.AddField(key, pairs[i+1])
		}
	}
	return builder.ToJSON()
}

// QuickArray creates array with items
// Usage: arr := easyjson.QuickArray("item1", "item2", "item3")
func QuickArray(items ...interface{}) *JSONValue {
	builder := NewArrayBuilder()
	builder.AddItems(items...)
	return builder.ToJSON()
}

// QuickAPIResponse creates standard API response
// Usage: response := easyjson.QuickAPIResponse("success", data, "Operation completed")
func QuickAPIResponse(status string, data interface{}, message string) *JSONValue {
	return NewBuilder().
		AddField("status", status).
		AddField("data", data).
		AddField("message", message).
		AddTimestamp("timestamp").
		ToJSON()
}

// QuickErrorResponse creates standard error response
// Usage: response := easyjson.QuickErrorResponse("INVALID_INPUT", "Validation failed", errors)
func QuickErrorResponse(code, message string, details interface{}) *JSONValue {
	builder := NewBuilder().
		AddField("status", "error").
		AddTimestamp("timestamp").
		AddObject("error", func(err *JSONBuilder) {
			err.AddField("code", code).
				AddField("message", message)
			if details != nil {
				err.AddField("details", details)
			}
		})
	return builder.ToJSON()
}

// QuickPaginatedResponse creates paginated response
// Usage: response := easyjson.QuickPaginatedResponse(items, 1, 100, 10)
func QuickPaginatedResponse(data interface{}, page, total, limit int) *JSONValue {
	return NewBuilder().
		AddField("status", "success").
		AddField("data", data).
		AddPaginationInfo(page, total, limit).
		AddTimestamp("timestamp").
		ToJSON()
}
