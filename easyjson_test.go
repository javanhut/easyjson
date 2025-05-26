package easyjson

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestLoads(t *testing.T) {
	jsonStr := `{"name": "John", "age": 30, "active": true}`
	jv, err := Loads(jsonStr)
	if err != nil {
		t.Fatalf("Loads failed: %v", err)
	}

	if jv.Get("name").AsString() != "John" {
		t.Errorf("Expected name 'John', got '%s'", jv.Get("name").AsString())
	}

	if jv.Get("age").AsInt() != 30 {
		t.Errorf("Expected age 30, got %d", jv.Get("age").AsInt())
	}

	if !jv.Get("active").AsBool() {
		t.Errorf("Expected active true, got false")
	}
}

func TestLoadsInvalidJSON(t *testing.T) {
	invalidJSON := `{"name": "John", "age":}`
	_, err := Loads(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoad(t *testing.T) {
	jsonBytes := []byte(`{"test": "value"}`)
	jv, err := Load(jsonBytes)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if jv.Get("test").AsString() != "value" {
		t.Errorf("Expected 'value', got '%s'", jv.Get("test").AsString())
	}
}

func TestDumps(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jv := New(data)

	result, err := jv.Dumps()
	if err != nil {
		t.Fatalf("Dumps failed: %v", err)
	}

	// Parse it back to verify
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		t.Fatalf("Failed to parse dumped JSON: %v", err)
	}

	if parsed["name"] != "John" || parsed["age"].(float64) != 30 {
		t.Error("Dumped JSON doesn't match original data")
	}
}

func TestDumpsIndent(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jv := New(data)

	result, err := jv.DumpsIndent("  ")
	if err != nil {
		t.Fatalf("DumpsIndent failed: %v", err)
	}

	if !strings.Contains(result, "  ") {
		t.Error("Indented JSON should contain indentation")
	}
}

func TestGet(t *testing.T) {
	// Test object access
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jv := New(data)

	name := jv.Get("name")
	if name.AsString() != "John" {
		t.Errorf("Expected 'John', got '%s'", name.AsString())
	}

	// Test array access
	arr := []interface{}{"a", "b", "c"}
	jvArr := New(arr)

	first := jvArr.Get(0)
	if first.AsString() != "a" {
		t.Errorf("Expected 'a', got '%s'", first.AsString())
	}

	// Test nonexistent key
	nonexistent := jv.Get("nonexistent")
	if !nonexistent.IsNull() {
		t.Error("Expected null for nonexistent key")
	}
}

func TestSet(t *testing.T) {
	// Test object set
	data := map[string]interface{}{"name": "John"}
	jv := New(data)

	err := jv.Set("age", 30)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if jv.Get("age").AsInt() != 30 {
		t.Error("Set didn't work correctly")
	}

	// Test array set
	arr := []interface{}{"a", "b", "c"}
	jvArr := New(arr)

	err = jvArr.Set(1, "modified")
	if err != nil {
		t.Fatalf("Array set failed: %v", err)
	}

	if jvArr.Get(1).AsString() != "modified" {
		t.Error("Array set didn't work correctly")
	}
}

func TestHas(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jv := New(data)

	if !jv.Has("name") {
		t.Error("Has should return true for existing key")
	}

	if jv.Has("nonexistent") {
		t.Error("Has should return false for nonexistent key")
	}

	// Test array
	arr := []interface{}{"a", "b", "c"}
	jvArr := New(arr)

	if !jvArr.Has(0) {
		t.Error("Has should return true for valid array index")
	}

	if jvArr.Has(10) {
		t.Error("Has should return false for invalid array index")
	}
}

func TestDelete(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
	}
	jv := New(data)

	err := jv.Delete("age")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if jv.Has("age") {
		t.Error("Key should be deleted")
	}
}

