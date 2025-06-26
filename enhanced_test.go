package easyjson

import (
	"strings"
	"testing"
)

// enhanced_tests.go - Tests for all new features

func TestGetStringWithDefault(t *testing.T) {
	data := New(map[string]interface{}{
		"name": "John",
		"age":  30,
	})

	// Test existing key
	name := data.GetString("name", "Default")
	if name != "John" {
		t.Errorf("Expected 'John', got '%s'", name)
	}

	// Test missing key with default
	email := data.GetString("email", "no-email@example.com")
	if email != "no-email@example.com" {
		t.Errorf("Expected default email, got '%s'", email)
	}

	// Test nested path
	nested := New(map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"name": "Jane",
			},
		},
	})

	userName := nested.GetString("user", "profile", "name", "Anonymous")
	if userName != "Jane" {
		t.Errorf("Expected 'Jane', got '%s'", userName)
	}

	missingName := nested.GetString("user", "profile", "missing", "Anonymous")
	if missingName != "Anonymous" {
		t.Errorf("Expected 'Anonymous', got '%s'", missingName)
	}
}

func TestTryPaths(t *testing.T) {
	data := New(map[string]interface{}{
		"title":   "Main Title",
		"user":    map[string]interface{}{"name": "John"},
		"product": map[string]interface{}{"label": "Product Label"},
	})

	// Should find first matching path
	result := data.TryPaths("title", "name", "label")
	if result.AsString() != "Main Title" {
		t.Errorf("Expected 'Main Title', got '%s'", result.AsString())
	}

	// Should find second path when first doesn't exist
	result = data.TryPaths("nonexistent", "user.name", "product.label")
	if result.AsString() != "John" {
		t.Errorf("Expected 'John', got '%s'", result.AsString())
	}

	// Should return null when no paths exist
	result = data.TryPaths("missing1", "missing2", "missing3")
	if !result.IsNull() {
		t.Error("Expected null result for all missing paths")
	}
}

func TestFindByField(t *testing.T) {
	data := New(map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
			map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
			map[string]interface{}{"id": 3, "name": "Charlie", "role": "admin"},
		},
	})

	// Find user by ID
	user := data.Get("users").FindByField("id", 2)
	if user.IsNull() {
		t.Error("Expected to find user with ID 2")
	}
	if user.GetString("name") != "Bob" {
		t.Errorf("Expected 'Bob', got '%s'", user.GetString("name"))
	}

	// Find user by role
	admin := data.Get("users").FindByField("role", "admin")
	if admin.IsNull() {
		t.Error("Expected to find admin user")
	}
	if admin.GetString("name") != "Alice" {
		t.Errorf("Expected first admin 'Alice', got '%s'", admin.GetString("name"))
	}

	// Try to find non-existent user
	missing := data.Get("users").FindByField("id", 999)
	if !missing.IsNull() {
		t.Error("Expected null for non-existent user")
	}
}

func TestFilterArray(t *testing.T) {
	data := New(map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "active": true, "age": 25},
			map[string]interface{}{"name": "Bob", "active": false, "age": 30},
			map[string]interface{}{"name": "Charlie", "active": true, "age": 35},
		},
	})

	// Filter active users
	activeUsers := data.Get("users").FilterArray(func(user *JSONValue) bool {
		return user.GetBool("active", false)
	})

	if activeUsers.Len() != 2 {
		t.Errorf("Expected 2 active users, got %d", activeUsers.Len())
	}

	// Check first active user
	firstActive := activeUsers.Get(0)
	if firstActive.GetString("name") != "Alice" {
		t.Errorf("Expected first active user 'Alice', got '%s'", firstActive.GetString("name"))
	}
}

func TestMapArray(t *testing.T) {
	data := New(map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 25},
			map[string]interface{}{"name": "Bob", "age": 30},
		},
	})

	// Extract names
	names := data.Get("users").MapArray(func(user *JSONValue) interface{} {
		return user.GetString("name")
	})

	if names.Len() != 2 {
		t.Errorf("Expected 2 names, got %d", names.Len())
	}

	if names.Get(0).AsString() != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", names.Get(0).AsString())
	}
}

func TestSuggestions(t *testing.T) {
	data := New(map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
		},
		"status": "active",
	})

	smart := WithSuggestions(data)
	suggestions := smart.SuggestPaths()

	// Should include actual keys
	foundUser := false
	for _, suggestion := range suggestions {
		if suggestion == "user" {
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Error("Suggestions should include 'user' key")
	}

	// Test partial completion
	completions := smart.CompletePartial("us")
	foundUserCompletion := false
	for _, completion := range completions {
		if completion == "user" {
			foundUserCompletion = true
			break
		}
	}
	if !foundUserCompletion {
		t.Error("Completions should include 'user' for partial 'us'")
	}
}

