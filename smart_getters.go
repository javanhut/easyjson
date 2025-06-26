package easyjson

// smart_getters.go - Smart default value getters

// GetString gets a string value with optional default
// Usage: data.GetString("user", "name", "Anonymous")
func (jv *JSONValue) GetString(keys ...interface{}) string {
	// Validate inputs
	if jv == nil {
		return ""
	}
	if len(keys) == 0 {
		return ""
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return ""
	}

	var defaultValue string
	var path []interface{}

	// Last parameter is default if it's a string and we have more than one param
	if len(keys) > 1 {
		if def, ok := keys[len(keys)-1].(string); ok {
			// Validate default string length
			if err := validateStringLength(def); err != nil {
				defaultValue = "" // Use empty string if default is too long
			} else {
				defaultValue = def
			}
			path = keys[:len(keys)-1]
		} else {
			path = keys
		}
	} else {
		path = keys
	}

	// Validate each key in path
	for _, key := range path {
		if err := validateKey(key); err != nil {
			return defaultValue
		}
	}

	result := jv.Q(path...)
	if result != nil && !result.IsNull() {
		resultStr := result.AsString()
		// Validate result string length
		if err := validateStringLength(resultStr); err != nil {
			return defaultValue
		}
		return resultStr
	}
	return defaultValue
}

// GetInt gets an integer value with optional default
// Usage: data.GetInt("user", "age", 18)
func (jv *JSONValue) GetInt(keys ...interface{}) int {
	// Validate inputs
	if jv == nil {
		return 0
	}
	if len(keys) == 0 {
		return 0
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return 0
	}

	var defaultValue int
	var path []interface{}

	if len(keys) > 1 {
		if def, ok := keys[len(keys)-1].(int); ok {
			defaultValue = def
			path = keys[:len(keys)-1]
		} else {
			path = keys
		}
	} else {
		path = keys
	}

	// Validate each key in path
	for _, key := range path {
		if err := validateKey(key); err != nil {
			return defaultValue
		}
	}

	result := jv.Q(path...)
	if result != nil && !result.IsNull() {
		return result.AsInt()
	}
	return defaultValue
}

// GetBool gets a boolean value with optional default
// Usage: data.GetBool("user", "active", true)
func (jv *JSONValue) GetBool(keys ...interface{}) bool {
	// Validate inputs
	if jv == nil {
		return false
	}
	if len(keys) == 0 {
		return false
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return false
	}

	var defaultValue bool
	var path []interface{}

	if len(keys) > 1 {
		if def, ok := keys[len(keys)-1].(bool); ok {
			defaultValue = def
			path = keys[:len(keys)-1]
		} else {
			path = keys
		}
	} else {
		path = keys
	}

	// Validate each key in path
	for _, key := range path {
		if err := validateKey(key); err != nil {
			return defaultValue
		}
	}

	result := jv.Q(path...)
	if result != nil && !result.IsNull() {
		return result.AsBool()
	}
	return defaultValue
}

// GetFloat gets a float64 value with optional default
// Usage: data.GetFloat("user", "rating", 0.0)
func (jv *JSONValue) GetFloat(keys ...interface{}) float64 {
	// Validate inputs
	if jv == nil {
		return 0.0
	}
	if len(keys) == 0 {
		return 0.0
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return 0.0
	}

	var defaultValue float64
	var path []interface{}

	if len(keys) > 1 {
		if def, ok := keys[len(keys)-1].(float64); ok {
			defaultValue = def
			path = keys[:len(keys)-1]
		} else {
			path = keys
		}
	} else {
		path = keys
	}

	// Validate each key in path
	for _, key := range path {
		if err := validateKey(key); err != nil {
			return defaultValue
		}
	}

	result := jv.Q(path...)
	if result != nil && !result.IsNull() {
		return result.AsFloat()
	}
	return defaultValue
}

// GetOr gets a value with smart type-matched default
// Usage: data.GetOr("user", "name", "Anonymous")
func (jv *JSONValue) GetOr(keys ...interface{}) interface{} {
	// Validate inputs
	if jv == nil {
		return nil
	}
	if len(keys) < 2 {
		return nil
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return nil
	}

	defaultValue := keys[len(keys)-1]
	path := keys[:len(keys)-1]

	// Validate default value type
	if err := validateJSONValueType(defaultValue); err != nil {
		return nil
	}

	// Validate each key in path
	for _, key := range path {
		if err := validateKey(key); err != nil {
			return defaultValue
		}
	}

	result := jv.Q(path...)
	if result != nil && !result.IsNull() {
		// Smart type matching with validation
		switch defaultValue.(type) {
		case string:
			resultStr := result.AsString()
			if err := validateStringLength(resultStr); err != nil {
				return defaultValue
			}
			return resultStr
		case int:
			return result.AsInt()
		case float64:
			return result.AsFloat()
		case bool:
			return result.AsBool()
		default:
			return result.Raw()
		}
	}
	return defaultValue
}

// IsEmptyOrNull checks if value is effectively empty
func (jv *JSONValue) IsEmptyOrNull() bool {
	if jv.IsNull() {
		return true
	}
	if jv.IsString() && jv.AsString() == "" {
		return true
	}
	if (jv.IsArray() || jv.IsObject()) && jv.Len() == 0 {
		return true
	}
	return false
}

// StringOrEmpty returns string value or empty string if null/missing
func (jv *JSONValue) StringOrEmpty() string {
	if jv.IsNull() {
		return ""
	}
	return jv.AsString()
}

// IntOrZero returns int value or 0 if null/missing
func (jv *JSONValue) IntOrZero() int {
	if jv.IsNull() {
		return 0
	}
	return jv.AsInt()
}

// BoolOrFalse returns bool value or false if null/missing
func (jv *JSONValue) BoolOrFalse() bool {
	if jv.IsNull() {
		return false
	}
	return jv.AsBool()
}

// FloatOrZero returns float64 value or 0.0 if null/missing
func (jv *JSONValue) FloatOrZero() float64 {
	if jv.IsNull() {
		return 0.0
	}
	return jv.AsFloat()
}
