package easyjson

import (
	"sort"
	"strings"
)

// array_operations.go - Enhanced array operations

// FindInArray searches array for item matching the predicate
// Usage: data.Get("users").FindInArray(func(user *JSONValue) bool { return user.GetString("role") == "admin" })
func (jv *JSONValue) FindInArray(matchFn func(*JSONValue) bool) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if matchFn == nil {
		return &JSONValue{data: nil}
	}
	if !jv.IsArray() {
		return &JSONValue{data: nil}
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return &JSONValue{data: nil}
	}

	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		if matchFn(item) {
			return item
		}
	}
	return &JSONValue{data: nil}
}

// FindByField finds first array item where field equals value
// Usage: data.Get("users").FindByField("id", 123)
func (jv *JSONValue) FindByField(fieldName string, value interface{}) *JSONValue {
	// Validate inputs
	if jv == nil {
		return &JSONValue{data: nil}
	}
	if err := validateStringKey(fieldName); err != nil {
		return &JSONValue{data: nil}
	}
	if err := validateJSONValueType(value); err != nil {
		return &JSONValue{data: nil}
	}

	return jv.FindInArray(func(item *JSONValue) bool {
		if item == nil {
			return false
		}
		field := item.Get(fieldName)
		if field == nil {
			return false
		}
		switch v := value.(type) {
		case string:
			return field.AsString() == v
		case int:
			return field.AsInt() == v
		case bool:
			return field.AsBool() == v
		case float64:
			return field.AsFloat() == v
		}
		return false
	})
}

// FindAllByField finds all array items where field equals value
// Usage: data.Get("users").FindAllByField("role", "admin")
func (jv *JSONValue) FindAllByField(fieldName string, value interface{}) []*JSONValue {
	// Validate inputs
	if jv == nil {
		return []*JSONValue{}
	}
	if err := validateStringKey(fieldName); err != nil {
		return []*JSONValue{}
	}
	if err := validateJSONValueType(value); err != nil {
		return []*JSONValue{}
	}
	if !jv.IsArray() {
		return []*JSONValue{}
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return []*JSONValue{}
	}

	var results []*JSONValue
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		field := item.Get(fieldName)
		if field == nil {
			continue
		}
		match := false
		switch v := value.(type) {
		case string:
			match = field.AsString() == v
		case int:
			match = field.AsInt() == v
		case bool:
			match = field.AsBool() == v
		case float64:
			match = field.AsFloat() == v
		}
		if match {
			results = append(results, item)
		}
	}
	return results
}

// FilterArray returns new JSONValue with filtered array items
// Usage: data.Get("users").FilterArray(func(user *JSONValue) bool { return user.GetBool("active") })
func (jv *JSONValue) FilterArray(filterFn func(*JSONValue) bool) *JSONValue {
	// Validate inputs
	if jv == nil {
		return NewArray()
	}
	if filterFn == nil {
		return NewArray()
	}
	if !jv.IsArray() {
		return NewArray()
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return NewArray()
	}

	var filtered []interface{}
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		if filterFn(item) {
			filtered = append(filtered, item.Raw())
		}
	}

	return &JSONValue{data: filtered}
}

// MapArray transforms array items and returns new JSONValue
// Usage: data.Get("users").MapArray(func(user *JSONValue) interface{} { return user.GetString("name") })
func (jv *JSONValue) MapArray(mapFn func(*JSONValue) interface{}) *JSONValue {
	// Validate inputs
	if jv == nil {
		return NewArray()
	}
	if mapFn == nil {
		return NewArray()
	}
	if !jv.IsArray() {
		return NewArray()
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return NewArray()
	}

	var mapped []interface{}
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			mapped = append(mapped, nil)
			continue
		}
		result := mapFn(item)
		// Validate the mapped result type
		if err := validateJSONValueType(result); err != nil {
			mapped = append(mapped, nil)
			continue
		}
		mapped = append(mapped, result)
	}

	return &JSONValue{data: mapped}
}