func TestTypeChecking(t *testing.T) {
	tests := []struct {
		value    interface{}
		isString bool
		isNumber bool
		isBool   bool
		isArray  bool
		isObject bool
		isNull   bool
	}{
		{"hello", true, false, false, false, false, false},
		{42.0, false, true, false, false, false, false},
		{true, false, false, true, false, false, false},
		{[]interface{}{1, 2, 3}, false, false, false, true, false, false},
		{map[string]interface{}{"key": "value"}, false, false, false, false, true, false},
		{nil, false, false, false, false, false, true},
	}

	for _, test := range tests {
		jv := New(test.value)

		if jv.IsString() != test.isString {
			t.Errorf("IsString() failed for %v", test.value)
		}
		if jv.IsNumber() != test.isNumber {
			t.Errorf("IsNumber() failed for %v", test.value)
		}
		if jv.IsBool() != test.isBool {
			t.Errorf("IsBool() failed for %v", test.value)
		}
		if jv.IsArray() != test.isArray {
			t.Errorf("IsArray() failed for %v", test.value)
		}
		if jv.IsObject() != test.isObject {
			t.Errorf("IsObject() failed for %v", test.value)
		}
		if jv.IsNull() != test.isNull {
			t.Errorf("IsNull() failed for %v", test.value)
		}
	}
}

func TestTypeConversion(t *testing.T) {
	// Test string conversion
	jv := New("42")
	if jv.AsInt() != 42 {
		t.Error("String to int conversion failed")
	}

	if jv.AsFloat() != 42.0 {
		t.Error("String to float conversion failed")
	}

	// Test number conversion
	jv = New(42.5)
	if jv.AsInt() != 42 {
		t.Error("Float to int conversion failed")
	}

	if jv.AsString() != "42.5" {
		t.Error("Number to string conversion failed")
	}

	// Test boolean conversion
	jv = New(true)
	if jv.AsString() != "true" {
		t.Error("Bool to string conversion failed")
	}

	jv = New("true")
	if !jv.AsBool() {
		t.Error("String to bool conversion failed")
	}
}

func TestKeys(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "NYC",
	}
	jv := New(data)

	keys := jv.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	expectedKeys := []string{"name", "age", "city"}
	for _, expected := range expectedKeys {
		found := false
		for _, key := range keys {
			if key == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected key '%s' not found", expected)
		}
	}
}

func TestValues(t *testing.T) {
	arr := []interface{}{"a", "b", "c"}
	jv := New(arr)

	values := jv.Values()
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	for i, expected := range []string{"a", "b", "c"} {
		if values[i].AsString() != expected {
			t.Errorf("Expected value '%s', got '%s'", expected, values[i].AsString())
		}
	}
}

func TestLen(t *testing.T) {
	// Test array length
	arr := []interface{}{1, 2, 3, 4, 5}
	jv := New(arr)
	if jv.Len() != 5 {
		t.Errorf("Expected length 5, got %d", jv.Len())
	}

	// Test object length
	obj := map[string]interface{}{"a": 1, "b": 2}
	jv = New(obj)
	if jv.Len() != 2 {
		t.Errorf("Expected length 2, got %d", jv.Len())
	}

	// Test string length
	str := "hello"
	jv = New(str)
	if jv.Len() != 5 {
		t.Errorf("Expected length 5, got %d", jv.Len())
	}
}

func TestAppend(t *testing.T) {
	arr := []interface{}{"a", "b"}
	jv := New(arr)

	err := jv.Append("c")
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if jv.Len() != 3 {
		t.Errorf("Expected length 3 after append, got %d", jv.Len())
	}

	if jv.Get(2).AsString() != "c" {
		t.Error("Appended value not found")
	}
}

func TestExtend(t *testing.T) {
	arr := []interface{}{"a", "b"}
	jv := New(arr)

	err := jv.Extend([]interface{}{"c", "d"})
	if err != nil {
		t.Fatalf("Extend failed: %v", err)
	}

	if jv.Len() != 4 {
		t.Errorf("Expected length 4 after extend, got %d", jv.Len())
	}
}

