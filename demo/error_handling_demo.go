package main

import (
	"fmt"
	"strings"

	"github.com/javanhut/easyjson"
)

func main() {
	fmt.Println("=== EasyJSON Error Handling Demonstration ===\n")

	// 1. Demonstrate array operations with validation
	fmt.Println("1. Array Operations Error Handling:")
	data := easyjson.QuickArray("item1", "item2", "item3")
	
	// Try to find with nil function - should handle gracefully
	result := data.FindInArray(nil)
	fmt.Printf("   FindInArray with nil function: %v (should be null)\n", result.IsNull())

	// 2. Demonstrate safe parsing with oversized input
	fmt.Println("\n2. Safe Parsing Error Handling:")
	oversizedJSON := "{" + strings.Repeat("\"key\":\"value\",", 100000) + "\"end\":true}"
	parseResult := easyjson.ParseSafely(oversizedJSON)
	fmt.Printf("   Oversized JSON parsing error: %v\n", parseResult.Error != nil)
	fmt.Printf("   Safe fallback data: %v\n", parseResult.Data.IsObject())

	// 3. Demonstrate builder with validation
	fmt.Println("\n3. Builder Error Handling:")
	builder := easyjson.NewBuilder()
	
	// Try to add field with empty key - should be handled
	builder.AddField("", "should_be_ignored")
	builder.AddField("valid_field", "valid_value")
	
	result2 := builder.ToJSON()
	fmt.Printf("   Builder with invalid key handled: %v fields\n", result2.Len())

	// 4. Demonstrate multi-path search with validation
	fmt.Println("\n4. Multi-Path Search Error Handling:")
	testData := easyjson.QuickObject("user", easyjson.QuickObject("name", "John"))
	
	// Try with empty paths array
	result3 := testData.TryPaths()
	fmt.Printf("   TryPaths with no paths: %v (should be null)\n", result3.IsNull())
	
	// Try with valid path
	result4 := testData.TryPaths("user.name", "user.title", "name")
	fmt.Printf("   TryPaths with valid fallback: %v\n", result4.AsString())

	// 5. Demonstrate smart getters with validation
	fmt.Println("\n5. Smart Getters Error Handling:")
	
	// Try with nil JSONValue
	var nilData *easyjson.JSONValue
	safeString := nilData.GetString("any_key", "default")
	fmt.Printf("   GetString on nil value: '%v' (should be default)\n", safeString)

	fmt.Println("\n=== All error handling demonstrations completed successfully! ===")
}