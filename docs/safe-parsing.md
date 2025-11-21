# Safe Parsing

EasyJSON provides multiple parsing strategies to handle invalid or messy JSON gracefully.

## ParseSafely - Never Fails

The safest parsing method that always returns valid data:

```go
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    fmt.Printf("Parse error: %v\n", result.Error)
    for _, suggestion := range result.Suggestions {
        fmt.Printf("Suggestion: %s\n", suggestion)
    }
}
data := result.Data // Always valid JSONValue, even on error
```

## TryParse - Boolean Result

Quick check if parsing succeeds:

```go
if data, ok := easyjson.TryParse(jsonString); ok {
    // Successfully parsed
    processData(data)
} else {
    // Parse failed
    handleError()
}
```

## ParseOrDefault - Fallback Value

Parse with a fallback:

```go
data := easyjson.ParseOrDefault(jsonString, easyjson.NewObject())
// Uses default if parsing fails
```

## ParseLenient - Very Forgiving

Extremely lenient parsing that tries multiple strategies:

```go
// Handles messy JSON with common issues
data := easyjson.ParseLenient(messyJSON)
```

## FixCommonIssues - Auto-Fix

Automatically fix common JSON issues:

```go
fixed := easyjson.FixCommonIssues(pythonStyleJSON)
// Fixes:
// - True/False -> true/false
// - None -> null
// - Single quotes -> double quotes
```

## ParseWithFixes - Parse After Fixing

Attempt fixes before parsing:

```go
data, err := easyjson.ParseWithFixes(messyJSON)
```

## Validation

### ValidateJSON

Check if string is valid JSON:

```go
if easyjson.ValidateJSON(jsonString) {
    // Valid JSON
}
```

### ValidateJSONWithDetails

Get detailed validation information:

```go
valid, err, suggestions := easyjson.ValidateJSONWithDetails(jsonString)
if !valid {
    fmt.Printf("Error: %v\n", err)
    fmt.Printf("Suggestions: %v\n", suggestions)
}
```

## Best Practices

### For User Input

```go
// Always use ParseSafely for untrusted input
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    return fmt.Errorf("invalid JSON: %w", result.Error)
}
```

### For Configuration Files

```go
// Use LoadFile with error handling
data, err := easyjson.LoadFile("config.json")
if err != nil {
    log.Fatal("Failed to load config:", err)
}
```

### For API Responses

```go
// ParseLenient for tolerating minor issues
data := easyjson.ParseLenient(apiResponse)
```

See [Getting Started](getting-started.md) for more examples.
