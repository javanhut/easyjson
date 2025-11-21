# Validation and Security

EasyJSON includes comprehensive validation and security features to protect against common attacks and ensure safe operation.

## Built-in Security Limits

### Path Validation
- **Max Path Length**: 1,000 characters
- **Max Path Depth**: 50 levels
- **Max Path Segment**: 200 characters

### Array Validation
- **Max Array Size**: 1,000,000 elements
- **Max Array Index**: 999,999

### Key Validation
- **Max Key Length**: 500 characters
- **Max Object Keys**: 10,000 keys

### Memory Limits
- **Max String Length**: 10 MB
- **Max JSON Size**: 100 MB
- **Max Recursion Depth**: 100 levels

### Operation Limits
- **Max Batch Operations**: 10,000 operations

## Validation Features

### Automatic Validation

All operations are automatically validated:

```go
// Path validation
data.Path("user.profile.email") // Validates path format and depth

// Array validation
arr.Get(index) // Validates index bounds

// Key validation
obj.Set("key", value) // Validates key format and length
```

### Input Sanitization

```go
// Removes sensitive fields
publicData := data.SanitizeForOutput()
// Removes: password, secret, token, key, ssn, credit_card, etc.
```

## Security Best Practices

### 1. Always Validate User Input

```go
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    return fmt.Errorf("invalid JSON: %w", result.Error)
}
// Use result.Data safely
```

### 2. Limit JSON Size

```go
// Check size before parsing
if len(jsonString) > MaxJSONSize {
    return errors.New("JSON too large")
}

data, err := easyjson.Loads(jsonString)
```

### 3. Validate Data Structure

```go
// Validate required fields
if !data.HasRequiredFields("id", "name", "email") {
    return errors.New("missing required fields")
}

// Validate field formats
email := data.Get("email")
if !email.IsValidEmail() {
    return errors.New("invalid email format")
}
```

### 4. Sanitize Output

```go
// Remove sensitive data before sending
publicResponse := data.SanitizeForOutput()
json, _ := publicResponse.Dumps()
w.Write([]byte(json))
```

### 5. Use Safe Access Methods

```go
// Safe: Returns default if missing
name := data.GetString("user", "name", "Guest")

// Unsafe: Could access nil
name := data.Get("user").Get("name").AsString()
```

## Protection Against Common Attacks

### 1. JSON Bomb Protection

Limits prevent deeply nested or extremely large JSON:

```go
// Protected against deep nesting (max 50 levels)
// Protected against large arrays (max 1M elements)
// Protected against huge strings (max 10MB)
```

### 2. Path Traversal Protection

```go
// Validates paths to prevent directory traversal
data.Path("../../etc/passwd") // Returns error
```

### 3. DoS Protection

Memory and operation limits prevent resource exhaustion:

```go
// Max JSON size: 100MB
// Max recursion depth: 100
// Max batch operations: 10,000
```

### 4. Injection Protection

All keys and values are properly validated:

```go
// Control characters are rejected
// Invalid UTF-8 is rejected
// Excessively long keys are rejected
```

## Error Handling

### Safe Error Messages

Never expose internal details in production:

```go
result := easyjson.ParseSafely(input)
if result.Error != nil {
    // In production, log but don't expose details
    log.Printf("Parse error: %v", result.Error)
    
    // Return generic error to user
    return errors.New("invalid JSON format")
}
```

### Validation Errors

```go
type ValidationError struct {
    Type    string
    Message string
    Value   interface{}
}
```

## Configuration

### Custom Validation Limits

```go
config := easyjson.DefaultValidationConfig()
config.MaxPathLength = 500
config.MaxArraySize = 100000

// Use custom config with validation
// (Note: Most users should use defaults)
```

## Monitoring and Logging

### Log Validation Failures

```go
result := easyjson.ParseSafely(input)
if result.Error != nil {
    log.Printf("Validation failed: %v, Input size: %d", 
        result.Error, len(input))
}
```

### Monitor Resource Usage

```go
// Track JSON sizes being processed
size := len(jsonString)
if size > expectedMaxSize {
    log.Warnf("Unusually large JSON: %d bytes", size)
}
```

## Secure Configuration Loading

```go
func LoadSecureConfig(path string) (*Config, error) {
    // Validate file path
    if !isValidConfigPath(path) {
        return nil, errors.New("invalid config path")
    }
    
    // Check file size
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }
    if info.Size() > MaxConfigSize {
        return nil, errors.New("config file too large")
    }
    
    // Load and validate
    data, err := easyjson.LoadFile(path)
    if err != nil {
        return nil, err
    }
    
    // Validate structure
    if !data.HasRequiredFields("version", "server") {
        return nil, errors.New("invalid config structure")
    }
    
    return parseConfig(data), nil
}
```

## Summary

EasyJSON provides defense-in-depth security:
- Automatic validation of all inputs
- Comprehensive size and depth limits
- Protection against common attacks
- Safe defaults and error handling
- Memory-safe operations

Always:
- Use `ParseSafely` for untrusted input
- Validate data structure after parsing
- Sanitize output with `SanitizeForOutput`
- Use smart getters with defaults
- Monitor for unusual patterns

See [Best Practices](best-practices.md) for more security guidance.
