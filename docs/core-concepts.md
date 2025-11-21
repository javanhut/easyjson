# Core Concepts

Understanding the fundamental concepts behind EasyJSON will help you use the library effectively.

## The JSONValue Type

The `JSONValue` is the core type in EasyJSON. It wraps Go's `interface{}` and provides a rich set of methods for safe JSON manipulation.

### What is JSONValue?

```go
type JSONValue struct {
    data interface{}
}
```

A `JSONValue` can represent any valid JSON type:
- Object (map)
- Array (slice)
- String
- Number (int, float64)
- Boolean
- Null

### Creating JSONValues

```go
// From parsing
data, _ := easyjson.Loads(`{"name": "John"}`)

// Empty structures
obj := easyjson.NewObject()  // Empty JSON object
arr := easyjson.NewArray()   // Empty JSON array

// From Go values
obj := easyjson.NewObjectFrom(map[string]interface{}{
    "name": "John",
    "age":  30,
})

arr := easyjson.NewArrayFrom([]interface{}{"a", "b", "c"})

// From any value
data := easyjson.New(anyValue)
```

## Access Patterns

EasyJSON provides four different ways to access nested data. Choose the one that fits your use case.

### 1. Traditional Chaining (Get)

Step-by-step access, good for debugging:

```go
hairColor := data.Get("users").Get(0).Get("profile").Get("hair_color").AsString()
```

**Pros:**
- Familiar to Go developers
- Easy to debug (inspect each step)
- Explicit

**Cons:**
- Verbose for deep nesting
- Returns nil JSONValue on any missing key

### 2. Fluent Query (Q)

Most Python-like, recommended for most cases:

```go
hairColor := data.Q("users", 0, "profile", "hair_color").AsString()
```

**Pros:**
- Clean and concise
- Supports mixed types (strings and integers)
- Most flexible

**Cons:**
- Less explicit than chaining

**Best For:** Most use cases, especially with mixed access patterns

### 3. Path Notation

String-based paths with dot notation:

```go
hairColor := data.Path("users.0.profile.hair_color").AsString()
```

**Pros:**
- Very concise
- Good for configuration paths
- Easy to store paths as strings

**Cons:**
- Only supports strings (array indices as "0", "1", etc.)
- Less type-safe

**Best For:** Configuration files, stored paths

### 4. Smart Getters

Safe access with default values:

```go
name := data.GetString("user", "name", "Anonymous")
age := data.GetInt("user", "age", 0)
active := data.GetBool("user", "active", false)
```

**Pros:**
- No nil checking needed
- Built-in defaults
- Very safe

**Cons:**
- Default must match expected type

**Best For:** When you need guaranteed non-null values

## Safety Guarantees

EasyJSON prioritizes safety and never panics in normal operation.

### Safe Type Conversions

All `As*()` methods return sensible defaults instead of panicking:

```go
// If value is not a string, returns ""
str := data.Get("name").AsString()

// If value is not a number, returns 0
num := data.Get("age").AsInt()

// If value is not a boolean, returns false
flag := data.Get("active").AsBool()

// If value is not an array, returns empty slice
arr := data.Get("items").AsArray()

// If value is not an object, returns empty map
obj := data.Get("config").AsObject()
```

### Null Safety

Missing keys and null values are handled gracefully:

```go
// Both return JSONValue with nil data
missingKey := data.Get("nonexistent")
nullValue := data.Get("nullField")

// Check for null
if data.Get("field").IsNull() {
    fmt.Println("Field is null or missing")
}

// Safe access with defaults
value := data.GetString("field", "default")
```

### Chain Safety

Chains never panic, even with missing keys:

```go
// Safe even if any intermediate key is missing
value := data.Get("a").Get("b").Get("c").Get("d").AsString()
// Returns "" if any key in chain is missing
```

## Type System

### Type Checking

Always check types before conversion when type matters:

```go
if data.IsString() {
    str := data.AsString()
}

if data.IsNumber() {
    num := data.AsInt()
}

if data.IsArray() {
    for i := 0; i < data.Len(); i++ {
        item := data.Get(i)
        // Process item
    }
}

if data.IsObject() {
    for _, key := range data.Keys() {
        value := data.Get(key)
        // Process key-value
    }
}
```

### Available Type Checks

```go
data.IsNull()    // true if nil or missing
data.IsString()  // true if string
data.IsNumber()  // true if any numeric type
data.IsBool()    // true if boolean
data.IsArray()   // true if array/slice
data.IsObject()  // true if object/map
```

### Type String

Get human-readable type name:

```go
typeStr := data.TypeString()
// Returns: "null", "string", "number", "boolean", "array", or "object"
```

## Modification

EasyJSON supports both immutable and mutable operations.

### Immutable Operations

These return new JSONValues:

```go
// Filter returns new array
filtered := data.FilterArray(predicate)

// Map returns new array
mapped := data.MapArray(transform)

// Clone creates deep copy
clone := data.Clone()

// Sort returns new sorted array
sorted := data.SortBy("name")
```

