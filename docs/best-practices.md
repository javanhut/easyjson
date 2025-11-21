# Best Practices

Guidelines for using EasyJSON effectively in production applications.

## Parsing

### Use the Right Parsing Method

```go
// User input: Use ParseSafely
result := easyjson.ParseSafely(userInput)

// Trusted sources: Use Loads with error handling
data, err := easyjson.Loads(trustedJSON)

// Configuration files: Use LoadFile
config, err := easyjson.LoadFile("config.json")
```

### Always Handle Parse Errors

```go
// Good
data, err := easyjson.Loads(jsonString)
if err != nil {
    return fmt.Errorf("parse failed: %w", err)
}

// Avoid
data, _ := easyjson.Loads(jsonString) // Ignoring errors
```

## Access Patterns

### Choose the Right Access Method

```go
// Known structure: Use Q
email := data.Q("user", "profile", "email").AsString()

// Optional fields: Use smart getters
name := data.GetString("user", "name", "Guest")

// Varying formats: Use TryPaths
id := data.TryPaths("id", "user_id", "userId").AsInt()

// Configuration: Use Path
dbURL := config.Path("database.url").AsString()
```

### Prefer Smart Getters

```go
// Good: Smart getter with default
port := config.GetInt("server", "port", 8080)

// Avoid: Manual nil checking
port := config.Q("server", "port").AsInt()
if port == 0 {
    port = 8080
}
```

## Type Safety

### Check Types When Uncertain

```go
// Good: Check before conversion
if data.IsArray() {
    items := data.AsArray()
    for _, item := range items {
        process(item)
    }
}

// Avoid: Assuming type
items := data.AsArray() // Could be empty if not array
```

### Use Type-Specific Methods

```go
// Good: Specific type method
age := data.GetInt("age", 0)

// Less good: Generic then convert
age := data.Get("age").AsInt()
```

## Performance

### Reuse Parsed Data

```go
// Good: Parse once, use many times
data, _ := easyjson.Loads(jsonString)
processUser(data.Get("user"))
processSettings(data.Get("settings"))
processPreferences(data.Get("preferences"))

// Avoid: Parsing repeatedly
processUser(easyjson.Loads(jsonString).Get("user"))
processSettings(easyjson.Loads(jsonString).Get("settings"))
```

### Cache Frequently Accessed Paths

```go
// Good: Cache frequently used values
users := data.Get("users")
for i := 0; i < users.Len(); i++ {
    processUser(users.Get(i))
}

// Less efficient: Re-accessing each time
for i := 0; i < data.Get("users").Len(); i++ {
    processUser(data.Get("users").Get(i))
}
```

## Building JSON

### Use Builders for Complex Structures

```go
// Good: Use builder
response := easyjson.NewBuilder().
    AddAPIStatus("success", "OK").
    AddObject("data", func(d *easyjson.JSONBuilder) {
        d.AddField("users", users)
    }).
    ToJSON()

// Avoid: Manual building for complex structures
response := easyjson.NewObject()
response.Set("status", "success")
// ... many more Set calls
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

## Error Handling

### Don't Ignore Errors

```go
// Good
err := data.SaveFile("output.json")
if err != nil {
    return fmt.Errorf("save failed: %w", err)
}

// Bad
data.SaveFile("output.json") // Ignoring error
```

### Provide Context in Errors

```go
// Good: Contextual error
if err := data.SaveFile(path); err != nil {
    return fmt.Errorf("failed to save config to %s: %w", path, err)
}

// Less helpful
if err != nil {
    return err
}
```

## Security

### Validate User Input

```go
// Always validate untrusted input
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    return errors.New("invalid JSON")
}

// Validate structure
if !result.Data.HasRequiredFields("name", "email") {
    return errors.New("missing required fields")
}
```

### Sanitize Output

```go
// Remove sensitive data before sending
publicData := userData.SanitizeForOutput()
json, _ := publicData.Dumps()
w.Write([]byte(json))
```

## Code Organization

### Extract Repeated Patterns

```go
// Good: Extract to helper
func buildUserResponse(user User) *easyjson.JSONValue {
    return easyjson.QuickObject(
        "id", user.ID,
        "name", user.Name,
        "email", user.Email,
    )
}

// Use in multiple places
response1 := buildUserResponse(user1)
response2 := buildUserResponse(user2)
```

### Use Type-Safe Wrappers

```go
type Config struct {
    data *easyjson.JSONValue
}

func (c *Config) GetPort() int {
    return c.data.GetInt("server", "port", 8080)
}

func (c *Config) GetHost() string {
    return c.data.GetString("server", "host", "localhost")
}
```

## Testing

### Test with Various Inputs

```go
func TestParseUser(t *testing.T) {
    testCases := []struct {
        name  string
        input string
        want  string
    }{
        {"standard format", `{"name":"John"}`, "John"},
        {"empty", `{}`, "Guest"},
        {"null name", `{"name":null}`, "Guest"},
        {"invalid JSON", `{invalid}`, "Guest"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            data := easyjson.ParseSafely(tc.input).Data
            got := data.GetString("name", "Guest")
            if got != tc.want {
                t.Errorf("got %q, want %q", got, tc.want)
            }
        })
    }
}
```

## Documentation

### Document Default Values

```go
// Good: Clear why default is chosen
// Default port 8080 avoids requiring root privileges
port := config.GetInt("port", 8080)

// Default timeout based on average API response time
timeout := config.GetInt("timeout", 30)
```

### Document Expected Structure

```go
// ProcessUserData expects JSON in the format:
// {
//   "user": {
//     "id": 123,
//     "name": "John Doe",
//     "email": "john@example.com"
//   }
// }
func ProcessUserData(data *easyjson.JSONValue) error {
    // ...
}
```

## Common Pitfalls to Avoid

### Don't Assume Types

```go
// Bad: Assuming it's an array
for i := 0; i < data.Get("items").Len(); i++ {
    // Will be 0 if not an array
}

// Good: Check first
items := data.Get("items")
if items.IsArray() {
    for i := 0; i < items.Len(); i++ {
        // Process
    }
}
```

### Don't Chain Without Checking

```go
// Risky: Any missing key breaks chain
name := data.Get("user").Get("profile").Get("name").AsString()

// Better: Use smart getter
name := data.GetString("user", "profile", "name", "Unknown")
```

### Don't Mutate Shared Data

```go
// Bad: Mutating shared JSONValue
func updateUser(data *easyjson.JSONValue) {
    data.Set("updated_at", time.Now())
}

// Good: Clone first
func updateUser(data *easyjson.JSONValue) *easyjson.JSONValue {
    updated := data.Clone()
    updated.Set("updated_at", time.Now())
    return updated
}
```

## Summary

Key principles:
1. **Safety first**: Use ParseSafely, validate inputs, check types
2. **Choose wisely**: Pick the right access method for your use case
3. **Be efficient**: Reuse parsed data, cache frequently accessed values
4. **Handle errors**: Always check errors, provide context
5. **Stay secure**: Validate, sanitize, use safe defaults
6. **Keep it clean**: Extract patterns, use type-safe wrappers
7. **Test thoroughly**: Cover edge cases, invalid inputs
8. **Document well**: Explain defaults, expected structures

See related guides:
- [Security](validation-security.md) for security best practices
- [Performance](performance.md) for optimization tips
- [Examples](examples.md) for real-world patterns
