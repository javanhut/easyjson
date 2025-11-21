# Getting Started with EasyJSON

This guide will help you get up and running with EasyJSON quickly.

## Installation

### Requirements

- Go 1.19 or later
- No external dependencies required

### Install via go get

```bash
go get github.com/javanhut/easyjson
```

### Verify Installation

Create a simple test file to verify the installation:

```go
package main

import (
    "fmt"
    "github.com/javanhut/easyjson"
)

func main() {
    data, _ := easyjson.Loads(`{"message": "Hello, EasyJSON!"}`)
    fmt.Println(data.GetString("message"))
}
```

Run it:
```bash
go run test.go
```

You should see: `Hello, EasyJSON!`

## Quick Start

### Basic Parsing

Parse JSON from strings, bytes, or files:

```go
// From string
data, err := easyjson.Loads(`{"name": "John", "age": 30}`)
if err != nil {
    log.Fatal(err)
}

// From bytes
jsonBytes := []byte(`{"name": "John", "age": 30}`)
data, err := easyjson.Load(jsonBytes)

// From file
data, err := easyjson.LoadFile("config.json")
```

### Accessing Data

EasyJSON provides multiple ways to access data:

```go
jsonStr := `{
    "user": {
        "name": "Alice",
        "age": 30,
        "profile": {
            "email": "alice@example.com"
        }
    }
}`

data, _ := easyjson.Loads(jsonStr)

// Method 1: Traditional chaining
name := data.Get("user").Get("name").AsString()

// Method 2: Fluent Query (recommended)
email := data.Q("user", "profile", "email").AsString()

// Method 3: Path notation
age := data.Path("user.age").AsInt()

// Method 4: Smart getters with defaults
name := data.GetString("user", "name", "Unknown")
```

### Type Conversions

All type conversions are safe and return sensible defaults:

```go
// String conversion
name := data.Q("name").AsString() // Returns "" if not a string

// Number conversions
age := data.Q("age").AsInt()         // Returns 0 if not a number
price := data.Q("price").AsFloat()   // Returns 0.0 if not a number

// Boolean conversion
active := data.Q("active").AsBool()  // Returns false if not a boolean

// Collection conversions
users := data.Q("users").AsArray()   // Returns empty slice if not array
config := data.Q("config").AsObject() // Returns empty map if not object
```

### Type Checking

Check types before conversion:

```go
if data.Q("age").IsNumber() {
    age := data.Q("age").AsInt()
    fmt.Printf("Age: %d\n", age)
}

if data.Q("users").IsArray() {
    users := data.Q("users").AsArray()
    fmt.Printf("Found %d users\n", len(users))
}
```

### Modifying Data

```go
// Create a new object
user := easyjson.NewObject()

// Set values
user.Set("name", "John Doe")
user.Set("age", 30)
user.Set("active", true)

// Set nested paths (creates intermediate objects)
user.SetPath("profile.email", "john@example.com")

// Append to arrays
tags := easyjson.NewArray()
tags.Append("golang")
tags.Append("json")
user.Set("tags", tags.Raw())
```

### Saving Data

```go
// To JSON string
jsonString, err := data.Dumps()

// Pretty-printed JSON
prettyJSON, err := data.DumpsIndent("  ")

// To file
err := data.SaveFile("output.json")

// Pretty-printed to file
err := data.SaveFileIndent("output.json", "  ")
```

## Your First Complete Program

Here's a complete example that demonstrates common operations:

```go
package main

import (
    "fmt"
    "log"
    "github.com/javanhut/easyjson"
)

func main() {
    // Sample JSON data
    jsonData := `{
        "users": [
            {"id": 1, "name": "Alice", "role": "admin", "active": true},
            {"id": 2, "name": "Bob", "role": "user", "active": false},
            {"id": 3, "name": "Charlie", "role": "admin", "active": true}
        ],
        "metadata": {
            "total": 3,
            "page": 1
        }
    }`

    // Parse JSON
    data, err := easyjson.Loads(jsonData)
    if err != nil {
        log.Fatal("Failed to parse JSON:", err)
    }

    // Access nested data
    total := data.GetInt("metadata", "total", 0)
    fmt.Printf("Total users: %d\n", total)

    // Work with arrays
    users := data.Get("users")
    
    // Filter active users
    activeUsers := users.FilterArray(func(user *easyjson.JSONValue) bool {
        return user.GetBool("active", false)
    })
    fmt.Printf("Active users: %d\n", activeUsers.Len())

    // Find admin users
    admins := users.FindAllByField("role", "admin")
    fmt.Printf("Admin users: %d\n", len(admins))

    // Extract all names
    names := users.PluckStrings("name")
    fmt.Printf("User names: %v\n", names)

    // Create new data
    newUser := easyjson.NewBuilder().
        AddField("id", 4).
        AddField("name", "Dave").
        AddField("role", "user").
        AddField("active", true).
        ToJSON()

    // Add to array
    users.Append(newUser.Raw())

    // Save modified data
    if err := data.SaveFileIndent("users.json", "  "); err != nil {
        log.Fatal("Failed to save:", err)
    }

    fmt.Println("Data saved successfully!")
}
```

