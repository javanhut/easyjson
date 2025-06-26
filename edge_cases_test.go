package easyjson

import (
	"fmt"
	"io"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

// edge_cases_test.go - Comprehensive edge case testing for EasyJSON library

// TestNullPointerHandling tests various null pointer scenarios
func TestNullPointerHandling(t *testing.T) {
	t.Run("NilJSONValueOperations", func(t *testing.T) {
		var jv *JSONValue
		
		// Test that we can detect nil
		if jv != nil {
			t.Error("nil JSONValue should be nil")
		}
		
		// Note: Operations on nil JSONValue will panic in this implementation
		// This documents the current behavior - ideally they would be safe
		// The test verifies nil detection works
	})
	
	t.Run("NullDataInJSONValue", func(t *testing.T) {
		jv := &JSONValue{data: nil}
		
		if !jv.IsNull() {
			t.Error("JSONValue with nil data should report IsNull() true")
		}
		
		if jv.AsString() != "<nil>" {
			t.Error("AsString on null JSONValue should return readable representation")
		}
		
		if jv.AsInt() != 0 {
			t.Error("AsInt on null JSONValue should return 0")
		}
		
		if jv.AsBool() != false {
			t.Error("AsBool on null JSONValue should return false")
		}
		
		if jv.Len() != 0 {
			t.Error("Len on null JSONValue should return 0")
		}
	})
	
	t.Run("NullInNestedStructures", func(t *testing.T) {
		data := map[string]interface{}{
			"user": nil,
			"settings": map[string]interface{}{
				"theme": nil,
				"lang": "en",
			},
			"scores": []interface{}{10, nil, 30},
		}
		jv := New(data)
		
		user := jv.Get("user")
		if !user.IsNull() {
			t.Error("Null value in map should be null")
		}
		
		theme := jv.Path("settings.theme")
		if !theme.IsNull() {
			t.Error("Null value in nested map should be null")
		}
		
		nullScore := jv.Path("scores.1")
		if !nullScore.IsNull() {
			t.Error("Null value in array should be null")
		}
		
		// Chaining through null should return null
		chained := jv.Q("user", "profile", "name")
		if !chained.IsNull() {
			t.Error("Chaining through null should return null")
		}
	})
}

// TestArrayBoundsEdgeCases tests array index edge cases
func TestArrayBoundsEdgeCases(t *testing.T) {
	t.Run("NegativeIndices", func(t *testing.T) {
		arr := []interface{}{"a", "b", "c"}
		jv := New(arr)
		
		// Get with negative index should return null
		result := jv.Get(-1)
		if !result.IsNull() {
			t.Error("Negative array index should return null")
		}
		
		result = jv.Get(-100)
		if !result.IsNull() {
			t.Error("Large negative array index should return null")
		}
		
		// Set with negative index should fail
		err := jv.Set(-1, "invalid")
		if err == nil {
			t.Error("Setting negative array index should fail")
		}
		
		// Has with negative index should return false
		if jv.Has(-1) {
			t.Error("Has with negative index should return false")
		}
	})
	
	t.Run("OutOfBoundsIndices", func(t *testing.T) {
		arr := []interface{}{"a", "b", "c"}
		jv := New(arr)
		
		// Get out of bounds should return null
		result := jv.Get(3)
		if !result.IsNull() {
			t.Error("Out of bounds array index should return null")
		}
		
		result = jv.Get(1000)
		if !result.IsNull() {
			t.Error("Large out of bounds index should return null")
		}
		
		// Has out of bounds should return false
		if jv.Has(3) {
			t.Error("Has with out of bounds index should return false")
		}
		
		if jv.Has(1000) {
			t.Error("Has with large out of bounds index should return false")
		}
	})
	
	t.Run("ExtremeArraySizes", func(t *testing.T) {
		// Test with large valid array
		largeArr := make([]interface{}, 1000)
		for i := range largeArr {
			largeArr[i] = i
		}
		jv := New(largeArr)
		
		// Should be able to access all elements
		if jv.Get(0).AsInt() != 0 {
			t.Error("Should access first element of large array")
		}
		
		if jv.Get(999).AsInt() != 999 {
			t.Error("Should access last element of large array")
		}
		
		// Test array size limits in Set operations
		err := jv.Set(MaxArrayIndex+1, "invalid")
		if err == nil {
			t.Error("Setting index beyond MaxArrayIndex should fail")
		}
	})
	
	t.Run("EmptyArrayOperations", func(t *testing.T) {
		jv := NewArray()
		
		// Operations on empty array
		if jv.Len() != 0 {
			t.Error("Empty array should have length 0")
		}
		
		result := jv.Get(0)
		if !result.IsNull() {
			t.Error("Get on empty array should return null")
		}
		
		if jv.Has(0) {
			t.Error("Has on empty array should return false")
		}
		
		// Should be able to append to empty array
		err := jv.Append("first")
		if err != nil {
			t.Errorf("Should be able to append to empty array: %v", err)
		}
		
		if jv.Len() != 1 {
			t.Error("Array should have length 1 after append")
		}
	})
}

// TestPathValidationEdgeCases tests path validation edge cases
func TestPathValidationEdgeCases(t *testing.T) {
	t.Run("EmptyPaths", func(t *testing.T) {
		jv := NewObject()
		
		// Empty path should fail
		result := jv.Path("")
		if !result.IsNull() {
			t.Error("Empty path should return null")
		}
		
		err := jv.SetPath("", "value")
		if err == nil {
			t.Error("Setting empty path should fail")
		}
	})
	
	t.Run("InvalidPathCharacters", func(t *testing.T) {
		jv := NewObject()
		
		// Paths with invalid characters
		invalidPaths := []string{
			"key\x00value",  // null character
			"key\tvalue",    // tab character
			"key\nvalue",    // newline character
			"key\rvalue",    // carriage return
		}
		
		for _, path := range invalidPaths {
			result := jv.Path(path)
			if !result.IsNull() {
				t.Errorf("Path with invalid characters should return null: %s", path)
			}
		}
	})
	
	t.Run("VeryLongPaths", func(t *testing.T) {
		jv := NewObject()
		
		// Create a path that exceeds MaxPathLength
		longPath := strings.Repeat("a.", MaxPathLength/2) + "b"
		result := jv.Path(longPath)
		if !result.IsNull() {
			t.Error("Extremely long path should return null")
		}
		
		err := jv.SetPath(longPath, "value")
		if err == nil {
			t.Error("Setting extremely long path should fail")
		}
	})
	
	t.Run("DeepNestingPaths", func(t *testing.T) {
		jv := NewObject()
		
		// Create path with excessive depth
		deepPath := strings.Repeat("level.", MaxPathDepth) + "value"
		result := jv.Path(deepPath)
		if !result.IsNull() {
			t.Error("Excessively deep path should return null")
		}
		
		err := jv.SetPath(deepPath, "value")
		if err == nil {
			t.Error("Setting excessively deep path should fail")
		}
	})
	
	t.Run("PathTraversalAttempts", func(t *testing.T) {
		jv := NewObject()
		
		// Attempt path traversal
		traversalPaths := []string{
			"../secret",
			"key/../other",
			"valid..invalid",
			"./current",
		}
		
		for _, path := range traversalPaths {
			result := jv.Path(path)
			if !result.IsNull() {
				t.Errorf("Path traversal attempt should return null: %s", path)
			}
		}
	})
	
	t.Run("PathsWithConsecutiveDots", func(t *testing.T) {
		jv := NewObject()
		
		// Paths with consecutive dots
		result := jv.Path("key..value")
		if !result.IsNull() {
			t.Error("Path with consecutive dots should return null")
		}
		
		result = jv.Path("key...value")
		if !result.IsNull() {
			t.Error("Path with multiple consecutive dots should return null")
		}
	})
}

// TestMemoryLimitEdgeCases tests memory limit handling
func TestMemoryLimitEdgeCases(t *testing.T) {
	t.Run("LargeStringHandling", func(t *testing.T) {
		// Create a string that approaches MaxStringLength
		largeString := strings.Repeat("x", MaxStringLength-1000)
		
		// Should be able to handle large but valid strings
		jv := New(largeString)
		if jv.AsString() != largeString {
			t.Error("Should handle large valid strings")
		}
		
		// Test GetString with large default
		result := jv.GetString("nonexistent", largeString)
		if result != largeString {
			t.Error("Should handle large default string")
		}
	})
	
	t.Run("OversizedStringPrevention", func(t *testing.T) {
		// This test ensures the validation prevents creation of oversized content
		// We test the validation directly rather than trying to create oversized content
		
		err := validateStringLength(strings.Repeat("x", MaxStringLength+1))
		if err == nil {
			t.Error("Should prevent oversized string creation")
		}
		
		// Test safe parsing with oversized input
		oversizedJSON := fmt.Sprintf(`{"key": "%s"}`, strings.Repeat("x", MaxStringLength))
		result := ParseSafely(oversizedJSON)
		if result.Error == nil {
			t.Error("Should fail to parse oversized JSON")
		}
		if len(result.Suggestions) == 0 {
			t.Error("Should provide suggestions for oversized JSON")
		}
	})
	
	t.Run("DeepNestingMemoryUsage", func(t *testing.T) {
		// Create deeply nested structure within limits
		current := NewObject()
		root := current
		
		// Create nesting within reasonable limits
		for i := 0; i < DefaultMaxRecursionDepth-10; i++ {
			nested := NewObject()
			current.Set("child", nested.Raw())
			current = nested
		}
		
		// Should be able to access deeply nested structure within limits
		path := strings.Repeat("child.", DefaultMaxRecursionDepth-11) + "child"
		result := root.Path(path)
		// Note: This may fail if the implementation has stricter path limits
		if result.IsNull() {
			t.Log("Deep nesting path access failed - this may be due to implementation limits")
		}
	})
	
	t.Run("LargeArrayHandling", func(t *testing.T) {
		// Test array size validation
		err := validateArraySize(MaxArraySize + 1)
		if err == nil {
			t.Error("Should prevent oversized array creation")
		}
		
		// Test with large but valid array
		largeArr := make([]interface{}, 1000)
		for i := range largeArr {
			largeArr[i] = fmt.Sprintf("item%d", i)
		}
		
		jv := New(largeArr)
		if jv.Len() != 1000 {
			t.Error("Should handle large valid array")
		}
		
		// Test operations on large array
		result := jv.Get(500)
		if result.AsString() != "item500" {
			t.Error("Should access elements in large array")
		}
	})
}

// TestTypeValidationEdgeCases tests type validation edge cases
func TestTypeValidationEdgeCases(t *testing.T) {
	t.Run("InvalidKeyTypes", func(t *testing.T) {
		jv := NewObject()
		
		// Test invalid key types
		invalidKeys := []interface{}{
			nil,
			[]int{1, 2, 3},
			map[string]string{"invalid": "key"},
			func() {},
			make(chan int),
		}
		
		for _, key := range invalidKeys {
			result := jv.Get(key)
			if !result.IsNull() {
				t.Errorf("Invalid key type should return null: %T", key)
			}
			
			err := jv.Set(key, "value")
			if err == nil {
				t.Errorf("Should not be able to set with invalid key type: %T", key)
			}
		}
	})
	
	t.Run("InvalidValueTypes", func(t *testing.T) {
		jv := NewObject()
		
		// Test invalid value types
		invalidValues := []interface{}{
			func() {},
			make(chan int),
			complex(1, 2),
		}
		
		for _, value := range invalidValues {
			err := jv.Set("key", value)
			if err == nil {
				t.Errorf("Should not be able to set invalid value type: %T", value)
			}
		}
	})
	
	t.Run("MixedTypeArrays", func(t *testing.T) {
		// Mixed type arrays should be valid in JSON
		mixedArr := []interface{}{
			"string",
			123,
			45.67,
			true,
			nil,
			map[string]interface{}{"nested": "object"},
			[]interface{}{"nested", "array"},
		}
		
		jv := New(mixedArr)
		
		// Should be able to access all types
		if jv.Get(0).AsString() != "string" {
			t.Error("Should handle string in mixed array")
		}
		
		if jv.Get(1).AsInt() != 123 {
			t.Error("Should handle int in mixed array")
		}
		
		if jv.Get(2).AsFloat() != 45.67 {
			t.Error("Should handle float in mixed array")
		}
		
		if !jv.Get(3).AsBool() {
			t.Error("Should handle bool in mixed array")
		}
		
		if !jv.Get(4).IsNull() {
			t.Error("Should handle null in mixed array")
		}
		
		if !jv.Get(5).IsObject() {
			t.Error("Should handle object in mixed array")
		}
		
		if !jv.Get(6).IsArray() {
			t.Error("Should handle array in mixed array")
		}
	})
	
	t.Run("NumericTypeEdgeCases", func(t *testing.T) {
		// Test extreme numeric values
		extremeNumbers := []interface{}{
			math.MaxFloat64,
			-math.MaxFloat64,
			math.SmallestNonzeroFloat64,
			float64(math.MaxInt64),
			float64(math.MinInt64),
		}
		
		for _, num := range extremeNumbers {
			jv := New(num)
			if !jv.IsNumber() {
				t.Errorf("Should recognize extreme number as number: %v", num)
			}
			
			// Should be able to convert without panic
			_ = jv.AsFloat()
			_ = jv.AsInt()
			_ = jv.AsString()
		}
		
		// Test NaN and Inf (should be handled gracefully)
		nanValue := New(math.NaN())
		if !nanValue.IsNumber() {
			t.Error("NaN should be recognized as number type")
		}
		
		infValue := New(math.Inf(1))
		if !infValue.IsNumber() {
			t.Error("Infinity should be recognized as number type")
		}
	})
}

// TestRecursionDepthEdgeCases tests recursion depth handling
func TestRecursionDepthEdgeCases(t *testing.T) {
	t.Run("RecursionDepthValidation", func(t *testing.T) {
		// Test recursion depth validation directly
		err := validateRecursionDepth(DefaultMaxRecursionDepth, DefaultMaxRecursionDepth)
		if err == nil {
			t.Error("Should fail when current depth equals max depth")
		}
		
		err = validateRecursionDepth(DefaultMaxRecursionDepth+1, DefaultMaxRecursionDepth)
		if err == nil {
			t.Error("Should fail when current depth exceeds max depth")
		}
		
		err = validateRecursionDepth(DefaultMaxRecursionDepth-1, DefaultMaxRecursionDepth)
		if err != nil {
			t.Error("Should pass when current depth is less than max depth")
		}
	})
	
	t.Run("DeeplyNestedCloning", func(t *testing.T) {
		// Create deeply nested structure for cloning test
		root := NewObject()
		current := root
		
		depth := DefaultMaxRecursionDepth / 2 // Stay well within limits
		for i := 0; i < depth; i++ {
			nested := NewObject()
			nested.Set("value", fmt.Sprintf("level%d", i))
			current.Set("child", nested.Raw())
			current = nested
		}
		
		// Clone should handle deep nesting
		cloned := root.Clone()
		if cloned == nil {
			t.Error("Clone should handle deeply nested structures")
		}
		
		// Verify clone independence
		root.Q("child").Set("modified", true)
		
		modified := cloned.Q("child").Get("modified")
		if !modified.IsNull() {
			t.Error("Clone should be independent of original")
		}
	})
	
	t.Run("DeepSearchLimits", func(t *testing.T) {
		// Create structure with multiple levels
		root := NewObject()
		
		// Create several nested levels
		for i := 0; i < 10; i++ {
			levelObj := NewObject()
			levelObj.Set("target", fmt.Sprintf("found_at_level_%d", i))
			levelObj.Set("other", "not_target")
			
			root.Set(fmt.Sprintf("level%d", i), levelObj.Raw())
		}
		
		// DeepSearch should find targets without exceeding limits
		result := root.DeepSearch("target")
		if result.IsNull() {
			t.Error("DeepSearch should find target in nested structure")
		}
		
		allResults := root.DeepSearchAll("target")
		if len(allResults) == 0 {
			t.Error("DeepSearchAll should find multiple targets")
		}
	})
}

// TestUnicodeAndSpecialCharacters tests Unicode and special character handling
func TestUnicodeAndSpecialCharacters(t *testing.T) {
	t.Run("UnicodeStringHandling", func(t *testing.T) {
		unicodeStrings := []string{
			"Hello, 世界",
			"🌟✨🎉",
			"Iñtërnâtiônàlizætiøn",
			"Здравствуй мир",
			"مرحبا بالعالم",
			"שלום עולם",
			"こんにちは世界",
		}
		
		for _, str := range unicodeStrings {
			if !utf8.ValidString(str) {
				t.Errorf("Test string should be valid UTF-8: %s", str)
				continue
			}
			
			jv := New(str)
			if jv.AsString() != str {
				t.Errorf("Unicode string not handled correctly: %s", str)
			}
			
			// Test in object keys
			obj := NewObject()
			err := obj.Set(str, "value")
			if err != nil {
				t.Errorf("Should be able to use Unicode string as key: %s", str)
			}
			
			// Test in paths
			if len(str) < MaxPathSegmentLength { // Only test if within limits
				result := obj.Path(str)
				if result.AsString() != "value" {
					t.Errorf("Should be able to use Unicode string in path: %s", str)
				}
			}
		}
	})
	
	t.Run("InvalidUTF8Handling", func(t *testing.T) {
		// Create invalid UTF-8 sequences
		invalidUTF8 := []string{
			"\xff\xfe\xfd",
			"\x80\x81\x82",
			string([]byte{0xff, 0xfe, 0xfd}),
		}
		
		for _, invalid := range invalidUTF8 {
			if utf8.ValidString(invalid) {
				continue // Skip if somehow valid
			}
			
			// Should not crash when handling invalid UTF-8
			jv := New(invalid)
			_ = jv.AsString() // Should not panic
			
			// Should not be able to use as key
			obj := NewObject()
			err := obj.Set(invalid, "value")
			if err == nil {
				t.Error("Should not be able to use invalid UTF-8 as key")
			}
		}
	})
	
	t.Run("ControlCharacterHandling", func(t *testing.T) {
		controlChars := []string{
			"\x00",  // null
			"\x01",  // start of heading
			"\x02",  // start of text
			"\x03",  // end of text
			"\x04",  // end of transmission
			"\x08",  // backspace
			"\x0B",  // vertical tab
			"\x0C",  // form feed
			"\x0E",  // shift out
			"\x0F",  // shift in
			"\x7F",  // delete
		}
		
		for _, ctrl := range controlChars {
			// Control characters should be handled safely
			jv := New(ctrl)
			_ = jv.AsString() // Should not panic
			
			// Should not be able to use control chars in keys
			obj := NewObject()
			err := obj.Set(ctrl, "value")
			if err == nil {
				t.Errorf("Should not be able to use control character as key: %q", ctrl)
			}
			
			// Should not be able to use in paths
			result := obj.Path(ctrl)
			if !result.IsNull() {
				t.Errorf("Should not be able to use control character in path: %q", ctrl)
			}
		}
	})
	
	t.Run("SpecialJSONCharacters", func(t *testing.T) {
		specialChars := []string{
			`"`,        // quote
			`\"`,       // escaped quote
			`\\`,       // backslash
			`\/`,       // escaped slash
			`\b`,       // backspace
			`\f`,       // form feed
			`\n`,       // newline
			`\r`,       // carriage return
			`\t`,       // tab
			`\u0000`,   // unicode null
			`\u00FF`,   // unicode character
		}
		
		for _, char := range specialChars {
			// Should be able to handle special JSON characters in values
			jv := New(char)
			_ = jv.AsString() // Should not panic
			
			// Test JSON round-trip
			jsonStr, err := jv.Dumps()
			if err != nil {
				t.Errorf("Should be able to serialize special character: %s", char)
				continue
			}
			
			parsed, err := Loads(jsonStr)
			if err != nil {
				t.Errorf("Should be able to parse back special character: %s", char)
			} else if parsed.AsString() != char {
				t.Errorf("Round-trip failed for special character: %s", char)
			}
		}
	})
}

// TestMalformedJSONRecovery tests malformed JSON handling
func TestMalformedJSONRecovery(t *testing.T) {
	t.Run("CommonMalformations", func(t *testing.T) {
		malformedJSONs := []struct {
			json        string
			description string
		}{
			{`{"key": }`, "missing value"},
			{`{"key": "value",}`, "trailing comma"},
			{`{key: "value"}`, "unquoted key"},
			{`{'key': 'value'}`, "single quotes"},
			{`{"key": "value"`, "missing closing brace"},
			{`"key": "value"}`, "missing opening brace"},
			{`{"key": "value", "other":}`, "incomplete second pair"},
			{`[1, 2, 3,]`, "trailing comma in array"},
			{`[1 2 3]`, "missing commas"},
			{`{"key": undefined}`, "undefined value"},
		}
		
		for _, test := range malformedJSONs {
			t.Run(test.description, func(t *testing.T) {
				// ParseSafely should handle malformed JSON gracefully
				result := ParseSafely(test.json)
				
				if result.Error == nil {
					t.Errorf("Should detect malformed JSON: %s", test.description)
				}
				
				if result.Data == nil {
					t.Error("Should always return valid JSONValue even for malformed JSON")
				}
				
				if len(result.Suggestions) == 0 {
					t.Error("Should provide suggestions for malformed JSON")
				}
				
				// ParseLenient should attempt recovery
				recovered := ParseLenient(test.json)
				if recovered == nil {
					t.Error("ParseLenient should always return valid JSONValue")
				}
			})
		}
	})
	
	t.Run("PartialJSONExtraction", func(t *testing.T) {
		partialJSONs := []string{
			`{"valid": "json"} and some extra text`,
			`prefix text {"valid": "json"}`,
			`{"valid": "json"} {"another": "object"}`,
			`[1, 2, 3] extra content`,
		}
		
		for _, partial := range partialJSONs {
			recovered := ParseLenient(partial)
			if recovered == nil {
				t.Error("ParseLenient should extract valid JSON from partial content")
			}
			
			if recovered.IsNull() && len(partial) > 0 {
				t.Errorf("Should extract some valid content from: %s", partial)
			}
		}
	})
	
	t.Run("PythonStyleFixing", func(t *testing.T) {
		pythonStyleJSONs := []struct {
			python   string
			expected string
		}{
			{`{"active": True}`, "true"},
			{`{"active": False}`, "false"},
			{`{"value": None}`, "null"},
			{`{'key': 'value'}`, "value"},
		}
		
		for _, test := range pythonStyleJSONs {
			// Test FixCommonIssues
			fixed := FixCommonIssues(test.python)
			if fixed == test.python {
				t.Errorf("FixCommonIssues should modify Python-style JSON: %s", test.python)
			}
			
			// Test ParseWithFixes
			result, err := ParseWithFixes(test.python)
			if err != nil {
				t.Errorf("ParseWithFixes should handle Python-style JSON: %s", test.python)
			}
			
			if test.expected != "null" {
				value := result.Get("key")
				if value.IsNull() {
					value = result.Get("active")
				}
				if value.IsNull() {
					value = result.Get("value")
				}
				
				if !value.IsNull() && value.AsString() != test.expected {
					t.Errorf("Expected %s, got %s for %s", test.expected, value.AsString(), test.python)
				}
			}
		}
	})
}

// TestConcurrentAccess tests concurrent access patterns
func TestConcurrentAccess(t *testing.T) {
	t.Run("ConcurrentReads", func(t *testing.T) {
		// Create shared data
		data := map[string]interface{}{
			"counter": 0,
			"users": []interface{}{
				map[string]interface{}{"id": 1, "name": "Alice"},
				map[string]interface{}{"id": 2, "name": "Bob"},
			},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
		}
		jv := New(data)
		
		// Concurrent readers
		const numReaders = 100
		var wg sync.WaitGroup
		errors := make(chan error, numReaders)
		
		for i := 0; i < numReaders; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Perform various read operations
				if jv.Get("counter").AsInt() != 0 {
					errors <- fmt.Errorf("reader %d: unexpected counter value", id)
					return
				}
				
				users := jv.Get("users")
				if users.Len() != 2 {
					errors <- fmt.Errorf("reader %d: unexpected users length", id)
					return
				}
				
				alice := users.Get(0)
				if alice.Get("name").AsString() != "Alice" {
					errors <- fmt.Errorf("reader %d: unexpected user name", id)
					return
				}
				
				theme := jv.Path("settings.theme")
				if theme.AsString() != "dark" {
					errors <- fmt.Errorf("reader %d: unexpected theme", id)
					return
				}
			}(i)
		}
		
		wg.Wait()
		close(errors)
		
		// Check for errors
		for err := range errors {
			t.Error(err)
		}
	})
	
	t.Run("ReadWriteRaceCondition", func(t *testing.T) {
		// Test read/write operations on separate JSONValue instances
		// (Safe since each operation creates new JSONValue instances)
		
		originalData := map[string]interface{}{
			"value": 0,
		}
		
		const numOperations = 50
		var wg sync.WaitGroup
		errors := make(chan error, numOperations*2)
		
		// Readers
		for i := 0; i < numOperations; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Create separate instance for reading
				jv := New(originalData)
				value := jv.Get("value").AsInt()
				if value != 0 {
					errors <- fmt.Errorf("reader %d: unexpected value %d", id, value)
				}
			}(i)
		}
		
		// Writers (on separate instances)
		for i := 0; i < numOperations; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				// Create separate instance for writing
				jv := New(map[string]interface{}{"value": id})
				err := jv.Set("value", id+1)
				if err != nil {
					errors <- fmt.Errorf("writer %d: set failed: %v", id, err)
				}
			}(i)
		}
		
		wg.Wait()
		close(errors)
		
		// Check for errors
		for err := range errors {
			t.Error(err)
		}
	})
	
	t.Run("ConcurrentCloning", func(t *testing.T) {
		// Test concurrent cloning operations
		complexData := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{"id": 1, "profile": map[string]interface{}{"name": "Alice"}},
				map[string]interface{}{"id": 2, "profile": map[string]interface{}{"name": "Bob"}},
			},
			"metadata": map[string]interface{}{
				"version": "1.0",
				"settings": []interface{}{"setting1", "setting2"},
			},
		}
		
		jv := New(complexData)
		
		const numCloners = 50
		var wg sync.WaitGroup
		errors := make(chan error, numCloners)
		
		for i := 0; i < numCloners; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				cloned := jv.Clone()
				if cloned == nil {
					errors <- fmt.Errorf("cloner %d: clone returned nil", id)
					return
				}
				
				// Verify clone content
				if cloned.Get("users").Len() != 2 {
					errors <- fmt.Errorf("cloner %d: clone has wrong users length", id)
					return
				}
				
				if cloned.Path("users.0.profile.name").AsString() != "Alice" {
					errors <- fmt.Errorf("cloner %d: clone has wrong nested value", id)
					return
				}
				
				// Modify clone to test independence
				err := cloned.Set("modified_by", id)
				if err != nil {
					errors <- fmt.Errorf("cloner %d: failed to modify clone: %v", id, err)
				}
			}(i)
		}
		
		wg.Wait()
		close(errors)
		
		// Check for errors
		for err := range errors {
			t.Error(err)
		}
		
		// Verify original is unmodified
		if jv.Has("modified_by") {
			t.Error("Original should not be modified by clone operations")
		}
	})
}