func TestUpdate(t *testing.T) {
	obj1 := map[string]interface{}{"a": 1, "b": 2}
	obj2 := map[string]interface{}{"b": 20, "c": 3}

	jv1 := New(obj1)
	jv2 := New(obj2)

	err := jv1.Update(jv2)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if jv1.Get("a").AsInt() != 1 {
		t.Error("Original value should remain")
	}

	if jv1.Get("b").AsInt() != 20 {
		t.Error("Value should be updated")
	}

	if jv1.Get("c").AsInt() != 3 {
		t.Error("New value should be added")
	}
}

func TestClone(t *testing.T) {
	original := map[string]interface{}{
		"name": "John",
		"data": map[string]interface{}{"nested": "value"},
	}
	jv := New(original)

	cloned := jv.Clone()

	// Modify original
	jv.Set("name", "Jane")
	jv.Get("data").Set("nested", "modified")

	// Check that clone is unaffected
	if cloned.Get("name").AsString() != "John" {
		t.Error("Clone should not be affected by original changes")
	}

	if cloned.Get("data").Get("nested").AsString() != "value" {
		t.Error("Deep clone should not be affected by nested changes")
	}
}

func TestPath(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"address": map[string]interface{}{
				"street": "123 Main St",
				"city":   "NYC",
			},
		},
		"scores": []interface{}{85, 90, 78},
	}
	jv := New(data)

	// Test nested object path
	name := jv.Path("user.name")
	if name.AsString() != "John" {
		t.Errorf("Expected 'John', got '%s'", name.AsString())
	}

	// Test deeper nesting
	street := jv.Path("user.address.street")
	if street.AsString() != "123 Main St" {
		t.Errorf("Expected '123 Main St', got '%s'", street.AsString())
	}

	// Test array access in path
	firstScore := jv.Path("scores.0")
	if firstScore.AsInt() != 85 {
		t.Errorf("Expected 85, got %d", firstScore.AsInt())
	}

	// Test nonexistent path
	nonexistent := jv.Path("user.nonexistent.path")
	if !nonexistent.IsNull() {
		t.Error("Nonexistent path should return null")
	}
}

func TestSetPath(t *testing.T) {
	jv := NewObject()

	// Test setting nested path (should create intermediate objects)
	err := jv.SetPath("user.address.street", "123 Main St")
	if err != nil {
		t.Fatalf("SetPath failed: %v", err)
	}

	street := jv.Path("user.address.street")
	if street.AsString() != "123 Main St" {
		t.Error("SetPath didn't create nested structure correctly")
	}

	// Test setting array path - first create the array structure
	jv.Set("scores", []interface{}{0, 0, 0}) // Pre-create array with elements
	err = jv.SetPath("scores.0", 95)
	if err != nil {
		t.Fatalf("SetPath for array failed: %v", err)
	}

	score := jv.Path("scores.0")
	if score.AsInt() != 95 {
		t.Error("SetPath for array didn't work correctly")
	}
}

func TestNewConstructors(t *testing.T) {
	// Test NewObject
	obj := NewObject()
	if !obj.IsObject() {
		t.Error("NewObject should create an object")
	}

	if obj.Len() != 0 {
		t.Error("NewObject should create empty object")
	}

	// Test NewArray
	arr := NewArray()
	if !arr.IsArray() {
		t.Error("NewArray should create an array")
	}

	if arr.Len() != 0 {
		t.Error("NewArray should create empty array")
	}

	// Test NewArrayFrom
	items := []interface{}{"a", "b", "c"}
	arr = NewArrayFrom(items)
	if arr.Len() != 3 {
		t.Error("NewArrayFrom should create array with correct length")
	}

	// Test NewObjectFrom
	objData := map[string]interface{}{"key": "value"}
	obj = NewObjectFrom(objData)
	if obj.Get("key").AsString() != "value" {
		t.Error("NewObjectFrom should create object with correct data")
	}
}