## Common Patterns

### Safe Parsing

Never let parsing errors crash your application:

```go
// Method 1: ParseSafely (never fails)
result := easyjson.ParseSafely(jsonString)
if result.Error != nil {
    fmt.Printf("Parse error: %v\n", result.Error)
    for _, suggestion := range result.Suggestions {
        fmt.Printf("Suggestion: %s\n", suggestion)
    }
}
data := result.Data // Always valid, even on error

// Method 2: Try parse
if data, ok := easyjson.TryParse(jsonString); ok {
    // Successfully parsed
}

// Method 3: Parse with fallback
data := easyjson.ParseOrDefault(jsonString, easyjson.NewObject())
```

### Working with API Responses

```go
func handleAPIResponse(jsonStr string) error {
    data, err := easyjson.Loads(jsonStr)
    if err != nil {
        return err
    }

    // Check API status
    if !data.GetBool("success", false) {
        return fmt.Errorf("API error: %s", data.GetString("error"))
    }

    // Extract data
    users := data.Get("data", "users")
    for i := 0; i < users.Len(); i++ {
        user := users.Get(i)
        fmt.Printf("User: %s (%s)\n", 
            user.GetString("name"),
            user.GetString("email"))
    }

    return nil
}
```

### Configuration Management

```go
type Config struct {
    data *easyjson.JSONValue
}

func LoadConfig(path string) (*Config, error) {
    data, err := easyjson.LoadFile(path)
    if err != nil {
        return nil, err
    }
    return &Config{data: data}, nil
}

func (c *Config) GetDatabaseURL() string {
    return c.data.GetString("database", "url", "sqlite://default.db")
}

func (c *Config) GetPort() int {
    return c.data.GetInt("server", "port", 8080)
}

func (c *Config) IsDebugMode() bool {
    return c.data.GetBool("debug", false)
}
```

## Next Steps

Now that you understand the basics, explore:

1. [Core Concepts](core-concepts.md) - Understand the philosophy and design
2. [API Reference](api-reference.md) - Complete method documentation
3. [Array Operations](array-operations.md) - Powerful array manipulation
4. [Examples](examples.md) - Real-world use cases

## Quick Reference

### Most Used Methods

```go
// Parsing
data, err := easyjson.Loads(jsonString)
data, err := easyjson.LoadFile("file.json")
result := easyjson.ParseSafely(jsonString)

// Access
value := data.Get("key")
value := data.Q("user", "profile", "name")
value := data.Path("user.profile.name")

// Safe access with defaults
str := data.GetString("key", "default")
num := data.GetInt("key", 0)
bool := data.GetBool("key", false)

// Type checking
data.IsString(), data.IsNumber(), data.IsBool()
data.IsArray(), data.IsObject(), data.IsNull()

// Type conversion
data.AsString(), data.AsInt(), data.AsFloat()
data.AsBool(), data.AsArray(), data.AsObject()

// Modification
data.Set("key", value)
data.SetPath("user.name", "John")
data.Delete("key")

// Arrays
data.Append(item)
data.FilterArray(predicate)
data.MapArray(transform)
data.FindByField("role", "admin")

// Serialization
jsonStr, err := data.Dumps()
prettyJSON, err := data.DumpsIndent("  ")
err := data.SaveFile("output.json")
```

## Getting Help

- Review the [Troubleshooting Guide](troubleshooting.md)
- Check out more [Examples](examples.md)
- Read the complete [API Reference](api-reference.md)
- Open an issue on GitHub
