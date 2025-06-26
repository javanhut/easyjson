package easyjson

import (
	"fmt"
	"log"
)

// ExampleValidationUsage demonstrates how to use the validation functions
func ExampleValidationUsage() {
	fmt.Println("=== EasyJSON Validation Examples ===")

	// Example 1: Path validation
	fmt.Println("\n1. Path Validation:")
	validPaths := []string{"user.name", "items.0.details", "data.nested.value"}
	invalidPaths := []string{"", "user..name", "user.\x00.name", "very.long.path.that.exceeds.maximum.depth.limit.and.should.fail.validation.because.it.has.too.many.segments.and.is.really.really.really.really.really.really.really.really.really.really.really.really.really.really.really.long"}

	for _, path := range validPaths {
		if err := validatePath(path); err != nil {
			fmt.Printf("  ❌ Path '%s' validation failed: %v\n", path, err)
		} else {
			fmt.Printf("  ✅ Path '%s' is valid\n", path)
		}
	}

	for _, path := range invalidPaths {
		if err := validatePath(path); err != nil {
			fmt.Printf("  ❌ Path '%s' validation failed: %v\n", path, err)
		} else {
			fmt.Printf("  ✅ Path '%s' is valid\n", path)
		}
	}

	// Example 2: Array index validation
	fmt.Println("\n2. Array Index Validation:")
	arrayLen := 10
	testIndices := []int{0, 5, 9, 10, -1, 1000000}

	for _, index := range testIndices {
		if err := validateArrayIndex(index, arrayLen); err != nil {
			fmt.Printf("  ❌ Index %d (array len %d): %v\n", index, arrayLen, err)
		} else {
			fmt.Printf("  ✅ Index %d (array len %d) is valid\n", index, arrayLen)
		}
	}

	// Example 3: Key validation
	fmt.Println("\n3. Key Validation:")
	validKeys := []interface{}{"username", "email", 42, 0}
	invalidKeys := []interface{}{"", -1, 3.14, nil, []int{1, 2, 3}}

	for _, key := range validKeys {
		if err := validateKey(key); err != nil {
			fmt.Printf("  ❌ Key %v (%T): %v\n", key, key, err)
		} else {
			fmt.Printf("  ✅ Key %v (%T) is valid\n", key, key)
		}
	}

	for _, key := range invalidKeys {
		if err := validateKey(key); err != nil {
			fmt.Printf("  ❌ Key %v (%T): %v\n", key, key, err)
		} else {
			fmt.Printf("  ✅ Key %v (%T) is valid\n", key, key)
		}
	}

	// Example 4: Recursion depth validation
	fmt.Println("\n4. Recursion Depth Validation:")
	maxDepth := 10
	testDepths := []int{0, 5, 9, 10, 15}

	for _, depth := range testDepths {
		if err := validateRecursionDepth(depth, maxDepth); err != nil {
			fmt.Printf("  ❌ Depth %d (max %d): %v\n", depth, maxDepth, err)
		} else {
			fmt.Printf("  ✅ Depth %d (max %d) is valid\n", depth, maxDepth)
		}
	}

	// Example 5: Memory limit validation
	fmt.Println("\n5. Memory Limit Validation:")
	memoryLimit := int64(1024 * 1024) // 1MB
	testSizes := []int64{1024, 512 * 1024, 1024 * 1024, 2 * 1024 * 1024}

	for _, size := range testSizes {
		if err := validateMemoryLimit(size, memoryLimit); err != nil {
			fmt.Printf("  ❌ Size %d bytes (limit %d): %v\n", size, memoryLimit, err)
		} else {
			fmt.Printf("  ✅ Size %d bytes (limit %d) is valid\n", size, memoryLimit)
		}
	}

	// Example 6: JSON value type validation
	fmt.Println("\n6. JSON Value Type Validation:")
	validValues := []interface{}{nil, true, "hello", 42, 3.14, []interface{}{1, 2, 3}, map[string]interface{}{"key": "value"}}
	invalidValues := []interface{}{func() {}, make(chan int), struct{}{}}

	for _, value := range validValues {
		if err := validateJSONValueType(value); err != nil {
			fmt.Printf("  ❌ Value %v (%T): %v\n", value, value, err)
		} else {
			fmt.Printf("  ✅ Value %v (%T) is valid JSON type\n", value, value)
		}
	}

	for _, value := range invalidValues {
		if err := validateJSONValueType(value); err != nil {
			fmt.Printf("  ❌ Value %T: %v\n", value, err)
		} else {
			fmt.Printf("  ✅ Value %T is valid JSON type\n", value)
		}
	}

	// Example 7: Composite validation for Get/Set operations
	fmt.Println("\n7. Composite Validation for Operations:")
	
	// Create test JSONValues
	obj := New(map[string]interface{}{"name": "John", "age": 30})
	arr := New([]interface{}{"apple", "banana", "cherry"})

	// Test Get validations
	fmt.Println("  Get Operation Validations:")
	getTests := []struct {
		jv  *JSONValue
		key interface{}
		desc string
	}{
		{obj, "name", "valid object key access"},
		{arr, 1, "valid array index access"},
		{obj, 123, "invalid key type for object"},
		{arr, "invalid", "invalid key type for array"},
		{arr, 10, "out of bounds array access"},
	}

	for _, test := range getTests {
		if err := validateForGet(test.jv, test.key); err != nil {
			fmt.Printf("    ❌ %s: %v\n", test.desc, err)
		} else {
			fmt.Printf("    ✅ %s: valid\n", test.desc)
		}
	}

	// Test Set validations
	fmt.Println("  Set Operation Validations:")
	setTests := []struct {
		jv    *JSONValue
		key   interface{}
		value interface{}
		desc  string
	}{
		{obj, "email", "john@example.com", "valid object set"},
		{arr, 0, "orange", "valid array set"},
		{obj, 123, "value", "invalid key type for object"},
		{arr, "key", "value", "invalid key type for array"},
		{obj, "func", func() {}, "invalid value type"},
	}

	for _, test := range setTests {
		if err := validateForSet(test.jv, test.key, test.value); err != nil {
			fmt.Printf("    ❌ %s: %v\n", test.desc, err)
		} else {
			fmt.Printf("    ✅ %s: valid\n", test.desc)
		}
	}
}

