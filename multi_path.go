package easyjson

import (
	"strconv"
	"strings"
)

// multi_path.go - Multi-path access and robust querying

// SearchConfig holds configuration for search operations
type SearchConfig struct {
	MaxDepth          int  // Maximum recursion depth
	EnableCycleDetect bool // Enable cycle detection (uses more memory)
}

// DefaultSearchConfig returns sensible defaults for search operations
func DefaultSearchConfig() *SearchConfig {
	return &SearchConfig{
		MaxDepth:          50,  // Increased from 10 for more flexibility
		EnableCycleDetect: false, // Disabled by default for performance
	}
}

// deepSearchState holds state for cycle detection
type deepSearchState struct {
	visited map[interface{}]bool
	config  *SearchConfig
}

// TryPaths attempts multiple paths until one returns a non-null result
// Usage: data.TryPaths("title", "name", "label", "header")
func (jv *JSONValue) TryPaths(paths ...string) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if len(paths) == 0 {
		return &JSONValue{data: nil}
	}

	// Validate batch size
	if err := validateBatchSize(len(paths)); err != nil {
		return &JSONValue{data: nil}
	}

	for _, path := range paths {
		// Validate each path
		if err := validatePath(path); err != nil {
			continue
		}
		if result := jv.Path(path); result != nil && !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// TryKeys attempts multiple keys at current level until one works
// Usage: data.TryKeys("name", "title", "label")
func (jv *JSONValue) TryKeys(keys ...string) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if len(keys) == 0 {
		return &JSONValue{data: nil}
	}

	// Validate batch size
	if err := validateBatchSize(len(keys)); err != nil {
		return &JSONValue{data: nil}
	}

	for _, key := range keys {
		// Validate each key
		if err := validateStringKey(key); err != nil {
			continue
		}
		if result := jv.Get(key); result != nil && !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// TryQueries attempts multiple Q-style queries until one works
// Usage: data.TryQueries([]interface{}{"user", "name"}, []interface{}{"profile", "name"})
func (jv *JSONValue) TryQueries(queries ...[]interface{}) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if len(queries) == 0 {
		return &JSONValue{data: nil}
	}

	// Validate batch size
	if err := validateBatchSize(len(queries)); err != nil {
		return &JSONValue{data: nil}
	}

	for _, query := range queries {
		if query == nil || len(query) == 0 {
			continue
		}
		// Validate each key in the query
		validQuery := true
		for _, key := range query {
			if err := validateKey(key); err != nil {
				validQuery = false
				break
			}
		}
		if !validQuery {
			continue
		}
		if result := jv.Q(query...); result != nil && !result.IsNull() {
			return result
		}
	}
	return &JSONValue{data: nil}
}

// DeepSearch searches for a key at any depth in the JSON structure
// Usage: data.DeepSearch("email") - finds first "email" key anywhere
func (jv *JSONValue) DeepSearch(key string) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if err := validateStringKey(key); err != nil {
		return &JSONValue{data: nil}
	}
	return jv.DeepSearchWithConfig(key, DefaultSearchConfig())
}

// DeepSearchWithConfig searches for a key with custom configuration
// Usage: data.DeepSearchWithConfig("email", &SearchConfig{MaxDepth: 20})
func (jv *JSONValue) DeepSearchWithConfig(key string, config *SearchConfig) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if err := validateStringKey(key); err != nil {
		return &JSONValue{data: nil}
	}
	if config == nil {
		config = DefaultSearchConfig()
	}

	// Validate recursion depth
	if err := validateRecursionDepth(0, config.MaxDepth); err != nil {
		return &JSONValue{data: nil}
	}
	
	if config.EnableCycleDetect {
		state := &deepSearchState{
			visited: make(map[interface{}]bool),
			config:  config,
		}
		return jv.deepSearchWithState(key, 0, state)
	}
	
	return jv.deepSearchRecursive(key, 0, config.MaxDepth)
}