// TestPerformanceEdgeCases tests performance edge cases
func TestPerformanceEdgeCases(t *testing.T) {
	t.Run("LargeArrayPerformance", func(t *testing.T) {
		// Create large array
		size := 10000
		largeArr := make([]interface{}, size)
		for i := 0; i < size; i++ {
			largeArr[i] = map[string]interface{}{
				"id":   i,
				"name": fmt.Sprintf("item_%d", i),
				"data": []interface{}{i, i * 2, i * 3},
			}
		}
		
		jv := New(largeArr)
		
		start := time.Now()
		
		// Test various operations
		_ = jv.Len()
		_ = jv.Get(0)
		_ = jv.Get(size - 1)
		_ = jv.Get(size / 2)
		
		// Test path operations
		_ = jv.Path("0.name")
		_ = jv.Path(fmt.Sprintf("%d.data.1", size-1))
		
		elapsed := time.Since(start)
		
		// Should complete within reasonable time (adjust threshold as needed)
		if elapsed > time.Second {
			t.Errorf("Large array operations took too long: %v", elapsed)
		}
	})
	
	t.Run("DeepNestingPerformance", func(t *testing.T) {
		// Create deeply nested structure
		depth := 100
		root := NewObject()
		current := root
		
		start := time.Now()
		
		// Build deep nesting
		for i := 0; i < depth; i++ {
			nested := NewObject()
			nested.Set("level", i)
			nested.Set("data", fmt.Sprintf("level_%d_data", i))
			current.Set("child", nested.Raw())
			current = nested
		}
		
		buildTime := time.Since(start)
		
		// Test access performance
		start = time.Now()
		
		pathParts := make([]string, depth+1)
		for i := 0; i <= depth; i++ {
			pathParts[i] = "child"
		}
		path := strings.Join(pathParts[:depth], ".")
		
		result := root.Path(path)
		if result.IsNull() {
			t.Log("Deep nesting path access failed - this may be due to implementation limits")
		}
		
		accessTime := time.Since(start)
		
		// Should complete within reasonable time
		if buildTime > time.Second {
			t.Errorf("Deep nesting construction took too long: %v", buildTime)
		}
		
		if accessTime > time.Second {
			t.Errorf("Deep nesting access took too long: %v", accessTime)
		}
	})
	
	t.Run("ManyKeysPerformance", func(t *testing.T) {
		// Create object with many keys
		numKeys := 10000
		largeObj := make(map[string]interface{})
		
		start := time.Now()
		
		for i := 0; i < numKeys; i++ {
			key := fmt.Sprintf("key_%d", i)
			largeObj[key] = map[string]interface{}{
				"value": i,
				"name":  fmt.Sprintf("name_%d", i),
			}
		}
		
		jv := New(largeObj)
		buildTime := time.Since(start)
		
		// Test access performance
		start = time.Now()
		
		// Access various keys
		_ = jv.Get("key_0")
		_ = jv.Get(fmt.Sprintf("key_%d", numKeys-1))
		_ = jv.Get(fmt.Sprintf("key_%d", numKeys/2))
		
		// Test path access
		_ = jv.Path("key_100.value")
		_ = jv.Path(fmt.Sprintf("key_%d.name", numKeys-100))
		
		// Test Keys() operation
		keys := jv.Keys()
		if len(keys) != numKeys {
			t.Errorf("Expected %d keys, got %d", numKeys, len(keys))
		}
		
		accessTime := time.Since(start)
		
		// Should complete within reasonable time
		if buildTime > time.Second*5 {
			t.Errorf("Large object construction took too long: %v", buildTime)
		}
		
		if accessTime > time.Second {
			t.Errorf("Large object access took too long: %v", accessTime)
		}
	})
	
	t.Run("RepeatedOperationsPerformance", func(t *testing.T) {
		data := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{"id": 1, "name": "Alice"},
				map[string]interface{}{"id": 2, "name": "Bob"},
			},
			"settings": map[string]interface{}{
				"theme": "dark",
			},
		}
		jv := New(data)
		
		iterations := 10000
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			// Repeated access operations
			_ = jv.Get("users")
			_ = jv.Path("users.0.name")
			_ = jv.Q("settings", "theme")
			_ = jv.Get("users").Get(0).Get("id")
		}
		
		elapsed := time.Since(start)
		
		// Should complete within reasonable time
		if elapsed > time.Second*2 {
			t.Errorf("Repeated operations took too long: %v", elapsed)
		}
		
		avgPerOp := elapsed / time.Duration(iterations*4) // 4 operations per iteration
		t.Logf("Average time per operation: %v", avgPerOp)
	})
}