### Mutable Operations

These modify the existing JSONValue:

```go
// Set modifies in place
data.Set("name", "John")

// SetPath creates intermediate objects
data.SetPath("user.profile.email", "john@example.com")

// Append modifies array
data.Append("new item")

// Delete removes key/index
data.Delete("oldField")

// Update merges objects
data.Update(otherJSONValue)
```

## Error Handling Philosophy

EasyJSON follows these principles:

### 1. Never Panic

Operations return safe defaults instead of panicking:

```go
// Never panics, even with invalid JSON
result := easyjson.ParseSafely(brokenJSON)
data := result.Data // Always valid JSONValue
```

### 2. Errors Where Appropriate

Operations that can genuinely fail return errors:

```go
// Parsing can fail
data, err := easyjson.Loads(jsonString)
if err != nil {
    // Handle parse error
}

// File operations can fail
err := data.SaveFile("output.json")
if err != nil {
    // Handle file error
}
```

### 3. Validation With Limits

All operations are validated with security limits:

```go
// Path validation
maxPathLength := 1000
maxPathDepth := 50

// Array validation
maxArraySize := 1000000
maxArrayIndex := 999999

// Memory validation
maxStringLength := 10 * 1024 * 1024  // 10MB
maxJSONSize := 100 * 1024 * 1024     // 100MB
```

These limits prevent DoS attacks and memory exhaustion.

## Working with Collections

### Objects (Maps)

```go
// Create
obj := easyjson.NewObject()

// Add fields
obj.Set("name", "John")
obj.Set("age", 30)

// Get all keys
keys := obj.Keys() // []string

// Get all values
values := obj.Values() // []*JSONValue

// Get key-value pairs
items := obj.Items() // map[string]*JSONValue

// Check key existence
exists := obj.Has("name") // bool

// Get length
count := obj.Len() // int
```

### Arrays (Slices)

```go
// Create
arr := easyjson.NewArray()

// Add items
arr.Append("item1")
arr.Append("item2")

// Access by index
first := arr.Get(0)

// Get length
length := arr.Len()

// Iterate
for i := 0; i < arr.Len(); i++ {
    item := arr.Get(i)
    // Process item
}

// Or use ForEach
arr.ForEach(func(i int, item *easyjson.JSONValue) {
    fmt.Printf("%d: %s\n", i, item.AsString())
})
```

## Raw Access

Sometimes you need the underlying Go value:

```go
// Get raw interface{} value
raw := data.Raw()

// Use in standard Go operations
switch v := raw.(type) {
case map[string]interface{}:
    // Handle object
case []interface{}:
    // Handle array
case string:
    // Handle string
case float64:
    // Handle number
case bool:
    // Handle boolean
case nil:
    // Handle null
}
```

## Performance Considerations

### Efficient Operations

```go
// Efficient: Single parse
data, _ := easyjson.Loads(jsonString)

// Efficient: Direct access
value := data.Q("user", "name")

// Efficient: Reuse JSONValue
for i := 0; i < 100; i++ {
    process(data) // Don't re-parse
}
```

### Less Efficient Operations

```go
// Less efficient: Multiple parses
for _, str := range jsonStrings {
    data, _ := easyjson.Loads(str) // Parse in loop
}

// Less efficient: Multiple conversions
for i := 0; i < arr.Len(); i++ {
    // Converting to string on each iteration
    str := arr.Get(i).Dumps()
}
```

## Best Practices

### 1. Choose the Right Access Pattern

```go
// Configuration: Use Path
dbURL := config.Path("database.url").AsString()

// Mixed access: Use Q
user := data.Q("users", 0, "profile", "name")

// With defaults: Use Smart Getters
name := data.GetString("name", "Unknown")
```

### 2. Check Types When Needed

```go
// When type is uncertain
if data.IsArray() {
    processArray(data)
} else if data.IsObject() {
    processObject(data)
}

// When type is known (from API spec)
users := data.Get("users") // Known to be array
```

### 3. Use Safe Parsing

```go
// For user input
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    showError(result.Error, result.Suggestions)
}

// For trusted input
data, err := easyjson.Loads(trustedInput)
if err != nil {
    return err
}
```

### 4. Prefer Immutable Operations

```go
// Good: Create new values
filtered := users.FilterArray(isActive)
sorted := filtered.SortBy("name")

// Avoid: Excessive mutation
data.Set("field1", val1)
data.Set("field2", val2)
// Better: Build with builder
data := easyjson.NewBuilder().
    AddField("field1", val1).
    AddField("field2", val2).
    ToJSON()
```

## Next Steps

- Learn about [Array Operations](array-operations.md) for powerful data manipulation
- Explore [Smart Getters](smart-getters.md) for safe access patterns
- Review [Best Practices](best-practices.md) for production use
- See [Examples](examples.md) for real-world patterns
