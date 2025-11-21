# Performance Tips

Guidelines for optimizing EasyJSON performance in your applications.

## General Principles

### 1. Parse Once, Use Many Times

```go
// Good: Parse once
data, _ := easyjson.Loads(jsonString)
processUser(data)
processSettings(data)
processPreferences(data)

// Bad: Parse repeatedly
processUser(easyjson.Loads(jsonString))
processSettings(easyjson.Loads(jsonString))
```

### 2. Cache Frequently Accessed Values

```go
// Good: Cache in variable
users := data.Get("users")
for i := 0; i < users.Len(); i++ {
    processUser(users.Get(i))
}

// Less efficient: Re-access each time
for i := 0; i < data.Get("users").Len(); i++ {
    processUser(data.Get("users").Get(i))
}
```

### 3. Use Appropriate Access Methods

```go
// Fast: Direct access with known structure
name := data.Q("user", "name").AsString()

// Slower: Path parsing
name := data.Path("user.name").AsString()

// Slowest: Deep search
name := data.DeepSearch("name").AsString()
```

## Method Performance

### Access Methods (Fastest to Slowest)

1. **Get()** - Direct key access
2. **Q()** - Fluent query (slight overhead for variadic args)
3. **GetString/GetInt/etc()** - Smart getters (includes default handling)
4. **Path()** - String parsing overhead
5. **TryPaths()** - Multiple path attempts
6. **DeepSearch()** - Recursive search

### When to Use Each

```go
// Known, static structure: Use Get or Q
email := data.Q("user", "profile", "email")

// Optional fields with defaults: Use smart getters
port := config.GetInt("server", "port", 8080)

// Varying API formats: Use TryPaths
id := data.TryPaths("id", "user_id", "userId")

// Unknown structure: Use DeepSearch (slowest)
email := data.DeepSearch("email") // Only when necessary
```

## Array Operations

### Efficient Filtering

```go
// Good: Single pass
activeAdmins := users.FilterArray(func(u *easyjson.JSONValue) bool {
    return u.GetBool("active") && u.GetString("role") == "admin"
})

// Less efficient: Multiple passes
active := users.FilterArray(isActive)
admins := active.FilterArray(isAdmin)
```

### Avoid Unnecessary Conversions

```go
// Good: Work with JSONValue
result := users.FilterArray(predicate).SortBy("name")

// Bad: Convert to Go slice and back
userSlice := users.AsArray()
// ... Go operations ...
result := easyjson.NewArrayFrom(userSlice)
```

## Memory Optimization

### 1. Clone Only When Necessary

```go
// Good: Clone only if modifying
func updateUser(user *easyjson.JSONValue) *easyjson.JSONValue {
    updated := user.Clone()
    updated.Set("updated_at", time.Now())
    return updated
}

// Wasteful: Cloning when not needed
func readUser(user *easyjson.JSONValue) string {
    copy := user.Clone() // Unnecessary
    return copy.GetString("name")
}
```

### 2. Reuse Builders

```go
// Good: Reuse builder for similar structures
builder := easyjson.NewBuilder()
for _, user := range users {
    builder.Reset()
    response := builder.
        AddField("id", user.ID).
        AddField("name", user.Name).
        ToJSON()
    responses = append(responses, response)
}
```

### 3. Avoid Large Intermediate Structures

```go
// Good: Stream processing
func processLargeArray(data *easyjson.JSONValue) {
    items := data.Get("items")
    items.ForEach(func(i int, item *easyjson.JSONValue) {
        processItem(item)
    })
}

// Bad: Creating large intermediate arrays
func processLargeArray(data *easyjson.JSONValue) {
    allItems := data.Get("items").AsArray() // Converts all at once
    for _, item := range allItems {
        processItem(item)
    }
}
```

## Parsing Optimization

### Choose the Right Parser

```go
// Strict validation (slowest, safest)
result := easyjson.ParseSafely(input)

// Standard parsing (balanced)
data, err := easyjson.Loads(input)

// From bytes (skip string allocation)
data, err := easyjson.Load(jsonBytes)

// From file (optimized I/O)
data, err := easyjson.LoadFile(path)
```

