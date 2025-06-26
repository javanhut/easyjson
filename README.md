# EasyJSON - Python-like JSON handling for Go

EasyJSON provides an intuitive, Python-like interface for working with JSON data in Go. It eliminates the complexity of type assertions and provides safe, chainable operations on JSON structures.

## Features

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
    
    // Save to file
    data.SaveFileIndent("output.json", "  ")
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
go test -v
```

Run benchmarks:

```bash
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

## Examples Repository

For more examples and use cases, check out the [examples directory](./example/main.go) in this repository.