// ReduceArray reduces array to single value
// Usage: data.Get("numbers").ReduceArray(0, func(acc interface{}, item *JSONValue) interface{} { return acc.(int) + item.AsInt() })
func (jv *JSONValue) ReduceArray(
	initial interface{},
	reduceFn func(interface{}, *JSONValue) interface{},
) interface{} {
	// Validate inputs
	if jv == nil {
		return initial
	}
	if reduceFn == nil {
		return initial
	}
	if !jv.IsArray() {
		return initial
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return initial
	}

	accumulator := initial
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		accumulator = reduceFn(accumulator, item)
	}
	return accumulator
}

// ForEach executes function for each array item
// Usage: data.Get("users").ForEach(func(i int, user *JSONValue) { fmt.Printf("%d: %s\n", i, user.GetString("name")) })
func (jv *JSONValue) ForEach(fn func(int, *JSONValue)) {
	// Validate inputs
	if jv == nil {
		return
	}
	if fn == nil {
		return
	}
	if !jv.IsArray() {
		return
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return
	}

	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		fn(i, item)
	}
}

// Some checks if at least one array item matches predicate
// Usage: data.Get("users").Some(func(user *JSONValue) bool { return user.GetString("role") == "admin" })
func (jv *JSONValue) Some(predicateFn func(*JSONValue) bool) bool {
	// Validate inputs
	if jv == nil {
		return false
	}
	if predicateFn == nil {
		return false
	}
	if !jv.IsArray() {
		return false
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return false
	}

	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		if predicateFn(item) {
			return true
		}
	}
	return false
}

// Every checks if all array items match predicate
// Usage: data.Get("users").Every(func(user *JSONValue) bool { return user.GetBool("active") })
func (jv *JSONValue) Every(predicateFn func(*JSONValue) bool) bool {
	// Validate inputs
	if jv == nil {
		return false
	}
	if predicateFn == nil {
		return false
	}
	if !jv.IsArray() {
		return false
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return false
	}

	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			return false
		}
		if item == nil {
			return false
		}
		if !predicateFn(item) {
			return false
		}
	}
	return true
}

// Pluck extracts specified field from all array items
// Usage: data.Get("users").Pluck("name") - returns array of all names
func (jv *JSONValue) Pluck(fieldName string) *JSONValue {
	return jv.MapArray(func(item *JSONValue) interface{} {
		return item.Get(fieldName).Raw()
	})
}

// PluckStrings extracts string field from all array items
// Usage: data.Get("users").PluckStrings("name") - returns []string of names
func (jv *JSONValue) PluckStrings(fieldName string) []string {
	// Validate inputs
	if jv == nil {
		return []string{}
	}
	if err := validateStringKey(fieldName); err != nil {
		return []string{}
	}
	if !jv.IsArray() {
		return []string{}
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return []string{}
	}

	var results []string
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			results = append(results, "")
			continue
		}
		results = append(results, item.GetString(fieldName))
	}
	return results
}

// PluckInts extracts integer field from all array items
// Usage: data.Get("users").PluckInts("age") - returns []int of ages
func (jv *JSONValue) PluckInts(fieldName string) []int {
	// Validate inputs
	if jv == nil {
		return []int{}
	}
	if err := validateStringKey(fieldName); err != nil {
		return []int{}
	}
	if !jv.IsArray() {
		return []int{}
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return []int{}
	}

	var results []int
	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			results = append(results, 0)
			continue
		}
		results = append(results, item.GetInt(fieldName))
	}
	return results
}

// GroupBy groups array items by field value
// Usage: data.Get("users").GroupBy("role") - returns map[string][]*JSONValue
func (jv *JSONValue) GroupBy(fieldName string) map[string][]*JSONValue {
	groups := make(map[string][]*JSONValue)

	// Validate inputs
	if jv == nil {
		return groups
	}
	if err := validateStringKey(fieldName); err != nil {
		return groups
	}
	if !jv.IsArray() {
		return groups
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return groups
	}

	for i, item := range jv.AsArray() {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, jv.Len()); err != nil {
			continue
		}
		if item == nil {
			continue
		}
		key := item.GetString(fieldName)
		groups[key] = append(groups[key], item)
	}

	return groups
}

