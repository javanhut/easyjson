# Troubleshooting

Common issues and their solutions when using EasyJSON.

## Parsing Issues

### Issue: Parse Error on Valid JSON

**Problem:**
```go
data, err := easyjson.Loads(jsonString)
// err: "unexpected end of JSON input"
```

**Solutions:**
1. Check for truncated JSON
2. Verify string encoding (UTF-8)
3. Look for unescaped quotes
4. Use `ParseSafely` for better error messages

```go
result := easyjson.ParseSafely(jsonString)
if result.Error != nil {
    fmt.Printf("Error: %v\n", result.Error)
    for _, suggestion := range result.Suggestions {
        fmt.Printf("Suggestion: %s\n", suggestion)
    }
}
```

### Issue: Python-Style JSON Won't Parse

**Problem:**
```json
{"active": True, "user": None}
```

**Solution:**
Use `FixCommonIssues` or `ParseWithFixes`:

```go
fixed := easyjson.FixCommonIssues(pythonJSON)
data, err := easyjson.Loads(fixed)

// Or
data, err := easyjson.ParseWithFixes(pythonJSON)
```

### Issue: Single Quotes in JSON

**Problem:**
```json
{'name': 'John'}
```

**Solution:**
```go
fixed := easyjson.FixCommonIssues(jsonString)
data, _ := easyjson.Loads(fixed)
```

## Access Issues

### Issue: Nil Pointer on Chained Access

**Problem:**
```go
name := data.Get("user").Get("profile").Get("name").AsString()
// Panics or returns empty if any key missing
```

**Solution:**
Use smart getters:

```go
name := data.GetString("user", "profile", "name", "Unknown")
```

### Issue: Wrong Type Returned

**Problem:**
```go
age := data.Get("age").AsInt() // Returns 0 unexpectedly
```

**Solutions:**
1. Check if value is actually a number
2. Verify JSON structure
3. Use type checking

```go
ageVal := data.Get("age")
if !ageVal.IsNumber() {
    fmt.Printf("Age is not a number, it's: %s\n", ageVal.TypeString())
}
```

### Issue: Array Index Out of Bounds

**Problem:**
```go
first := data.Get("items").Get(0) // Returns null if not array or empty
```

**Solution:**
Check length first:

```go
items := data.Get("items")
if items.IsArray() && items.Len() > 0 {
    first := items.Get(0)
}
```

## Path Issues

### Issue: Path Not Working

**Problem:**
```go
value := data.Path("user.profile.name") // Returns null
```

**Solutions:**
1. Verify path syntax (use dots, not slashes)
2. Check if keys exist
3. Use validation

```go
// Debug path
parts := []string{"user", "profile", "name"}
current := data
for _, part := range parts {
    if current.Has(part) {
        current = current.Get(part)
    } else {
        fmt.Printf("Missing key: %s\n", part)
        break
    }
}
```

### Issue: Array Index in Path

**Problem:**
```go
value := data.Path("items[0]") // Wrong syntax
```

**Solution:**
Use numeric index without brackets:

```go
value := data.Path("items.0") // Correct
// Or
value := data.Q("items", 0) // Better for arrays
```

## Building Issues

### Issue: Nested Object Not Created

**Problem:**
```go
data := easyjson.NewObject()
data.SetPath("user.profile.name", "John")
// profile object not created
```

**Solution:**
SetPath creates intermediate objects automatically for most cases, but for arrays you need to create them first:

```go
// For objects (works automatically)
data.SetPath("user.profile.name", "John") // OK

// For arrays (create first)
data.Set("items", []interface{}{})
data.SetPath("items.0", "first item")
```

### Issue: Builder Not Producing Expected JSON

**Problem:**
```go
builder := easyjson.NewBuilder().
    AddField("name", "John")
// Where's my JSON?
```

**Solution:**
Call `ToJSON()` to get the JSONValue:

```go
jsonValue := builder.ToJSON()
jsonString, _ := jsonValue.Dumps()
```

## Type Issues

### Issue: Number Always Returns 0

**Problem:**
```go
age := data.Get("age").AsInt() // Always 0
```

**Solutions:**
1. JSON numbers are float64 by default
2. Check actual type

```go
ageVal := data.Get("age")
fmt.Printf("Type: %s, Value: %v\n", ageVal.TypeString(), ageVal.Raw())

// If it's a float
age := int(ageVal.AsFloat())

// Or use smart getter
age := data.GetInt("age", 0)
```

### Issue: Boolean Conversion Unexpected

**Problem:**
```go
active := data.Get("active").AsBool() // Returns false for "true"
```

**Solution:**
Check if value is actually boolean:

```go
activeVal := data.Get("active")
if activeVal.IsString() && activeVal.AsString() == "true" {
    active = true
} else {
    active = activeVal.AsBool()
}

// Or use smart getter
active := data.GetBool("active", false)
```

## Performance Issues

### Issue: Slow JSON Processing