func TestSafeParsing(t *testing.T) {
	// Test valid JSON
	result := ParseSafely(`{"name": "John", "age": 30}`)
	if result.Error != nil {
		t.Errorf("Expected no error for valid JSON, got: %v", result.Error)
	}
	if result.Data.GetString("name") != "John" {
		t.Error("Failed to parse valid JSON correctly")
	}

	// Test invalid JSON
	result = ParseSafely(`{"name": "John", "age":}`)
	if result.Error == nil {
		t.Error("Expected error for invalid JSON")
	}
	if result.Data == nil {
		t.Error("Should always return valid JSONValue even on error")
	}
	if len(result.Suggestions) == 0 {
		t.Error("Should provide suggestions for invalid JSON")
	}

	// Test Python-style boolean
	result = ParseSafely(`{"active": True}`)
	if result.Error == nil {
		t.Error("Expected error for Python-style boolean")
	}
	// Check if suggestion mentions boolean case
	foundBoolSuggestion := false
	for _, suggestion := range result.Suggestions {
		if strings.Contains(strings.ToLower(suggestion), "true") {
			foundBoolSuggestion = true
			break
		}
	}
	if !foundBoolSuggestion {
		t.Error("Should suggest using lowercase 'true' instead of 'True'")
	}
}

func TestBuilder(t *testing.T) {
	// Test object building
	obj := NewBuilder().
		AddField("name", "John").
		AddField("age", 30).
		AddObject("profile", func(profile *JSONBuilder) {
			profile.AddField("bio", "Developer").
				AddField("location", "NYC")
		}).
		AddArray("skills", func(skills *JSONBuilder) {
			skills.AddItem("Go").
				AddItem("JavaScript").
				AddItem("Python")
		})

	result := obj.ToJSON()

	if result.GetString("name") != "John" {
		t.Error("Builder failed to set name")
	}

	if result.GetInt("age") != 30 {
		t.Error("Builder failed to set age")
	}

	if result.Get("profile").GetString("bio") != "Developer" {
		t.Error("Builder failed to create nested object")
	}

	if result.Get("skills").Len() != 3 {
		t.Error("Builder failed to create array with correct length")
	}
}

func TestCommonPatterns(t *testing.T) {
	data := New(map[string]interface{}{
		"user_id":    123,
		"full_name":  "John Doe",
		"email":      "john@example.com",
		"user_role":  "admin",
		"created_at": "2023-01-01T00:00:00Z",
	})

	// Test user info extraction
	userInfo := data.GetUserInfo()
	if userInfo["id"] != "123" {
		t.Errorf("Expected user ID '123', got '%s'", userInfo["id"])
	}
	if userInfo["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", userInfo["name"])
	}
	if userInfo["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", userInfo["email"])
	}
	if userInfo["role"] != "admin" {
		t.Errorf("Expected role 'admin', got '%s'", userInfo["role"])
	}

	// Test timestamps extraction
	timestamps := data.GetTimestamps()
	if timestamps["created"] != "2023-01-01T00:00:00Z" {
		t.Errorf("Expected created timestamp, got '%s'", timestamps["created"])
	}
}

func TestEmailValidation(t *testing.T) {
	// Valid emails
	validEmail := New("john@example.com")
	if !validEmail.IsValidEmail() {
		t.Error("Should recognize valid email")
	}

	// Invalid emails
	invalidEmail := New("not-an-email")
	if invalidEmail.IsValidEmail() {
		t.Error("Should not recognize invalid email")
	}

	emptyEmail := New("")
	if emptyEmail.IsValidEmail() {
		t.Error("Should not recognize empty email as valid")
	}
}

func TestDeepSearch(t *testing.T) {
	data := New(map[string]interface{}{
		"level1": map[string]interface{}{
			"level2": map[string]interface{}{
				"target": "found",
				"level3": map[string]interface{}{
					"target": "found_deeper",
				},
			},
		},
		"other": map[string]interface{}{
			"target": "found_other",
		},
	})

	// Should find first occurrence
	result := data.DeepSearch("target")
	if result.IsNull() {
		t.Error("DeepSearch should find target key")
	}

	// The exact value depends on iteration order, but it should find something
	value := result.AsString()
	if value != "found" && value != "found_other" {
		t.Errorf("DeepSearch found unexpected value: %s", value)
	}

	// Should find all occurrences
	allResults := data.DeepSearchAll("target")
	if len(allResults) < 2 {
		t.Errorf("DeepSearchAll should find at least 2 occurrences, found %d", len(allResults))
	}
}

func TestQuickBuilders(t *testing.T) {
	// Test QuickObject
	obj := QuickObject("name", "John", "age", 30, "active", true)
	if obj.GetString("name") != "John" {
		t.Error("QuickObject failed to set name")
	}
	if obj.GetInt("age") != 30 {
		t.Error("QuickObject failed to set age")
	}
	if !obj.GetBool("active") {
		t.Error("QuickObject failed to set active")
	}

	// Test QuickArray
	arr := QuickArray("item1", "item2", "item3")
	if arr.Len() != 3 {
		t.Error("QuickArray should have 3 items")
	}
	if arr.Get(0).AsString() != "item1" {
		t.Error("QuickArray failed to set first item")
	}

	// Test QuickAPIResponse
	response := QuickAPIResponse(
		"success",
		map[string]interface{}{"id": 123},
		"Operation completed",
	)
	if response.GetString("status") != "success" {
		t.Error("QuickAPIResponse failed to set status")
	}
	if response.GetString("message") != "Operation completed" {
		t.Error("QuickAPIResponse failed to set message")
	}
}