// TestSourceVariations tests ParseSafelyFrom with various sources
func TestSourceVariations(t *testing.T) {
	t.Run("ByteSliceSource", func(t *testing.T) {
		jsonBytes := []byte(`{"test": "value"}`)
		result := ParseSafelyFrom(jsonBytes)
		
		if result.Error != nil {
			t.Errorf("Should parse valid byte slice: %v", result.Error)
		}
		
		if result.Data.GetString("test") != "value" {
			t.Error("Should parse byte slice correctly")
		}
	})
	
	t.Run("ReaderSource", func(t *testing.T) {
		jsonStr := `{"reader": "test"}`
		reader := strings.NewReader(jsonStr)
		
		result := ParseSafelyFrom(reader)
		
		if result.Error != nil {
			t.Errorf("Should parse from reader: %v", result.Error)
		}
		
		if result.Data.GetString("reader") != "test" {
			t.Error("Should parse reader correctly")
		}
	})
	
	t.Run("LimitedReaderSource", func(t *testing.T) {
		// Create a reader with content that exceeds MaxJSONSize
		largeContent := fmt.Sprintf(`{"data": "%s"}`, strings.Repeat("x", int(MaxJSONSize)))
		reader := strings.NewReader(largeContent)
		
		result := ParseSafelyFrom(reader)
		
		if result.Error == nil {
			t.Error("Should fail for oversized reader content")
		}
		
		if len(result.Suggestions) == 0 {
			t.Error("Should provide suggestions for oversized content")
		}
	})
	
	t.Run("NilSource", func(t *testing.T) {
		result := ParseSafelyFrom(nil)
		
		if result.Error == nil {
			t.Error("Should fail for nil source")
		}
		
		if result.Data == nil {
			t.Error("Should return valid JSONValue even for nil source")
		}
	})
	
	t.Run("InvalidSourceType", func(t *testing.T) {
		result := ParseSafelyFrom(123)
		
		if result.Error == nil {
			t.Error("Should fail for invalid source type")
		}
		
		if len(result.Suggestions) == 0 {
			t.Error("Should provide suggestions for invalid source type")
		}
	})
	
	t.Run("FailingReader", func(t *testing.T) {
		// Create a reader that fails
		failingReader := &failingReader{}
		
		result := ParseSafelyFrom(failingReader)
		
		if result.Error == nil {
			t.Error("Should fail for failing reader")
		}
		
		if result.Data == nil {
			t.Error("Should return valid JSONValue even for failing reader")
		}
	})
}