### Size Matters

```go
// Pre-check size for large inputs
if len(input) > 10*1024*1024 { // 10MB
    return errors.New("input too large")
}
data, err := easyjson.Loads(input)
```

## Building Optimization

### Batch Operations

```go
// Good: Batch with AddFields
builder := easyjson.NewBuilder().
    AddFields(map[string]interface{}{
        "id":    123,
        "name":  "John",
        "email": "john@example.com",
    })

// Less efficient: Individual AddField calls
builder := easyjson.NewBuilder().
    AddField("id", 123).
    AddField("name", "John").
    AddField("email", "john@example.com")
```

### Use Quick Builders for Simple Cases

```go
// Good: Quick builder
user := easyjson.QuickObject("id", 123, "name", "John")

// Overkill: Full builder
user := easyjson.NewBuilder().
    AddField("id", 123).
    AddField("name", "John").
    ToJSON()
```

## Serialization

### Choose the Right Output Method

```go
// Fastest: Direct to bytes
bytes, _ := data.Dump()

// Fast: String without formatting
str, _ := data.Dumps()

// Slower: Pretty printing
pretty, _ := data.DumpsIndent("  ")

// Slowest: File I/O
data.SaveFileIndent("output.json", "  ")
```

### Reuse Serialized Data

```go
// Good: Serialize once
jsonBytes, _ := data.Dump()
for _, client := range clients {
    client.Send(jsonBytes)
}

// Bad: Serialize repeatedly
for _, client := range clients {
    jsonBytes, _ := data.Dump()
    client.Send(jsonBytes)
}
```

## Benchmarking

### Profile Your Code

```go
import "testing"

func BenchmarkDataAccess(b *testing.B) {
    data, _ := easyjson.Loads(testJSON)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = data.GetString("user", "name")
    }
}

// Run with:
// go test -bench=. -benchmem
```

### Common Bottlenecks

1. **Excessive parsing**: Cache parsed data
2. **Deep searches**: Use direct access when possible
3. **Large arrays**: Process in chunks or streams
4. **Unnecessary cloning**: Clone only when modifying
5. **String operations**: Use bytes when possible

## Production Tips

### 1. Connection Pooling

```go
// Reuse HTTP clients
var httpClient = &http.Client{
    Timeout: 10 * time.Second,
}

func fetchJSON(url string) (*easyjson.JSONValue, error) {
    resp, err := httpClient.Get(url)
    // ... parse response
}
```

### 2. Worker Pools

```go
// Process JSON in parallel
func processJSONFiles(files []string) {
    jobs := make(chan string, len(files))
    results := make(chan Result, len(files))
    
    // Start workers
    for w := 0; w < runtime.NumCPU(); w++ {
        go worker(jobs, results)
    }
    
    // Send jobs
    for _, file := range files {
        jobs <- file
    }
    close(jobs)
    
    // Collect results
    for range files {
        <-results
    }
}

func worker(jobs <-chan string, results chan<- Result) {
    for file := range jobs {
        data, err := easyjson.LoadFile(file)
        // Process data...
        results <- Result{data, err}
    }
}
```

### 3. Memory Limits

```go
// Enforce size limits
const MaxJSONSize = 10 * 1024 * 1024 // 10MB

func safeLoad(input string) (*easyjson.JSONValue, error) {
    if len(input) > MaxJSONSize {
        return nil, errors.New("input too large")
    }
    return easyjson.Loads(input)
}
```

## Summary

**Key Takeaways:**
1. Parse once, use many times
2. Cache frequently accessed values
3. Use appropriate access methods for your use case
4. Avoid unnecessary conversions and clones
5. Profile your code to find bottlenecks
6. Set reasonable size limits
7. Use worker pools for batch processing

**Performance Hierarchy (Fast → Slow):**
- Direct access (Get, Q)
- Smart getters (GetString, GetInt)
- Path notation (Path)
- Multi-path (TryPaths)
- Deep search (DeepSearch)

See [Best Practices](best-practices.md) for general usage guidelines.