func (jv *JSONValue) deepSearchRecursive(key string, currentDepth, maxDepth int) *JSONValue {
	// Validate recursion depth
	if err := validateRecursionDepth(currentDepth, maxDepth); err != nil {
		return &JSONValue{data: nil}
	}
	if jv == nil {
		return &JSONValue{data: nil}
	}

	// Check current level first
	if jv.Has(key) {
		return jv.Get(key)
	}

	// Search in nested objects
	if jv.IsObject() {
		// Validate object size
		if err := validateObjectKeyCount(jv.Len()); err != nil {
			return &JSONValue{data: nil}
		}
		for _, k := range jv.Keys() {
			// Validate key
			if err := validateStringKey(k); err != nil {
				continue
			}
			child := jv.Get(k)
			if child == nil {
				continue
			}
			if result := child.deepSearchRecursive(key, currentDepth+1, maxDepth); result != nil && !result.IsNull() {
				return result
			}
		}
	}

	// Search in arrays
	if jv.IsArray() {
		// Validate array size
		if err := validateArraySize(jv.Len()); err != nil {
			return &JSONValue{data: nil}
		}
		for i := 0; i < jv.Len(); i++ {
			// Validate array bounds
			if err := validateArrayIndex(i, jv.Len()); err != nil {
				break
			}
			child := jv.Get(i)
			if child == nil {
				continue
			}
			if result := child.deepSearchRecursive(key, currentDepth+1, maxDepth); result != nil && !result.IsNull() {
				return result
			}
		}
	}

	return &JSONValue{data: nil}
}

// DeepSearchAll finds all occurrences of a key at any depth
// Usage: data.DeepSearchAll("id") - returns all "id" values found
func (jv *JSONValue) DeepSearchAll(key string) []*JSONValue {
	return jv.DeepSearchAllWithConfig(key, DefaultSearchConfig())
}

// DeepSearchAllWithConfig finds all occurrences with custom configuration
// Usage: data.DeepSearchAllWithConfig("id", &SearchConfig{MaxDepth: 30})
func (jv *JSONValue) DeepSearchAllWithConfig(key string, config *SearchConfig) []*JSONValue {
	if config == nil {
		config = DefaultSearchConfig()
	}
	
	var results []*JSONValue
	if config.EnableCycleDetect {
		state := &deepSearchState{
			visited: make(map[interface{}]bool),
			config:  config,
		}
		jv.deepSearchAllWithState(key, &results, 0, state)
	} else {
		jv.deepSearchAllRecursive(key, &results, 0, config.MaxDepth)
	}
	
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
			newPath := currentPath + "." + strconv.Itoa(i)
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

// deepSearchWithState performs deep search with cycle detection
func (jv *JSONValue) deepSearchWithState(key string, currentDepth int, state *deepSearchState) *JSONValue {
	if currentDepth > state.config.MaxDepth {
		return &JSONValue{data: nil}
	}

	// Check for cycles if enabled
	if state.config.EnableCycleDetect {
		if state.visited[jv.data] {
			return &JSONValue{data: nil} // Cycle detected
		}
		state.visited[jv.data] = true
		defer delete(state.visited, jv.data) // Clean up after recursion
	}

	// Check current level first
	if jv.Has(key) {
		return jv.Get(key)
	}

	// Search in nested objects
	if jv.IsObject() {
		for _, k := range jv.Keys() {
			child := jv.Get(k)
			if result := child.deepSearchWithState(key, currentDepth+1, state); !result.IsNull() {
				return result
			}
		}
	}

	// Search in arrays
	if jv.IsArray() {
		for i := 0; i < jv.Len(); i++ {
			child := jv.Get(i)
			if result := child.deepSearchWithState(key, currentDepth+1, state); !result.IsNull() {
				return result
			}
		}
	}

	return &JSONValue{data: nil}
}

// deepSearchAllWithState performs deep search all with cycle detection
func (jv *JSONValue) deepSearchAllWithState(key string, results *[]*JSONValue, currentDepth int, state *deepSearchState) {
	if currentDepth > state.config.MaxDepth {
		return
	}

	// Check for cycles if enabled
	if state.config.EnableCycleDetect {
		if state.visited[jv.data] {
			return // Cycle detected
		}
		state.visited[jv.data] = true
		defer delete(state.visited, jv.data) // Clean up after recursion
	}

	// Check current level
	if jv.Has(key) {
		*results = append(*results, jv.Get(key))
	}

	// Search in nested objects
	if jv.IsObject() {
		for _, k := range jv.Keys() {
			child := jv.Get(k)
			child.deepSearchAllWithState(key, results, currentDepth+1, state)
		}
	}

	// Search in arrays
	if jv.IsArray() {
		for i := 0; i < jv.Len(); i++ {
			child := jv.Get(i)
			child.deepSearchAllWithState(key, results, currentDepth+1, state)
		}
	}
}
