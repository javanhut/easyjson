# EasyJSON - Python-like JSON handling for Go

EasyJSON provides an intuitive, Python-like interface for working with JSON data in Go. It eliminates the complexity of type assertions and provides safe, chainable operations on JSON structures.

## Features

- **Python-like API**: Familiar `loads()`, `dumps()`, and intuitive access patterns
- **Fluent Query Syntax**: Chain access with `data.Q("users", 0, "profile", "hair_color")` 
- **Safe operations**: No panics on missing keys or invalid operations
- **Type flexibility**: Automatic type conversions with fallback defaults
- **Path notation**: Access nested values with dot notation (`data.Path("user.address.street")`)
- **Chainable operations**: Fluent interface for complex manipulations
- **Zero external dependencies**: Uses only Go standard library

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
    // Parse JSON string
    data, _ := easyjson.Loads(`{
        "name": "John Doe",
        "age": 30,
        "hobbies": ["reading", "coding"],
        "address": {
            "street": "123 Main St",
            "city": "New York"
        }
    }`)

    // Access values easily
    fmt.Println("Name:", data.Get("name").AsString())
    fmt.Println("Age:", data.Get("age").AsInt())
    fmt.Println("First hobby:", data.Get("hobbies").Get(0).AsString())
    
    // Use path notation for nested access
    fmt.Println("City:", data.Path("address.city").AsString())
    
    // Or use fluent query syntax (most Python-like)
    fmt.Println("City:", data.Q("address", "city").AsString())
    
    // Modify data
    data.Set("age", 31)
    data.Get("hobbies").Append("photography")
    
    // Convert back to JSON
    result, _ := data.DumpsIndent("  ")
    fmt.Println(result)
}
```

## API Reference

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

// Convert to JSON bytes
jsonBytes, err := data.Dump()
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

### Important Notes

**SetPath Limitations:**
- `SetPath` automatically creates intermediate **objects** when paths don't exist
- For arrays, you need to create the array structure first before using `SetPath` with numeric indices
- Example: Create `data.Set("items", []interface{}{})` before using `data.SetPath("items.0", value)`

### Accessing Data

```go
// Get values by key (objects) or index (arrays)
value := data.Get("key")
firstItem := data.Get(0)

// Fluent query syntax - most Python-like approach
hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
score := data.Q("players", 1, "stats", "score").AsInt()

// Check if key/index exists
exists := data.Has("key")
exists := data.Has(0)

// Access nested data with paths
street := data.Path("user.address.street")
score := data.Path("users.0.scores.1")

// Fluent query syntax (most Python-like)
hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
age := data.Q("users", 0, "age").AsInt()
```

### Modifying Data

```go
// Set values
data.Set("key", "value")
data.Set(0, "new first item")

// Set nested paths (creates intermediate objects)
data.SetPath("user.address.street", "456 Oak Ave")

// For arrays, create the structure first, then set elements
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

### Type Checking

```go
if data.IsString() { /* ... */ }
if data.IsNumber() { /* ... */ }
if data.IsBool() { /* ... */ }
if data.IsArray() { /* ... */ }
if data.IsObject() { /* ... */ }
if data.IsNull() { /* ... */ }
```

### Type Conversion

All conversion methods provide safe defaults for invalid conversions:

```go
str := data.AsString()    // Returns "" for non-strings
num := data.AsInt()       // Returns 0 for non-numbers
flt := data.AsFloat()     // Returns 0.0 for non-numbers
bln := data.AsBool()      // Returns false for non-bools
arr := data.AsArray()     // Returns empty slice for non-arrays
obj := data.AsObject()    // Returns empty map for non-objects
raw := data.Raw()         // Returns underlying Go value
```

### Collection Operations

```go
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
    // Handle missing or invalid data
    name = "Unknown"
}

// Check existence before access
if data.Has("optional_field") {
    value := data.Get("optional_field").AsString()
    // Process value
}

// Fluent query is always safe
missing := data.Q("nonexistent", "path", "here").AsString() // Returns ""
```

## Testing

Run the comprehensive test suite:

```bash
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

## Performance

EasyJSON is designed for ease of use while maintaining reasonable performance. For high-performance scenarios where you need maximum speed and minimal allocations, consider using Go's standard `encoding/json` package directly.

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

## Access Pattern Comparison

EasyJSON provides three different ways to access nested data - choose what feels most natural:

| Pattern | Syntax | Best For |
|---------|--------|----------|
| **Fluent Query** | `data.Q("users", 0, "name").AsString()` | Python-like access, mixed key types |
| **Path Notation** | `data.Path("users.0.name").AsString()` | String-based paths, simple cases |
| **Traditional** | `data.Get("users").Get(0).Get("name").AsString()` | Step-by-step access, debugging |

```go
// All three are equivalent:
name1 := data.Q("users", 0, "name").AsString()           // Fluent (recommended)
name2 := data.Path("users.0.name").AsString()            // Path notation
name3 := data.Get("users").Get(0).Get("name").AsString() // Traditional chaining
```

| Operation | Standard Go | EasyJSON |
|-----------|-------------|----------|
| Parse JSON | `json.Unmarshal(data, &v)` | `easyjson.Loads(jsonStr)` |
| Access nested | `v["user"].(map[string]interface{})["name"].(string)` | `data.Path("user.name").AsString()` |
| Type assertion | `value, ok := v.(string)` | `data.AsString()` (safe) |
| Check existence | Complex nested checks | `data.Has("key")` |
| Modify nested | Manual map/slice operations | `data.SetPath("user.name", "John")` |

EasyJSON trades some performance for significantly improved developer experience and code readability.
