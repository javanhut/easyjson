package main

import (
	"fmt"
	"log"

	"github.com/javanhut/easyjson"
)

func main() {
	fmt.Println("=== EasyJSON Library Demo ===")

	// Example 1: Basic JSON parsing and access
	basicExample()

	// Example 2: Working with arrays
	arrayExample()

	// Example 3: Nested object manipulation
	nestedExample()

	// Example 4: Building JSON from scratch
	buildingExample()

	// Example 5: Safe operations and error handling
	safetyExample()

	// Example 6: Type conversions
	conversionExample()

	// Example 7: Advanced path operations
	pathExample()

	// Example 8: Fluent query syntax
	fluentQueryExample()
}

func basicExample() {
	fmt.Println("1. Basic JSON Parsing and Access")
	fmt.Println("================================")

	jsonStr := `{
		"name": "John Doe",
		"age": 30,
		"email": "john@example.com",
		"active": true
	}`

	data, err := easyjson.Loads(jsonStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", data.Get("name").AsString())
	fmt.Printf("Age: %d\n", data.Get("age").AsInt())
	fmt.Printf("Email: %s\n", data.Get("email").AsString())
	fmt.Printf("Active: %t\n", data.Get("active").AsBool())

	// Modify data
	data.Set("age", 31)
	data.Set("last_login", "2025-05-25")

	result, _ := data.DumpsIndent("  ")
	fmt.Println("\nModified JSON:")
	fmt.Println(result)
	fmt.Println()
}

func arrayExample() {
	fmt.Println("2. Working with Arrays")
	fmt.Println("======================")

	jsonStr := `{
		"fruits": ["apple", "banana", "orange"],
		"numbers": [1, 2, 3, 4, 5]
	}`

	data, _ := easyjson.Loads(jsonStr)

	fruits := data.Get("fruits")
	fmt.Printf("Number of fruits: %d\n", fruits.Len())
	fmt.Printf("First fruit: %s\n", fruits.Get(0).AsString())

	// Add new fruit
	fruits.Append("grape")
	fruits.Extend([]interface{}{"kiwi", "mango"})

	// Iterate through fruits
	fmt.Println("All fruits:")
	for i, fruit := range fruits.AsArray() {
		fmt.Printf("  %d: %s\n", i+1, fruit.AsString())
	}

	// Work with numbers
	numbers := data.Get("numbers")
	sum := 0
	for _, num := range numbers.AsArray() {
		sum += num.AsInt()
	}
	fmt.Printf("Sum of numbers: %d\n", sum)
	fmt.Println()
}

func nestedExample() {
	fmt.Println("3. Nested Object Manipulation")
	fmt.Println("=============================")

	jsonStr := `{
		"user": {
			"personal": {
				"name": "Alice Smith",
				"age": 28
			},
			"contact": {
				"email": "alice@example.com",
				"phone": "555-0123"
			},
			"preferences": {
				"theme": "dark",
				"notifications": true
			}
		}
	}`

	data, _ := easyjson.Loads(jsonStr)

	// Access nested values using Path
	fmt.Printf("Name: %s\n", data.Path("user.personal.name").AsString())
	fmt.Printf("Email: %s\n", data.Path("user.contact.email").AsString())
	fmt.Printf("Theme: %s\n", data.Path("user.preferences.theme").AsString())

	// Modify nested values
	data.SetPath("user.personal.age", 29)
	data.SetPath("user.contact.address", "123 Main St")
	data.SetPath("user.preferences.language", "en")

	// Show all user keys
	userObj := data.Get("user")
	fmt.Println("\nUser object keys:")
	for _, key := range userObj.Keys() {
		fmt.Printf("  - %s\n", key)
	}

	result, _ := data.DumpsIndent("  ")
	fmt.Println("\nUpdated nested structure:")
	fmt.Println(result)
	fmt.Println()
}

func buildingExample() {
	fmt.Println("4. Building JSON from Scratch")
	fmt.Println("=============================")

	// Create a new response object
	response := easyjson.NewObject()
	response.Set("status", "success")
	response.Set("timestamp", 1716638400) // Unix timestamp

	// Create user data
	user := easyjson.NewObject()
	user.Set("id", 12345)
	user.Set("username", "johndoe")
	user.Set("verified", true)

	// Create permissions array
	permissions := easyjson.NewArrayFrom([]interface{}{
		"read", "write", "delete",
	})

	// Create settings object
	settings := easyjson.NewObject()
	settings.Set("theme", "light")
	settings.Set("language", "en")
	settings.Set("notifications", true)

	// Assemble the structure
	user.Set("permissions", permissions.Raw())
	user.Set("settings", settings.Raw())
	response.Set("user", user.Raw())

	// Add metadata
	metadata := easyjson.NewObject()
	metadata.Set("version", "1.0")
	metadata.Set("server", "api-01")
	response.Set("metadata", metadata.Raw())

	result, _ := response.DumpsIndent("  ")
	fmt.Println("Built JSON structure:")
	fmt.Println(result)
	fmt.Println()
}

func safetyExample() {
	fmt.Println("5. Safe Operations and Error Handling")
	fmt.Println("=====================================")

	// Example with missing fields
	jsonStr := `{"user": {"name": "Bob"}}`
	data, _ := easyjson.Loads(jsonStr)

	// Safe access to missing fields
	age := data.Path("user.age").AsInt()                // Returns 0 for missing field
	email := data.Path("user.contact.email").AsString() // Returns "" for missing nested field
	active := data.Path("user.active").AsBool()         // Returns false for missing field

	fmt.Printf("Age (missing): %d\n", age)
	fmt.Printf("Email (missing): '%s'\n", email)
	fmt.Printf("Active (missing): %t\n", active)

	// Check existence before access
	if data.Path("user.name").IsNull() {
		fmt.Println("Name is missing")
	} else {
		fmt.Printf("Name exists: %s\n", data.Path("user.name").AsString())
	}

	// Safe type checking
	userObj := data.Get("user")
	fmt.Printf("User is object: %t\n", userObj.IsObject())
	fmt.Printf("User is array: %t\n", userObj.IsArray())
	fmt.Printf("User is null: %t\n", userObj.IsNull())

	// Handle invalid JSON gracefully
	invalidJSON := `{"broken": json}`
	_, err := easyjson.Loads(invalidJSON)
	if err != nil {
		fmt.Printf("Handled JSON parsing error: %v\n", err)
	}
	fmt.Println()
}

func conversionExample() {
	fmt.Println("6. Type Conversions")
	fmt.Println("==================")

	jsonStr := `{
		"string_number": "42",
		"float_number": 3.14159,
		"bool_string": "true",
		"number_zero": 0,
		"mixed_array": [1, "2", 3.0, true]
	}`

	data, _ := easyjson.Loads(jsonStr)

	// String to number conversion
	stringNum := data.Get("string_number")
	fmt.Printf("String '42' as int: %d\n", stringNum.AsInt())
	fmt.Printf("String '42' as float: %.2f\n", stringNum.AsFloat())

	// Number precision
	floatNum := data.Get("float_number")
	fmt.Printf("Float %.5f as int: %d\n", floatNum.AsFloat(), floatNum.AsInt())

	// String to boolean
	boolStr := data.Get("bool_string")
	fmt.Printf("String 'true' as bool: %t\n", boolStr.AsBool())

	// Zero/falsy values
	zero := data.Get("number_zero")
	fmt.Printf("Number 0 as bool: %t\n", zero.AsBool())

	// Mixed array conversions
	mixedArray := data.Get("mixed_array")
	fmt.Println("Mixed array conversions:")
	for i, item := range mixedArray.AsArray() {
		fmt.Printf("  [%d] as string: '%s', as int: %d, as bool: %t\n",
			i, item.AsString(), item.AsInt(), item.AsBool())
	}
	fmt.Println()
}

func pathExample() {
	fmt.Println("7. Advanced Path Operations")
	fmt.Println("===========================")

	jsonStr := `{
		"company": {
			"departments": [
				{
					"name": "Engineering",
					"employees": [
						{"name": "Alice", "role": "Senior Engineer", "skills": ["Go", "Python", "Docker"]},
						{"name": "Bob", "role": "DevOps", "skills": ["AWS", "Kubernetes", "Terraform"]}
					]
				},
				{
					"name": "Marketing",
					"employees": [
						{"name": "Carol", "role": "Marketing Manager", "skills": ["SEO", "Analytics", "Content"]}
					]
				}
			]
		}
	}`

	data, _ := easyjson.Loads(jsonStr)

	// Access deeply nested data
	firstDept := data.Path("company.departments.0.name").AsString()
	fmt.Printf("First department: %s\n", firstDept)

	firstEmployee := data.Path("company.departments.0.employees.0.name").AsString()
	fmt.Printf("First employee: %s\n", firstEmployee)

	firstSkill := data.Path("company.departments.0.employees.0.skills.0").AsString()
	fmt.Printf("First skill of first employee: %s\n", firstSkill)

	// Iterate through complex structures
	departments := data.Path("company.departments").AsArray()
	for deptIdx, dept := range departments {
		deptName := dept.Get("name").AsString()
		fmt.Printf("\nDepartment %d: %s\n", deptIdx+1, deptName)

		employees := dept.Get("employees").AsArray()
		for empIdx, emp := range employees {
			name := emp.Get("name").AsString()
			role := emp.Get("role").AsString()
			fmt.Printf("  Employee %d: %s (%s)\n", empIdx+1, name, role)

			skills := emp.Get("skills").AsArray()
			fmt.Printf("    Skills: ")
			for skillIdx, skill := range skills {
				if skillIdx > 0 {
					fmt.Print(", ")
				}
				fmt.Print(skill.AsString())
			}
			fmt.Println()
		}
	}

	// Add new data using paths
	newEmployee := easyjson.NewObject()
	newEmployee.Set("name", "Dave")
	newEmployee.Set("role", "Junior Engineer")
	newEmployee.Set("skills", []interface{}{"JavaScript", "React", "Node.js"})

	// This would require extending the employees array - showing correct approach
	fmt.Println("\nAdding new employee to Engineering department...")
	engDept := data.Path("company.departments.0")
	engEmployees := engDept.Get("employees")
	engEmployees.Append(newEmployee.Raw())

	// Alternative: If you want to use SetPath with arrays, create structure first
	// For example, if adding a new department:
	data.SetPath(
		"company.departments",
		append(data.Path("company.departments").Raw().([]interface{}), map[string]interface{}{
			"name":      "Sales",
			"employees": []interface{}{},
		}),
	)

	result, _ := data.DumpsIndent("  ")
	fmt.Println("Updated company structure:")
	fmt.Println(result)
}

func fluentQueryExample() {
	fmt.Println("\n8. Fluent Query Syntax (Python-like)")
	fmt.Println("====================================")

	jsonStr := `{
		"users": [
			{
				"name": "Alice",
				"profile": {
					"hair_color": "Red",
					"age": 25,
					"preferences": {
						"theme": "dark",
						"notifications": true
					}
				}
			},
			{
				"name": "Bob",
				"profile": {
					"hair_color": "Brown", 
					"age": 30,
					"preferences": {
						"theme": "light",
						"notifications": false
					}
				}
			}
		],
		"metadata": {
			"total_users": 2,
			"last_updated": "2025-05-25"
		}
	}`

	data, _ := easyjson.Loads(jsonStr)

	// Python-like access: data["users"][0]["profile"]["hair_color"]
	// Go equivalent: data.Q("users", 0, "profile", "hair_color").AsString()
	fmt.Println("=== Fluent Query Access ===")

	// Single chain access
	hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
	fmt.Printf("Alice's hair color: %s\n", hairColor)

	// Deep nested access
	theme := data.Q("users", 1, "profile", "preferences", "theme").AsString()
	fmt.Printf("Bob's theme preference: %s\n", theme)

	// Mixed types in chain
	totalUsers := data.Q("metadata", "total_users").AsInt()
	fmt.Printf("Total users: %d\n", totalUsers)

	// Safe access to missing keys
	missing := data.Q("users", 0, "profile", "nonexistent_field").AsString()
	fmt.Printf("Missing field (safe): '%s'\n", missing) // Returns empty string

	// Check if value exists before using
	notifications := data.Q("users", 0, "profile", "preferences", "notifications")
	if !notifications.IsNull() {
		fmt.Printf("Alice's notifications: %t\n", notifications.AsBool())
	}

	fmt.Println("\n=== Comparison of Access Methods ===")

	// Traditional method
	traditional := data.Get("users").Get(0).Get("profile").Get("hair_color").AsString()
	fmt.Printf("Traditional: %s\n", traditional)

	// Path method
	pathMethod := data.Path("users.0.profile.hair_color").AsString()
	fmt.Printf("Path method: %s\n", pathMethod)

	// Fluent query method
	fluentMethod := data.Q("users", 0, "profile", "hair_color").AsString()
	fmt.Printf("Fluent query: %s\n", fluentMethod)

	fmt.Println("\n=== Iterating with Fluent Access ===")

	// Get users array and iterate
	usersArray := data.Q("users").AsArray()
	for i, user := range usersArray {
		name := user.Q("name").AsString()
		age := user.Q("profile", "age").AsInt()
		theme := user.Q("profile", "preferences", "theme").AsString()

		fmt.Printf("User %d: %s (age %d, prefers %s theme)\n",
			i+1, name, age, theme)
	}

	fmt.Println("\n=== Error-Safe Deep Access ===")

	// These won't panic even with invalid paths
	invalid1 := data.Q("nonexistent").AsString()
	invalid2 := data.Q("users", 999, "name").AsString()
	invalid3 := data.Q("users", 0, "invalid", "chain", "here").AsString()

	fmt.Printf("Nonexistent key: '%s'\n", invalid1)
	fmt.Printf("Invalid array index: '%s'\n", invalid2)
	fmt.Printf("Invalid nested chain: '%s'\n", invalid3)
}