// SortBy sorts array by field value (returns new JSONValue)
// Usage: data.Get("users").SortBy("name") - sorts by name alphabetically
func (jv *JSONValue) SortBy(fieldName string) *JSONValue {
	if !jv.IsArray() {
		return NewArray()
	}

	items := jv.AsArray()

	// Efficient sort using Go's sort.Slice
	sorted := make([]*JSONValue, len(items))
	copy(sorted, items)

	sort.Slice(sorted, func(i, j int) bool {
		val1 := sorted[i].GetString(fieldName)
		val2 := sorted[j].GetString(fieldName)
		return strings.ToLower(val1) < strings.ToLower(val2) // Case-insensitive sort
	})

	// Convert back to interface{} slice
	result := make([]interface{}, len(sorted))
	for i, item := range sorted {
		result[i] = item.Raw()
	}

	return &JSONValue{data: result}
}

// Unique returns array with duplicate items removed
// Usage: data.Get("tags").Unique() - removes duplicate tags
func (jv *JSONValue) Unique() *JSONValue {
	if !jv.IsArray() {
		return NewArray()
	}

	seen := make(map[string]bool)
	var unique []interface{}

	for _, item := range jv.AsArray() {
		// Convert item to string for comparison
		str := item.String()
		if !seen[str] {
			seen[str] = true
			unique = append(unique, item.Raw())
		}
	}

	return &JSONValue{data: unique}
}

// First returns first array item or null if empty
// Usage: data.Get("users").First()
func (jv *JSONValue) First() *JSONValue {
	if !jv.IsArray() || jv.Len() == 0 {
		return &JSONValue{data: nil}
	}
	return jv.Get(0)
}

// Last returns last array item or null if empty
// Usage: data.Get("users").Last()
func (jv *JSONValue) Last() *JSONValue {
	if !jv.IsArray() || jv.Len() == 0 {
		return &JSONValue{data: nil}
	}
	return jv.Get(jv.Len() - 1)
}

// Take returns first N items from array
// Usage: data.Get("users").Take(5) - first 5 users
func (jv *JSONValue) Take(n int) *JSONValue {
	// Validate inputs
	if jv == nil {
		return NewArray()
	}
	if n < 0 {
		return NewArray()
	}
	if !jv.IsArray() {
		return NewArray()
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return NewArray()
	}

	length := jv.Len()
	if n > length {
		n = length
	}

	// Validate the take size
	if err := validateArraySize(n); err != nil {
		return NewArray()
	}

	var taken []interface{}
	for i := 0; i < n; i++ {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, length); err != nil {
			break
		}
		item := jv.Get(i)
		if item != nil {
			taken = append(taken, item.Raw())
		} else {
			taken = append(taken, nil)
		}
	}

	return &JSONValue{data: taken}
}

// Skip returns array without first N items
// Usage: data.Get("users").Skip(10) - all users except first 10
func (jv *JSONValue) Skip(n int) *JSONValue {
	// Validate inputs
	if jv == nil {
		return NewArray()
	}
	if n < 0 {
		return NewArray()
	}
	if !jv.IsArray() {
		return NewArray()
	}

	// Validate array size
	if err := validateArraySize(jv.Len()); err != nil {
		return NewArray()
	}

	length := jv.Len()
	if n >= length {
		return NewArray()
	}

	// Validate the skip index
	if err := validateArrayIndex(n, length); err != nil {
		return NewArray()
	}

	var remaining []interface{}
	for i := n; i < length; i++ {
		// Validate array bounds during iteration
		if err := validateArrayIndex(i, length); err != nil {
			break
		}
		item := jv.Get(i)
		if item != nil {
			remaining = append(remaining, item.Raw())
		} else {
			remaining = append(remaining, nil)
		}
	}

	return &JSONValue{data: remaining}
}