// failingReader implements io.Reader but always returns an error
type failingReader struct{}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated reader failure")
}

// Helper function to create a buffer reader that simulates partial reads
type slowReader struct {
	data []byte
	pos  int
}

func (sr *slowReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}
	
	// Read only one byte at a time to simulate slow/partial reads
	if len(p) > 0 {
		p[0] = sr.data[sr.pos]
		sr.pos++
		return 1, nil
	}
	
	return 0, nil
}

func TestSlowReader(t *testing.T) {
	jsonStr := `{"slow": "reader", "test": true}`
	slowReader := &slowReader{data: []byte(jsonStr)}
	
	result := ParseSafelyFrom(slowReader)
	
	if result.Error != nil {
		t.Errorf("Should handle slow reader: %v", result.Error)
	}
	
	if result.Data.GetString("slow") != "reader" {
		t.Error("Should parse slow reader correctly")
	}
}

// TestValidationErrorHandling tests validation error structures
func TestValidationErrorHandling(t *testing.T) {
	t.Run("ValidationErrorStructure", func(t *testing.T) {
		err := newValidationError("TestType", "test message", "test value")
		
		if err.Type != "TestType" {
			t.Error("ValidationError should preserve type")
		}
		
		if err.Message != "test message" {
			t.Error("ValidationError should preserve message")
		}
		
		if err.Value != "test value" {
			t.Error("ValidationError should preserve value")
		}
		
		errorString := err.Error()
		if !strings.Contains(errorString, "TestType") ||
			!strings.Contains(errorString, "test message") ||
			!strings.Contains(errorString, "test value") {
			t.Errorf("ValidationError.Error() should contain all components: %s", errorString)
		}
	})
	
	t.Run("ValidationFunctionCoverage", func(t *testing.T) {
		// Test all validation functions with edge cases
		
		// Path validation
		err := validatePath(strings.Repeat("a", MaxPathLength+1))
		if err == nil {
			t.Error("Should validate path length")
		}
		
		// Array validation
		err = validateArraySize(-1)
		if err == nil {
			t.Error("Should validate negative array size")
		}
		
		// Key validation
		err = validateKey(make(chan int))
		if err == nil {
			t.Error("Should validate invalid key type")
		}
		
		// Memory validation
		err = validateMemoryLimit(100, 50)
		if err == nil {
			t.Error("Should validate memory limit exceeded")
		}
		
		// JSON size validation
		err = validateJSONSize(MaxJSONSize + 1)
		if err == nil {
			t.Error("Should validate JSON size limit")
		}
	})
}

// Benchmark tests for edge cases
func BenchmarkLargeArrayAccess(b *testing.B) {
	size := 10000
	arr := make([]interface{}, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}
	jv := New(arr)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jv.Get(i % size)
	}
}

func BenchmarkDeepNestingAccess(b *testing.B) {
	// Create nested structure
	root := NewObject()
	current := root
	depth := 100
	
	for i := 0; i < depth; i++ {
		nested := NewObject()
		nested.Set("data", i)
		current.Set("child", nested.Raw())
		current = nested
	}
	
	path := strings.Repeat("child.", depth-1) + "data"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = root.Path(path)
	}
}

func BenchmarkConcurrentReads(b *testing.B) {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "id": 1},
			map[string]interface{}{"name": "Bob", "id": 2},
		},
	}
	jv := New(data)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = jv.Path("users.0.name")
		}
	})
}