// ExampleValidationInPractice shows how to integrate validation into real operations
func ExampleValidationInPractice() {
	fmt.Println("\n=== Validation in Practice ===")

	// Create a JSONValue
	data := New(map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
			map[string]interface{}{"id": 2, "name": "Bob"},
		},
	})

	// Safe path access with validation
	path := "users.0.name"
	fmt.Printf("\nAccessing path: %s\n", path)
	
	if err := validateForPath(path); err != nil {
		log.Printf("Path validation failed: %v", err)
		return
	}

	result := data.Path(path)
	if !result.IsNull() {
		fmt.Printf("Result: %s\n", result.AsString())
	}

	// Safe array access with validation
	usersArray := data.Get("users")
	index := 0
	fmt.Printf("\nAccessing array index: %d\n", index)
	
	if usersArray.IsArray() {
		if err := validateArrayIndex(index, usersArray.Len()); err != nil {
			log.Printf("Array index validation failed: %v", err)
			return
		}

		user := usersArray.Get(index)
		fmt.Printf("User: %s\n", user.String())
	}

	// Safe key validation before setting
	newKey := "email"
	newValue := "alice@example.com"
	fmt.Printf("\nSetting key: %s = %s\n", newKey, newValue)

	user := data.Path("users.0")
	if err := validateForSet(user, newKey, newValue); err != nil {
		log.Printf("Set operation validation failed: %v", err)
		return
	}

	if err := user.Set(newKey, newValue); err != nil {
		log.Printf("Set operation failed: %v", err)
		return
	}

	fmt.Printf("Updated user: %s\n", user.String())
}

// ExampleCustomValidationLimits shows how to use custom validation limits
func ExampleCustomValidationLimits() {
	fmt.Println("\n=== Custom Validation Limits ===")

	// Get default configuration
	config := DefaultValidationConfig()
	fmt.Printf("Default max path length: %d\n", config.MaxPathLength)
	fmt.Printf("Default max recursion depth: %d\n", config.MaxRecursionDepth)
	fmt.Printf("Default max array size: %d\n", config.MaxArraySize)

	// Example of checking against limits
	testPath := "very.long.path.with.many.segments"
	if len(testPath) > config.MaxPathLength {
		fmt.Printf("Path '%s' exceeds maximum length\n", testPath)
	} else {
		fmt.Printf("Path '%s' is within limits\n", testPath)
	}

	// Test memory estimation
	largeObject := map[string]interface{}{
		"data": make([]interface{}, 1000),
	}
	estimatedSize := estimateJSONSize(largeObject)
	fmt.Printf("Estimated size of large object: %d bytes\n", estimatedSize)

	if estimatedSize > config.MaxJSONSize {
		fmt.Println("Object size exceeds memory limits")
	} else {
		fmt.Println("Object size is within memory limits")
	}
}