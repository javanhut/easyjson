# Migration Guide

Guide for migrating to EasyJSON or upgrading between versions.

## From Standard encoding/json

### Before (encoding/json)

```go
var data map[string]interface{}
json.Unmarshal([]byte(jsonStr), &data)

user := data["user"].(map[string]interface{})
name := user["name"].(string)
age := user["age"].(float64)
```

### After (EasyJSON)

```go
data, _ := easyjson.Loads(jsonStr)

name := data.GetString("user", "name", "")
age := data.GetInt("user", "age", 0)
```

## From Other JSON Libraries

### From gjson

```go
// gjson
name := gjson.Get(jsonStr, "user.name").String()

// EasyJSON
data, _ := easyjson.Loads(jsonStr)
name := data.Path("user.name").AsString()
```

### From jsonparser

```go
// jsonparser
value, _ := jsonparser.GetString(data, "user", "name")

// EasyJSON
data, _ := easyjson.Load(jsonBytes)
value := data.GetString("user", "name", "")
```

## Version Upgrades

### From v1.x to v2.x

EasyJSON v2.0 introduces many new features with **zero breaking changes**.

All v1.x code continues to work:

```go
// v1.x code still works in v2.x
data, _ := easyjson.Loads(jsonStr)
name := data.Get("user").Get("name").AsString()
```

New v2.x features to adopt:

```go
// Smart getters (new in v2.x)
name := data.GetString("user", "name", "Unknown")

// Array operations (new in v2.x)
active := users.FilterArray(func(u *easyjson.JSONValue) bool {
    return u.GetBool("active", false)
})

// Builder API (new in v2.x)
response := easyjson.NewBuilder().
    AddAPIStatus("success", "OK").
    ToJSON()

// Multi-path access (new in v2.x)
id := data.TryPaths("id", "user_id", "userId").AsInt()

// Safe parsing (new in v2.x)
result := easyjson.ParseSafely(userInput)
```

## Migration Strategies

### Gradual Migration

Migrate incrementally without breaking existing code:

```go
// Phase 1: Replace parsing
// Old: json.Unmarshal
data, _ := easyjson.Loads(jsonStr)

// Phase 2: Adopt smart getters for new code
name := data.GetString("user", "name", "Guest")

// Phase 3: Refactor existing access patterns over time
// Old: data.Get("user").Get("name").AsString()
// New: data.GetString("user", "name", "")
```

### Testing During Migration

```go
func TestMigration(t *testing.T) {
    jsonStr := `{"user": {"name": "John"}}`
    
    // Old way still works
    data1, _ := easyjson.Loads(jsonStr)
    name1 := data1.Get("user").Get("name").AsString()
    
    // New way
    data2, _ := easyjson.Loads(jsonStr)
    name2 := data2.GetString("user", "name", "")
    
    // Both produce same result
    if name1 != name2 {
        t.Error("Migration failed")
    }
}
```

## Common Patterns

### Pattern 1: Type Assertions

```go
// Before
if user, ok := data["user"].(map[string]interface{}); ok {
    if name, ok := user["name"].(string); ok {
        fmt.Println(name)
    }
}

// After
name := data.GetString("user", "name", "")
if name != "" {
    fmt.Println(name)
}
```

### Pattern 2: Arrays

```go
// Before
if users, ok := data["users"].([]interface{}); ok {
    for _, u := range users {
        if user, ok := u.(map[string]interface{}); ok {
            name := user["name"].(string)
            fmt.Println(name)
        }
    }
}

// After
users := data.Get("users")
users.ForEach(func(i int, user *easyjson.JSONValue) {
    fmt.Println(user.GetString("name"))
})
```

### Pattern 3: Building JSON

```go
// Before
response := map[string]interface{}{
    "status": "success",
    "data": map[string]interface{}{
        "user": map[string]interface{}{
            "id":   123,
            "name": "John",
        },
    },
}
jsonBytes, _ := json.Marshal(response)

// After
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddObject("data", func(d *easyjson.JSONBuilder) {
        d.AddObject("user", func(u *easyjson.JSONBuilder) {
            u.AddField("id", 123).
              AddField("name", "John")
        })
    }).
    ToJSON()
jsonBytes, _ := response.Dump()
```

## Compatibility Notes

- EasyJSON is fully backward compatible within major versions
- No breaking changes in v2.x from v1.x
- All v1.x methods remain available
- New features are additive only
- Performance characteristics remain similar
- Memory usage is comparable to standard library

## Getting Help

If you encounter issues during migration:
- Check the [API Reference](api-reference.md)
- Review [Examples](examples.md)
- See [Troubleshooting](troubleshooting.md)
- Open an issue on GitHub