func TestAsArray(t *testing.T) {
	arr := []interface{}{"a", "b", "c"}
	jv := New(arr)

	asArray := jv.AsArray()
	if len(asArray) != 3 {
		t.Errorf("Expected array length 3, got %d", len(asArray))
	}

	for i, expected := range []string{"a", "b", "c"} {
		if asArray[i].AsString() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, asArray[i].AsString())
		}
	}
}

func TestAsObject(t *testing.T) {
	obj := map[string]interface{}{"name": "John", "age": 30}
	jv := New(obj)

	asObject := jv.AsObject()
	if len(asObject) != 2 {
		t.Errorf("Expected object length 2, got %d", len(asObject))
	}

	if asObject["name"].AsString() != "John" {
		t.Error("AsObject conversion failed")
	}
}

func TestRaw(t *testing.T) {
	original := map[string]interface{}{"key": "value"}
	jv := New(original)

	raw := jv.Raw()
	if !reflect.DeepEqual(raw, original) {
		t.Error("Raw() should return original data")
	}
}

func TestComplexNestedOperations(t *testing.T) {
	complexJSON := `{
		"users": [
			{"id": 1, "name": "Alice", "scores": [85, 92, 78]},
			{"id": 2, "name": "Bob", "scores": [76, 89, 94]}
		],
		"metadata": {
			"total": 2,
			"active": true
		}
	}`

	jv, err := Loads(complexJSON)
	if err != nil {
		t.Fatalf("Failed to parse complex JSON: %v", err)
	}

	// Test accessing nested array elements
	firstUser := jv.Path("users.0")
	if firstUser.Get("name").AsString() != "Alice" {
		t.Error("Failed to access nested array element")
	}

	// Test accessing deeply nested values
	firstScore := jv.Path("users.0.scores.0")
	if firstScore.AsInt() != 85 {
		t.Error("Failed to access deeply nested array value")
	}

	// Test modifying nested structures
	err = jv.SetPath("users.0.scores.0", 90)
	if err != nil {
		t.Fatalf("Failed to set nested array value: %v", err)
	}

	modifiedScore := jv.Path("users.0.scores.0")
	if modifiedScore.AsInt() != 90 {
		t.Error("Nested array modification failed")
	}
}

func TestFluentQuery(t *testing.T) {
	jsonStr := `{
		"users": [
			{
				"name": "Alice",
				"profile": {
					"hair_color": "Red",
					"age": 25
				}
			},
			{
				"name": "Bob", 
				"profile": {
					"hair_color": "Brown",
					"age": 30
				}
			}
		]
	}`

	data, err := Loads(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Test fluent query access: data.Q("users", 0, "profile", "hair_color")
	hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
	if hairColor != "Red" {
		t.Errorf("Expected 'Red', got '%s'", hairColor)
	}

	// Test with multiple levels
	age := data.Q("users", 1, "profile", "age").AsInt()
	if age != 30 {
		t.Errorf("Expected 30, got %d", age)
	}

	// Test with missing path (should return null)
	missing := data.Q("users", 0, "profile", "nonexistent")
	if !missing.IsNull() {
		t.Error("Missing path should return null")
	}

	// Test with invalid index
	invalid := data.Q("users", 10, "name")
	if !invalid.IsNull() {
		t.Error("Invalid index should return null")
	}

	// Test single key access
	firstUser := data.Q("users").Q(0).Q("name").AsString()
	if firstUser != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", firstUser)
	}
}

// Benchmark tests
func BenchmarkLoads(b *testing.B) {
	jsonStr := `{"name": "John", "age": 30, "city": "NYC", "hobbies": ["reading", "swimming"]}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Loads(jsonStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	data := map[string]interface{}{"name": "John", "age": 30}
	jv := New(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jv.Get("name")
	}
}

func BenchmarkPath(b *testing.B) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"address": map[string]interface{}{
				"street": "123 Main St",
			},
		},
	}
	jv := New(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jv.Path("user.address.street")
	}
}
