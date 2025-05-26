package easyjson

// array_operations.go - Enhanced array operations

// FindInArray searches array for item matching the predicate
// Usage: data.Get("users").FindInArray(func(user *JSONValue) bool { return user.GetString("role") == "admin" })
func (jv *JSONValue) FindInArray(matchFn func(*JSONValue) bool) *JSONValue {
	if !jv.IsArray() {
		return &JSONValue{data: nil}
	}

	for _, item := range jv.AsArray() {
		if matchFn(item) {
			return item
		}
	}
	return &JSONValue{data: nil}
}

// FindByField finds first array item where field equals value
// Usage: data.Get("users").FindByField("id", 123)
func (jv *JSONValue) FindByField(fieldName string, value interface{}) *JSONValue {
	return jv.FindInArray(func(item *JSONValue) bool {
		field := item.Get(fieldName)
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
	if !jv.IsArray() {
		return []*JSONValue{}
	}

	var results []*JSONValue
	for _, item := range jv.AsArray() {
		field := item.Get(fieldName)
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
	if !jv.IsArray() {
		return NewArray()
	}

	var filtered []interface{}
	for _, item := range jv.AsArray() {
		if filterFn(item) {
			filtered = append(filtered, item.Raw())
		}
	}

	return &JSONValue{data: filtered}
}

// MapArray transforms array items and returns new JSONValue
// Usage: data.Get("users").MapArray(func(user *JSONValue) interface{} { return user.GetString("name") })
func (jv *JSONValue) MapArray(mapFn func(*JSONValue) interface{}) *JSONValue {
	if !jv.IsArray() {
		return NewArray()
	}

	var mapped []interface{}
	for _, item := range jv.AsArray() {
		mapped = append(mapped, mapFn(item))
	}

	return &JSONValue{data: mapped}
}

// ReduceArray reduces array to single value
// Usage: data.Get("numbers").ReduceArray(0, func(acc interface{}, item *JSONValue) interface{} { return acc.(int) + item.AsInt() })
func (jv *JSONValue) ReduceArray(
	initial interface{},
	reduceFn func(interface{}, *JSONValue) interface{},
) interface{} {
	if !jv.IsArray() {
		return initial
	}

	accumulator := initial
	for _, item := range jv.AsArray() {
		accumulator = reduceFn(accumulator, item)
	}
	return accumulator
}

// ForEach executes function for each array item
// Usage: data.Get("users").ForEach(func(i int, user *JSONValue) { fmt.Printf("%d: %s\n", i, user.GetString("name")) })
func (jv *JSONValue) ForEach(fn func(int, *JSONValue)) {
	if !jv.IsArray() {
		return
	}

	for i, item := range jv.AsArray() {
		fn(i, item)
	}
}

// Some checks if at least one array item matches predicate
// Usage: data.Get("users").Some(func(user *JSONValue) bool { return user.GetString("role") == "admin" })
func (jv *JSONValue) Some(predicateFn func(*JSONValue) bool) bool {
	if !jv.IsArray() {
		return false
	}

	for _, item := range jv.AsArray() {
		if predicateFn(item) {
			return true
		}
	}
	return false
}

// Every checks if all array items match predicate
// Usage: data.Get("users").Every(func(user *JSONValue) bool { return user.GetBool("active") })
func (jv *JSONValue) Every(predicateFn func(*JSONValue) bool) bool {
	if !jv.IsArray() {
		return false
	}

	for _, item := range jv.AsArray() {
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
	if !jv.IsArray() {
		return []string{}
	}

	var results []string
	for _, item := range jv.AsArray() {
		results = append(results, item.GetString(fieldName))
	}
	return results
}

// PluckInts extracts integer field from all array items
// Usage: data.Get("users").PluckInts("age") - returns []int of ages
func (jv *JSONValue) PluckInts(fieldName string) []int {
	if !jv.IsArray() {
		return []int{}
	}

	var results []int
	for _, item := range jv.AsArray() {
		results = append(results, item.GetInt(fieldName))
	}
	return results
}

// GroupBy groups array items by field value
// Usage: data.Get("users").GroupBy("role") - returns map[string][]*JSONValue
func (jv *JSONValue) GroupBy(fieldName string) map[string][]*JSONValue {
	groups := make(map[string][]*JSONValue)

	if !jv.IsArray() {
		return groups
	}

	for _, item := range jv.AsArray() {
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

	// Simple bubble sort for now (can be optimized later)
	n := len(items)
	sorted := make([]*JSONValue, n)
	copy(sorted, items)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			val1 := sorted[j].GetString(fieldName)
			val2 := sorted[j+1].GetString(fieldName)
			if val1 > val2 {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	// Convert back to interface{} slice
	var result []interface{}
	for _, item := range sorted {
		result = append(result, item.Raw())
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
	if !jv.IsArray() {
		return NewArray()
	}

	length := jv.Len()
	if n > length {
		n = length
	}

	var taken []interface{}
	for i := 0; i < n; i++ {
		taken = append(taken, jv.Get(i).Raw())
	}

	return &JSONValue{data: taken}
}

// Skip returns array without first N items
// Usage: data.Get("users").Skip(10) - all users except first 10
func (jv *JSONValue) Skip(n int) *JSONValue {
	if !jv.IsArray() {
		return NewArray()
	}

	length := jv.Len()
	if n >= length {
		return NewArray()
	}

	var remaining []interface{}
	for i := n; i < length; i++ {
		remaining = append(remaining, jv.Get(i).Raw())
	}

	return &JSONValue{data: remaining}
}
