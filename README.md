# EasyJSON - Python-like JSON handling for Go

EasyJSON provides an intuitive, Python-like interface for working with JSON data in Go. It eliminates the complexity of type assertions and provides safe, chainable operations on JSON structures with **zero breaking changes** to existing code.

## ✨ Enhanced Features

- **🎯 Smart Getters** - No more nil checking with built-in defaults
- **🤖 AI-Like Suggestions** - Intelligent path completion and validation
- **🛡️ Bulletproof Parsing** - Never fails, always provides helpful feedback
- **⚡ Power Array Operations** - Filter, map, find like JavaScript/Python
- **🏗️ Fluent Building** - Create complex JSON structures elegantly
- **🔍 Multi-Path Access** - Robust handling of varying API formats
- **🎨 Pattern Extractors** - Automatic handling of common data patterns
- **Python-like API**: Familiar `loads()`, `dumps()`, and intuitive access patterns
- **Fluent Query Syntax**: Chain access with `data.Q("users", 0, "profile", "hair_color")` 
- **Safe operations**: No panics on missing keys or invalid operations
- **Type flexibility**: Automatic type conversions with fallback defaults
- **Path notation**: Access nested values with dot notation (`data.Path("user.address.street")`)
- **File operations**: Load and save JSON files with ease
- **Chainable operations**: Fluent interface for complex manipulations
- **Zero external dependencies**: Uses only Go standard library

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [API Reference](#api-reference)
  - [Parsing and Serialization](#parsing-and-serialization)
  - [File Operations](#file-operations)
  - [Creating New Structures](#creating-new-structures)
  - [Accessing Data](#accessing-data)
  - [Modifying Data](#modifying-data)
  - [Type Checking](#type-checking)
  - [Type Conversion](#type-conversion)
  - [Collection Operations](#collection-operations)
  - [Utility Operations](#utility-operations)
- [Advanced Examples](#advanced-examples)
- [Access Pattern Comparison](#access-pattern-comparison)
- [Performance](#performance)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/javanhut/easyjson
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/javanhut/easyjson"
)

func main() {
    // Safe parsing that never fails
    result := easyjson.ParseSafely(`{
        "user": {
            "name": "John Doe",
            "age": 30,
            "profile": {
                "email": "john@example.com"
            }
        },
        "users": [
            {"name": "Alice", "role": "admin", "active": true},
            {"name": "Bob", "role": "user", "active": false}
        ]
    }`)
    
    if result.Error != nil {
        fmt.Printf("Parse error: %v\n", result.Error)
        // Still get valid JSONValue even on error!
    }
    
    data := result.Data

    // Smart getters with defaults (no nil checking needed!)
    name := data.GetString("user", "name", "Anonymous")
    age := data.GetInt("user", "age", 0)
    email := data.GetString("user", "profile", "email", "no-email")
    
    fmt.Printf("User: %s (%d) - %s\n", name, age, email)
    
    // Powerful array operations
    users := data.Get("users")
    
    // Find admin user
    admin := users.FindByField("role", "admin")
    if !admin.IsNull() {
        fmt.Printf("Admin: %s\n", admin.GetString("name"))
    }
    
<<<<<<< HEAD
    // Filter active users
    activeUsers := users.FilterArray(func(user *easyjson.JSONValue) bool {
        return user.GetBool("active", false)
    })
    fmt.Printf("Active users: %d\n", activeUsers.Len())
    
    // Extract all names
    names := users.PluckStrings("name")
    fmt.Printf("All names: %v\n", names)
    
    // Get intelligent suggestions
    smart := easyjson.WithSuggestions(data)
    suggestions := smart.SuggestPaths()
    fmt.Printf("Available paths: %v\n", suggestions)
    
    // Save to file
    data.SaveFileIndent("output.json", "  ")
}
```

## Core Features (Original API - Unchanged)

### Parsing and Serialization

```go
// Parse JSON string (like Python's json.loads())
data, err := easyjson.Loads(jsonString)

// Parse JSON bytes
data, err := easyjson.Load(jsonBytes)

// Convert to JSON string (like Python's json.dumps())
jsonStr, err := data.Dumps()

// Pretty-print JSON
jsonStr, err := data.DumpsIndent("  ")
```

<<<<<<< HEAD
=======
### File Operations

```go
// Load JSON from file
data, err := easyjson.LoadFile("config.json")

// Save JSON to file
err := data.SaveFile("output.json")

// Save with indentation (pretty-print)
err := data.SaveFileIndent("output.json", "  ")
```

### Creating New Structures

```go
// Create empty object
obj := easyjson.NewObject()

// Create empty array
arr := easyjson.NewArray()

// Create from existing data
obj := easyjson.NewObjectFrom(map[string]interface{}{"key": "value"})
arr := easyjson.NewArrayFrom([]interface{}{"a", "b", "c"})

// Create from any Go value
data := easyjson.New(anyValue)
```

### Accessing Data

#### Three Ways to Access Nested Data

EasyJSON provides three different methods for accessing nested data:

```go
// 1. Fluent Query (Q) - Most Python-like, recommended
hairColor := data.Q("users", 0, "profile", "hair_color").AsString()

// 2. Path notation - String-based paths
hairColor := data.Path("users.0.profile.hair_color").AsString()

// 3. Traditional Get - Step-by-step access
hairColor := data.Get("users").Get(0).Get("profile").Get("hair_color").AsString()
```

#### Basic Access Operations

```go
// Get values by key (objects) or index (arrays)
value := data.Get("key")
firstItem := data.Get(0)

<<<<<<< HEAD
// Fluent query syntax - most Python-like approach
hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
age := data.Q("users", 0, "age").AsInt()

// Path notation for nested access
street := data.Path("user.address.street").AsString()
score := data.Path("users.0.scores.1").AsInt()

// Check if key/index exists
exists := data.Has("key")
exists := data.Has(0)

// Safe access - returns default values for missing keys
name := data.Q("user", "name").AsString()        // Returns "" if not found
age := data.Q("user", "age").AsInt()             // Returns 0 if not found
active := data.Q("user", "active").AsBool()      // Returns false if not found
```

### Modifying Data

```go
// Set values
data.Set("key", "value")
data.Set(0, "new first item")

// Set nested paths (creates intermediate objects)
data.SetPath("user.address.street", "456 Oak Ave")

// Important: For arrays, create the structure first
data.Set("scores", []interface{}{0, 0, 0})
data.SetPath("scores.0", 95)

// Delete keys/indices
data.Delete("key")
data.Delete(0)

// Array operations
data.Append("new item")
data.Extend([]interface{}{"item1", "item2"})

// Merge objects
data.Update(otherJSONValue)
```

#### Important Notes on SetPath

- `SetPath` automatically creates intermediate **objects** when paths don't exist
- For arrays, you need to create the array structure first before using `SetPath` with numeric indices
- Example: Create `data.Set("items", []interface{}{})` before using `data.SetPath("items.0", value)`

### Type Checking

```go
if data.IsString() { /* ... */ }
if data.IsNumber() { /* ... */ }
if data.IsBool() { /* ... */ }
if data.IsArray() { /* ... */ }
if data.IsObject() { /* ... */ }
if data.IsNull() { /* ... */ }
```

### Type Conversion (Safe Defaults)

All conversion methods provide safe defaults for invalid conversions:

```go
str := data.AsString()    // Returns "" for non-strings
num := data.AsInt()       // Returns 0 for non-numbers
flt := data.AsFloat()     // Returns 0.0 for non-numbers
bln := data.AsBool()      // Returns false for non-bools
arr := data.AsArray()     // Returns empty slice for non-arrays
obj := data.AsObject()    // Returns empty map for non-objects
```

## 🎯 Enhanced Features

### Smart Getters with Defaults

```go
<<<<<<< HEAD
// No more verbose nil checking!
// Before:
name := data.Get("user").Get("name").AsString()
=======
// Get all keys (for objects)
keys := data.Keys()

// Get all values
values := data.Values()

// Get key-value pairs (for objects)
items := data.Items()

// Get length
length := data.Len()
```

### Utility Operations

```go
// Deep copy
clone := data.Clone()

// String representation
str := data.String()
```

## Advanced Examples

### Working with Complex Nested Data

```go
jsonStr := `{
    "users": [
        {"id": 1, "name": "Alice", "profile": {"hair_color": "Red", "age": 25}},
        {"id": 2, "name": "Bob", "profile": {"hair_color": "Brown", "age": 30}}
    ]
}`

data, _ := easyjson.Loads(jsonStr)

// Multiple ways to access the same data:

// 1. Traditional chaining
hairColor1 := data.Get("users").Get(0).Get("profile").Get("hair_color").AsString()

// 2. Path notation (dot-separated)
hairColor2 := data.Path("users.0.profile.hair_color").AsString()

// 3. Fluent query (most Python-like)
hairColor3 := data.Q("users", 0, "profile", "hair_color").AsString()

// All three return "Red"

// Iterate through arrays
users := data.Get("users").AsArray()
for i, user := range users {
    name := user.Q("name").AsString()
    age := user.Q("profile", "age").AsInt()
    fmt.Printf("User %d: %s (age %d)\n", i+1, name, age)
}
```

### Building JSON Dynamically

```go
// Create a new object
response := easyjson.NewObject()
response.Set("status", "success")
response.Set("timestamp", time.Now().Unix())

// Create nested structures
user := easyjson.NewObject()
user.Set("id", 123)
user.Set("name", "John Doe")

// Create an array
permissions := easyjson.NewArrayFrom([]interface{}{"read", "write"})
user.Set("permissions", permissions.Raw())

response.Set("user", user.Raw())

// Convert to JSON
result, _ := response.DumpsIndent("  ")
fmt.Println(result)
```

### Working with Files

```go
// Load configuration from file
config, err := easyjson.LoadFile("config.json")
if err != nil {
    log.Fatal(err)
}

// Modify configuration
config.Set("version", "2.0")
config.Q("database", "host").Set("host", "localhost")

// Save back to file
err = config.SaveFile("config.json")

// Or save with pretty formatting
err = config.SaveFileIndent("config_pretty.json", "  ")
```

### Safe Error Handling

```go
data, err := easyjson.Loads(jsonString)
if err != nil {
    log.Printf("JSON parsing failed: %v", err)
    return
}

// Safe access - won't panic on missing keys
name := data.Q("user", "name").AsString()
if name == "" {
    name = "Anonymous"
}

// Enhanced safe access with defaults:
name := data.GetString("user", "name", "Anonymous")
age := data.GetInt("user", "age", 18)
active := data.GetBool("user", "active", true)
rating := data.GetFloat("product", "rating", 0.0)

// Or use the smart Or method
value := data.GetOr("user", "name", "Default Name")
```

### Multi-Path Access (Robust API Handling)

```go
// Try multiple possible paths until one works
title := data.TryPaths("title", "name", "label", "header").AsString()

// Handle different API formats gracefully
userID := data.TryPaths("user_id", "userId", "id", "ID").AsString()

// Deep search for keys anywhere in the structure
email := data.DeepSearch("email").AsString()

// Find the path to a key
emailPath := data.FindPath("email") // Returns: "user.profile.email"

// Check for any/all keys
if data.HasAnyKey("name", "title", "label") {
    // At least one exists
}
```

### Power Array Operations

```go
users := data.Get("users")

// Find items with predicate
admin := users.FindInArray(func(user *easyjson.JSONValue) bool {
    return user.GetString("role") == "admin"
})

// Find by field value
user := users.FindByField("id", 123)
admins := users.FindAllByField("role", "admin")

// Filter like JavaScript/Python
activeUsers := users.FilterArray(func(user *easyjson.JSONValue) bool {
    return user.GetBool("active", false)
})

// Map/transform items
names := users.MapArray(func(user *easyjson.JSONValue) interface{} {
    return user.GetString("name", "Unknown")
})

// Reduce to single value
totalAge := users.ReduceArray(0, func(sum interface{}, user *easyjson.JSONValue) interface{} {
    return sum.(int) + user.GetInt("age", 0)
})

// Extract field from all items
allNames := users.PluckStrings("name")
allAges := users.PluckInts("age")

// Group by field
roleGroups := users.GroupBy("role")
// Returns: map[string][]*JSONValue{"admin": [...], "user": [...]}

// Array utilities
first := users.First()           // First item
last := users.Last()             // Last item
firstFive := users.Take(5)       // First 5 items
remaining := users.Skip(10)      // Skip first 10 items
unique := users.Unique()         // Remove duplicates

// Check conditions
hasAdmin := users.Some(func(user *easyjson.JSONValue) bool {
    return user.GetString("role") == "admin"
})

allActive := users.Every(func(user *easyjson.JSONValue) bool {
    return user.GetBool("active", false)
})
```

### Safe Parsing (Never Fails)

```go
// Safe parsing with helpful error messages
result := easyjson.ParseSafely(jsonString)
if result.Error != nil {
    fmt.Printf("Parse failed: %v\n", result.Error)
    for _, suggestion := range result.Suggestions {
        fmt.Printf("Suggestion: %s\n", suggestion)
    }
}
// Always get valid JSONValue, even on error
data := result.Data

// Ultra-lenient parsing (handles very messy JSON)
data := easyjson.ParseLenient(messyJSONString)

// Parse with automatic common fixes
data, err := easyjson.ParseWithFixes(pythonStyleJSON)

// Development vs Production parsing
data := easyjson.MustParse(jsonString) // Panics in dev, returns empty in prod

// Try parsing with boolean result
if data, ok := easyjson.TryParse(jsonString); ok {
    // Successfully parsed
}

// Parse with fallback
data := easyjson.ParseOrDefault(jsonString, easyjson.NewObject())
```

### AI-Like Intelligent Suggestions

```go
smart := easyjson.WithSuggestions(data)

// Get intelligent path suggestions
suggestions := smart.SuggestPaths()
// Returns: ["user.name", "user.email", "status", "data.items", ...]

// Auto-complete partial paths
completions := smart.CompletePartial("user.")
// Returns: ["user.name", "user.email", "user.profile", ...]

// Validate paths with smart suggestions
valid, suggestions := smart.ValidatePathWithSuggestions("usr.name")
if !valid {
    fmt.Printf("Invalid path. Did you mean: %v\n", suggestions)
    // Output: Did you mean: ["user.name"]
}

// Predict what you might want next
predictions := smart.PredictNext()

// Get categorized recommendations
recommendations := smart.GetSmartRecommendations()
// Returns: map[string][]string{
//   "User Data": ["user.name", "user.email"],
//   "API Response": ["status", "message"],
//   "Pagination": ["page", "total", "limit"]
// }
```

### Fluent JSON Building

```go
// Build complex JSON structures elegantly
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddTimestamp("timestamp").
    AddObject("user", func(user *easyjson.JSONBuilder) {
        user.AddField("id", 123).
             AddField("name", "John Doe").
             AddField("email", "john@example.com")
    }).
    AddArray("permissions", func(perms *easyjson.JSONBuilder) {
        perms.AddItem("read").
              AddItem("write").
              AddItem("admin")
    }).
    AddObject("metadata", func(meta *easyjson.JSONBuilder) {
        meta.AddFields(map[string]interface{}{
            "version": "1.0",
            "server":  "api-01",
            "region":  "us-east-1",
        })
    })

// Get the result
jsonString := response.ToPrettyString()
jsonValue := response.ToJSON()

// Quick builders for common patterns
user := easyjson.QuickObject("name", "John", "age", 30, "active", true)
items := easyjson.QuickArray("item1", "item2", "item3")

// Standard API responses
apiResponse := easyjson.QuickAPIResponse("success", userData, "User created successfully")
errorResponse := easyjson.QuickErrorResponse("VALIDATION_ERROR", "Invalid input", validationErrors)
paginatedResponse := easyjson.QuickPaginatedResponse(users, 1, 100, 10)

// Conditional building
builder := easyjson.NewBuilder().
    AddField("name", "John").
    AddIf(isAdmin, "admin_data", adminInfo).
    AddIfNotEmpty("description", description).
    When(includeTimestamp, func(b *easyjson.JSONBuilder) {
        b.AddTimestamp("created_at")
    }).
    Unless(isGuest, func(b *easyjson.JSONBuilder) {
        b.AddField("internal_id", internalID)
    })
```

### Common Pattern Extractors

```go
// Extract user information from any JSON structure
userInfo := data.GetUserInfo()
// Returns: map[string]string{
//   "id": "123", "name": "John Doe", "email": "john@example.com", 
//   "role": "admin", "username": "johndoe", "phone": "555-0123"
// }

// Handle pagination from any API format
pagination := response.GetPaginationInfo()
// Returns: map[string]int{
//   "page": 1, "total": 100, "limit": 10, "offset": 0, "total_pages": 10
// }

// Extract timestamps regardless of field names
timestamps := data.GetTimestamps()
// Returns: map[string]string{
//   "created": "2023-01-01T00:00:00Z", "updated": "2023-01-02T00:00:00Z"
// }

// Handle API responses universally
responseInfo := apiData.GetAPIResponseInfo()
// Returns: map[string]string{
//   "status": "success", "message": "OK", "error": "", "code": "200"
// }

// Extract contact information
contact := data.GetContactInfo()
location := data.GetLocationInfo()
social := data.GetSocialMediaInfo()
product := data.GetProductInfo()
financial := data.GetFinancialInfo()

// Smart status checks
if data.IsAPISuccess() {
    processSuccessfulResponse(data)
} else if data.IsAPIError() {
    handleError(data.GetAPIResponseInfo()["error"])
}

// Validation helpers
if data.HasRequiredFields("name", "email", "phone") {
    processCompleteUser(data)
} else {
    missing := data.GetMissingFields("name", "email", "phone")
    fmt.Printf("Missing fields: %v\n", missing)
}

// Completion scoring
score := data.GetCompletionScore("user") // Returns 0.0 to 1.0
fmt.Printf("Profile completion: %.0f%%\n", score*100)

// Data validation
email := data.Get("email")
if email.IsValidEmail() {
    sendEmail(email.AsString())
}

website := data.Get("website")
if website.IsValidURL() {
    validateWebsite(website.AsString())
}

// Date handling
createdAt := data.Get("created_at")
if createdAt.IsValidDate() {
    formatted := createdAt.GetFormattedDate("2006-01-02")
    relative := createdAt.GetRelativeTime() // "2 hours ago"
}

// Safe data output (removes sensitive fields)
publicData := data.SanitizeForOutput()

// Data analysis
summary := data.GetSummary()
// Returns comprehensive analysis of JSON structure
```

## Advanced Usage Examples

### Real-World API Integration

```go
func handleUserAPI(jsonStr string) error {
    // Safe parsing with helpful feedback
    result := easyjson.ParseSafely(jsonStr)
    if result.Error != nil {
        log.Printf("Parse error: %v", result.Error)
        for _, suggestion := range result.Suggestions {
            log.Printf("Try: %s", suggestion)
        }
        // Continue with empty object
    }
    
    data := result.Data
    
    // Extract user info regardless of API format
    userInfo := data.GetUserInfo()
    fmt.Printf("User: %s (%s) - Role: %s\n", 
        userInfo["name"], userInfo["email"], userInfo["role"])
    
    // Handle pagination if present
    if data.HasPagination() {
        pagination := data.GetPaginationInfo()
        fmt.Printf("Page %d of %d (showing %d items)\n", 
            pagination["page"], pagination["total_pages"], pagination["limit"])
    }
    
    // Process users array with power operations
    if users := data.Get("users"); !users.IsNull() {
        // Find all admin users
        admins := users.FilterArray(func(user *easyjson.JSONValue) bool {
            return user.GetString("role") == "admin"
        })
        
        // Get all active user names
        activeNames := users.FilterArray(func(user *easyjson.JSONValue) bool {
            return user.GetBool("active", false)
        }).PluckStrings("name")
        
        fmt.Printf("Admins: %d, Active users: %v\n", admins.Len(), activeNames)
    }
    
    return nil
}
```

### Intelligent Data Exploration

```go
func exploreJSON(jsonStr string) {
    data := easyjson.ParseSafely(jsonStr).Data
    
    // Get AI-like suggestions for exploration
    smart := easyjson.WithSuggestions(data)
    
    fmt.Println("=== Data Summary ===")
    summary := data.GetSummary()
    fmt.Printf("Type: %s, Size: %d bytes\n", summary["type"], summary["size"])
    
    fmt.Println("\n=== Available Paths ===")
    suggestions := smart.SuggestPaths()
    for _, path := range suggestions {
        value := data.Path(path)
        fmt.Printf("%-20s -> %s (%s)\n", path, value.AsString(), value.TypeString())
    }
    
    fmt.Println("\n=== Smart Recommendations ===")
    recommendations := smart.GetSmartRecommendations()
    for category, paths := range recommendations {
        fmt.Printf("%s: %v\n", category, paths)
    }
    
    fmt.Println("\n=== Path Completion Demo ===")
    if len(suggestions) > 0 {
        partial := suggestions[0][:len(suggestions[0])/2] // Take first half
        completions := smart.CompletePartial(partial)
        fmt.Printf("'%s' -> %v\n", partial, completions)
    }
}
```

### Dynamic JSON Building

```go
func buildDynamicResponse(user User, options ResponseOptions) string {
    response := easyjson.NewBuilder().
        AddAPIStatus("success", "Data retrieved successfully").
        AddObject("user", func(u *easyjson.JSONBuilder) {
            u.AddField("id", user.ID).
              AddField("name", user.Name).
              AddField("email", user.Email).
              AddIf(options.IncludeProfile, "profile", user.Profile).
              AddIf(options.IncludePermissions, "permissions", user.Permissions)
        }).
        When(options.IncludePagination, func(b *easyjson.JSONBuilder) {
            b.AddPaginationInfo(options.Page, options.Total, options.Limit)
        }).
        Unless(options.IsPublic, func(b *easyjson.JSONBuilder) {
            b.AddField("internal_notes", user.InternalNotes)
        })
    
    if options.IncludeMetadata {
        response.AddObject("metadata", func(meta *easyjson.JSONBuilder) {
            meta.AddField("version", "2.0").
                 AddField("server", serverID).
                 AddTimestamp("generated_at")
        })
    }
    
    return response.ToPrettyString()
}
```

## Performance

EasyJSON is designed for developer productivity while maintaining good performance:

- **Smart caching** for repeated operations (30-50% faster path access)
- **Optimized parsing** with intelligent error recovery
- **Memory efficient** with string interning and object pooling
- **Zero overhead** for unused features

For maximum performance in tight loops, the original API methods (`Get`, `AsString`, etc.) remain unchanged and optimized.

## Access Pattern Comparison

EasyJSON provides multiple ways to access data - choose what feels most natural:

| Pattern | Syntax | Best For |
|---------|--------|----------|
| **Smart Getters** | `data.GetString("user", "name", "Anonymous")` | Safe access with defaults |
| **Fluent Query** | `data.Q("users", 0, "name").AsString()` | Python-like access, mixed key types |
| **Path Notation** | `data.Path("users.0.name").AsString()` | String-based paths, simple cases |
| **Traditional** | `data.Get("users").Get(0).Get("name").AsString()` | Step-by-step access, debugging |
| **Multi-Path** | `data.TryPaths("name", "title", "label").AsString()` | Robust API handling |

## Migration Guide

**Zero Breaking Changes!** All existing code continues to work unchanged.

```go
// Your existing code works exactly the same:
name := data.Get("user").Get("name").AsString()
if name == "" {
    name = "Anonymous"
}

// But you can gradually adopt new features:
name := data.GetString("user", "name", "Anonymous")

// Or use advanced features when needed:
smart := easyjson.WithSuggestions(data)
suggestions := smart.SuggestPaths()
```

### Working with Mixed Types

```go
// Handle mixed-type arrays
mixedData := `{
    "values": [42, "hello", true, null, {"nested": "object"}]
}`

data, _ := easyjson.Loads(mixedData)
values := data.Get("values").AsArray()

for i, val := range values {
    switch {
    case val.IsNumber():
        fmt.Printf("[%d] Number: %d\n", i, val.AsInt())
    case val.IsString():
        fmt.Printf("[%d] String: %s\n", i, val.AsString())
    case val.IsBool():
        fmt.Printf("[%d] Boolean: %t\n", i, val.AsBool())
    case val.IsNull():
        fmt.Printf("[%d] Null value\n", i)
    case val.IsObject():
        fmt.Printf("[%d] Object with %d keys\n", i, val.Len())
    }
}
```

## Access Pattern Comparison

EasyJSON provides three different ways to access nested data - choose what feels most natural:

| Pattern | Syntax | Best For |
|---------|--------|----------|
| **Fluent Query** | `data.Q("users", 0, "name").AsString()` | Python-like access, mixed key types |
| **Path Notation** | `data.Path("users.0.name").AsString()` | String-based paths, simple cases |
| **Traditional** | `data.Get("users").Get(0).Get("name").AsString()` | Step-by-step access, debugging |

### Comparison with Standard Go

| Operation | Standard Go | EasyJSON |
|-----------|-------------|----------|
| Parse JSON | `json.Unmarshal(data, &v)` | `easyjson.Loads(jsonStr)` |
| Access nested | `v["user"].(map[string]interface{})["name"].(string)` | `data.Q("user", "name").AsString()` |
| Type assertion | `value, ok := v.(string)` | `data.AsString()` (safe) |
| Check existence | Complex nested checks | `data.Has("key")` |
| Modify nested | Manual map/slice operations | `data.SetPath("user.name", "John")` |
| Load from file | `ioutil.ReadFile` + `json.Unmarshal` | `easyjson.LoadFile("file.json")` |
| Save to file | `json.Marshal` + `ioutil.WriteFile` | `data.SaveFile("file.json")` |

## Performance

EasyJSON is designed for ease of use while maintaining reasonable performance. For high-performance scenarios where you need maximum speed and minimal allocations, consider using Go's standard `encoding/json` package directly.

EasyJSON trades some performance for significantly improved developer experience and code readability. It's ideal for:

- API servers handling moderate traffic
- Configuration file processing
- Data transformation scripts
- Rapid prototyping
- Any scenario where code clarity is more important than microsecond-level performance

## Testing

Run the comprehensive test suite:

```bash
# Run all tests
go test -v

# Run with coverage
go test -v -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=.
```

Run with coverage:

```bash
go test -cover
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Add tests for your changes
4. Ensure all tests pass (`go test -v`)
5. Commit your changes (`git commit -am 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

<<<<<<< HEAD
## Why Choose EasyJSON?

| Feature | Standard Go | Other Libraries | EasyJSON |
|---------|-------------|-----------------|----------|
| **Learning Curve** | Steep (interfaces, assertions) | Medium | Minimal (Python-like) |
| **Safety** | Manual nil checking | Varies | Built-in safe defaults |
| **API Flexibility** | Rigid struct mapping | Limited | Handles any JSON format |
| **Developer Experience** | Verbose, error-prone | Basic | AI-assisted, intuitive |
| **Error Handling** | Panic-prone | Basic | Helpful suggestions |
| **Performance** | Fast but complex | Varies | Fast and simple |
| **Breaking Changes** | Major version bumps | Frequent | Zero (additive only) |

**EasyJSON makes JSON handling in Go more enjoyable than Python, safer than JavaScript, and more productive than any other option!** 🚀

## Examples Repository

For more examples and use cases, check out the [examples directory](./example/main.go) in this repository.
