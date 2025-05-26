package easyjson

import "strings"

// multi_path.go - Multi-path access and robust querying

// TryPaths attempts multiple paths until one returns a non-null result
// Usage: data.TryPaths("title", "name", "label", "header")
func (jv *JSONValue) TryPaths(paths ...string) *JSONValue {
	for _, path := range paths {
		if result := jv.Path(path); !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// TryKeys attempts multiple keys at current level until one works
// Usage: data.TryKeys("name", "title", "label")
func (jv *JSONValue) TryKeys(keys ...string) *JSONValue {
	for _, key := range keys {
		if result := jv.Get(key); !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// TryQueries attempts multiple Q-style queries until one works
// Usage: data.TryQueries([]interface{}{"user", "name"}, []interface{}{"profile", "name"})
func (jv *JSONValue) TryQueries(queries ...[]interface{}) *JSONValue {
	for _, query := range queries {
		if result := jv.Q(query...); !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// DeepSearch searches for a key at any depth in the JSON structure
// Usage: data.DeepSearch("email") - finds first "email" key anywhere
func (jv *JSONValue) DeepSearch(key string) *JSONValue {
	return jv.deepSearchRecursive(key, 0, 10) // Max depth of 10 to prevent infinite loops
}

func (jv *JSONValue) deepSearchRecursive(key string, currentDepth, maxDepth int) *JSONValue {
	if currentDepth > maxDepth {
		return &JSONValue{data: nil}
	}

	// Check current level first
	if jv.Has(key) {
		return jv.Get(key)
	}

	// Search in nested objects
	if jv.IsObject() {
		for _, k := range jv.Keys() {
			child := jv.Get(k)
			if result := child.deepSearchRecursive(key, currentDepth+1, maxDepth); !result.IsNull() {
				return result
			}
		}
	}

	// Search in arrays
	if jv.IsArray() {
		for i := 0; i < jv.Len(); i++ {
			child := jv.Get(i)
			if result := child.deepSearchRecursive(key, currentDepth+1, maxDepth); !result.IsNull() {
				return result
			}
		}
	}

	return &JSONValue{data: nil}
}

// DeepSearchAll finds all occurrences of a key at any depth
// Usage: data.DeepSearchAll("id") - returns all "id" values found
func (jv *JSONValue) DeepSearchAll(key string) []*JSONValue {
	var results []*JSONValue
	jv.deepSearchAllRecursive(key, &results, 0, 10)
	return results
}

func (jv *JSONValue) deepSearchAllRecursive(
	key string,
	results *[]*JSONValue,
	currentDepth, maxDepth int,
) {
	if currentDepth > maxDepth {
		return
	}

	// Check current level
	if jv.Has(key) {
		*results = append(*results, jv.Get(key))
	}

	// Search in nested objects
	if jv.IsObject() {
		for _, k := range jv.Keys() {
			child := jv.Get(k)
			child.deepSearchAllRecursive(key, results, currentDepth+1, maxDepth)
		}
	}

	// Search in arrays
	if jv.IsArray() {
		for i := 0; i < jv.Len(); i++ {
			child := jv.Get(i)
			child.deepSearchAllRecursive(key, results, currentDepth+1, maxDepth)
		}
	}
}

// FindPath returns the path to the first occurrence of a key
// Usage: data.FindPath("email") might return "user.profile.email"
func (jv *JSONValue) FindPath(key string) string {
	path := jv.findPathRecursive(key, "", 0, 10)
	return strings.TrimPrefix(path, ".")
}

func (jv *JSONValue) findPathRecursive(key, currentPath string, currentDepth, maxDepth int) string {
	if currentDepth > maxDepth {
		return ""
	}

	// Check current level
	if jv.Has(key) {
		return currentPath + "." + key
	}

	// Search in nested objects
	if jv.IsObject() {
		for _, k := range jv.Keys() {
			child := jv.Get(k)
			newPath := currentPath + "." + k
			if result := child.findPathRecursive(key, newPath, currentDepth+1, maxDepth); result != "" {
				return result
			}
		}
	}

	// Search in arrays
	if jv.IsArray() {
		for i := 0; i < jv.Len(); i++ {
			child := jv.Get(i)
			newPath := currentPath + "." + string(rune(i+'0')) // Simple index conversion
			if result := child.findPathRecursive(key, newPath, currentDepth+1, maxDepth); result != "" {
				return result
			}
		}
	}

	return ""
}

// HasAnyKey checks if any of the provided keys exist at current level
// Usage: data.HasAnyKey("name", "title", "label")
func (jv *JSONValue) HasAnyKey(keys ...string) bool {
	for _, key := range keys {
		if jv.Has(key) {
			return true
		}
	}
	return false
}

// HasAllKeys checks if all provided keys exist at current level
// Usage: data.HasAllKeys("name", "email", "id")
func (jv *JSONValue) HasAllKeys(keys ...string) bool {
	for _, key := range keys {
		if !jv.Has(key) {
			return false
		}
	}
	return true
}

// GetFirstAvailable returns the first non-null value from multiple keys
// Usage: data.GetFirstAvailable("name", "title", "label")
func (jv *JSONValue) GetFirstAvailable(keys ...string) *JSONValue {
	for _, key := range keys {
		if result := jv.Get(key); !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}