**Solutions:**
1. Parse once, reuse many times
2. Cache frequently accessed values
3. Avoid deep searches

```go
// Bad
for i := 0; i < 1000; i++ {
    data, _ := easyjson.Loads(jsonString) // Parsing 1000 times
}

// Good
data, _ := easyjson.Loads(jsonString) // Parse once
for i := 0; i < 1000; i++ {
    processData(data)
}
```

### Issue: High Memory Usage

**Solutions:**
1. Don't clone unnecessarily
2. Process large arrays in chunks
3. Use streaming for very large files

```go
// Instead of loading entire large file
items := data.Get("items")
items.ForEach(func(i int, item *easyjson.JSONValue) {
    processItem(item) // Process one at a time
})
```

## Validation Issues

### Issue: Validation Errors

**Problem:**
```go
// Error: "path exceeds maximum length"
```

**Solution:**
Check against limits:
- Max path length: 1,000 characters
- Max path depth: 50 levels
- Max array size: 1,000,000 elements

### Issue: Key Too Long Error

**Solution:**
Limit key lengths to 500 characters or restructure data.

## File I/O Issues

### Issue: File Not Found

**Problem:**
```go
data, err := easyjson.LoadFile("config.json")
// err: "no such file or directory"
```

**Solutions:**
1. Check file path (relative vs absolute)
2. Verify file exists
3. Check permissions

```go
// Use absolute path
absPath, _ := filepath.Abs("config.json")
data, err := easyjson.LoadFile(absPath)
```

### Issue: Permission Denied on Save

**Solution:**
Check write permissions and directory exists:

```go
// Ensure directory exists
dir := filepath.Dir(outputPath)
os.MkdirAll(dir, 0755)

// Save file
err := data.SaveFile(outputPath)
```

## Common Mistakes

### Mistake 1: Ignoring Errors

```go
// Bad
data, _ := easyjson.Loads(jsonString)

// Good
data, err := easyjson.Loads(jsonString)
if err != nil {
    return fmt.Errorf("parse failed: %w", err)
}
```

### Mistake 2: Assuming Type

```go
// Bad: Assuming it's an array
for i := 0; i < data.Get("items").Len(); i++ {
    // Len() returns 0 if not array
}

// Good: Check first
items := data.Get("items")
if items.IsArray() {
    for i := 0; i < items.Len(); i++ {
        // ...
    }
}
```

### Mistake 3: Not Using Smart Getters

```go
// Bad: Manual default handling
name := data.Get("name").AsString()
if name == "" {
    name = "Guest"
}

// Good: Smart getter
name := data.GetString("name", "Guest")
```

### Mistake 4: Mutating Shared Data

```go
// Bad: Modifying original
func processData(data *easyjson.JSONValue) {
    data.Set("processed", true) // Modifies original!
}

// Good: Clone first
func processData(data *easyjson.JSONValue) *easyjson.JSONValue {
    processed := data.Clone()
    processed.Set("processed", true)
    return processed
}
```

## Getting Help

If you're still stuck:

1. **Check Documentation**
   - [Getting Started](getting-started.md)
   - [API Reference](api-reference.md)
   - [Examples](examples.md)

2. **Enable Debug Output**
```go
result := easyjson.ParseSafely(input)
if result.Error != nil {
    fmt.Printf("Parse error: %v\n", result.Error)
    fmt.Printf("Suggestions: %v\n", result.Suggestions)
    fmt.Printf("Input (first 100 chars): %s\n", input[:min(100, len(input))])
}
```

3. **Inspect Data Structure**
```go
summary := data.GetSummary()
fmt.Printf("Summary: %+v\n", summary)

// Check available paths
smart := easyjson.WithSuggestions(data)
paths := smart.SuggestPaths()
fmt.Printf("Available paths: %v\n", paths)
```

4. **Create Minimal Reproduction**
```go
package main

import (
    "fmt"
    "github.com/javanhut/easyjson"
)

func main() {
    jsonStr := `{"test": "value"}`
    data, err := easyjson.Loads(jsonStr)
    if err != nil {
        panic(err)
    }
    
    // Your issue here
    fmt.Println(data.GetString("test"))
}
```

5. **Open an Issue**
   - Go to GitHub Issues
   - Include minimal reproduction
   - Describe expected vs actual behavior
   - Include EasyJSON version

## Quick Fixes

| Issue | Quick Fix |
|-------|-----------|
| Parse fails | Use `ParseSafely` for detailed errors |
| Nil panic | Use smart getters (`GetString`, etc.) |
| Wrong type | Check with `IsString()`, `IsNumber()`, etc. |
| Path not found | Use `TryPaths` for multiple options |
| Slow performance | Cache parsed data, avoid deep searches |
| Memory issues | Process in chunks, don't clone unnecessarily |

## Summary

Most issues can be avoided by:
1. Always handling errors
2. Using smart getters for optional fields
3. Checking types before conversion
4. Using `ParseSafely` for better error messages
5. Following [Best Practices](best-practices.md)
