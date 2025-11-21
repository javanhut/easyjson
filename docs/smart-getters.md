# Smart Getters and Safe Access

Smart getters provide the safest and most convenient way to access JSON data in EasyJSON, eliminating the need for nil checking and providing built-in default values.

## Table of Contents

- [Why Smart Getters?](#why-smart-getters)
- [Basic Smart Getters](#basic-smart-getters)
- [GetOr Method](#getor-method)
- [Safety Helpers](#safety-helpers)
- [Comparison with Other Methods](#comparison-with-other-methods)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Why Smart Getters?

### The Problem

Traditional access requires multiple steps and nil checking:

```go
// Traditional approach - verbose and error-prone
var name string
if user := data.Get("user"); !user.IsNull() {
    if nameVal := user.Get("name"); !nameVal.IsNull() {
        name = nameVal.AsString()
    } else {
        name = "Anonymous"
    }
} else {
    name = "Anonymous"
}
```

### The Solution

Smart getters do it all in one line:

```go
// Smart getter - clean and safe
name := data.GetString("user", "name", "Anonymous")
```

## Basic Smart Getters

### GetString

Get string values with optional default:

```go
func (jv *JSONValue) GetString(keys ...interface{}) string
```

**Usage:**
```go
// Single key with default
name := data.GetString("name", "Unknown")

// Nested path with default
email := data.GetString("user", "profile", "email", "no-email@example.com")

// Deep nesting
color := data.GetString("user", "settings", "theme", "colors", "primary", "#000000")

// Without default (returns "" if not found)
title := data.GetString("title")
```

**Examples:**
```go
// API response with optional fields
response := easyjson.ParseSafely(apiJSON).Data

// Safe extraction with defaults
username := response.GetString("user", "username", "guest")
displayName := response.GetString("user", "display_name", username)
bio := response.GetString("user", "bio", "No bio available")
avatar := response.GetString("user", "avatar_url", "/default-avatar.png")
```

### GetInt

Get integer values with optional default:

```go
func (jv *JSONValue) GetInt(keys ...interface{}) int
```

**Usage:**
```go
// Single key with default
age := data.GetInt("age", 0)

// Nested path with default
port := data.GetInt("server", "port", 8080)
timeout := data.GetInt("config", "timeout", 30)

// Without default (returns 0 if not found)
count := data.GetInt("count")
```

**Examples:**
```go
// Configuration with sensible defaults
config := easyjson.LoadFile("config.json")

port := config.GetInt("server", "port", 8080)
maxConnections := config.GetInt("server", "max_connections", 100)
timeout := config.GetInt("server", "timeout", 30)
retryCount := config.GetInt("retry", "count", 3)
retryDelay := config.GetInt("retry", "delay_ms", 1000)

// Pagination with defaults
page := request.GetInt("page", 1)
limit := request.GetInt("limit", 10)
offset := (page - 1) * limit
```

### GetBool

Get boolean values with optional default:

```go
func (jv *JSONValue) GetBool(keys ...interface{}) bool
```

**Usage:**
```go
// Single key with default
active := data.GetBool("active", false)

// Nested path with default
enabled := data.GetBool("features", "notifications", "enabled", true)
debug := data.GetBool("config", "debug", false)

// Without default (returns false if not found)
verified := data.GetBool("verified")
```

**Examples:**
```go
// Feature flags with defaults
features := data.Get("features")

darkMode := features.GetBool("dark_mode", false)
notifications := features.GetBool("notifications", true)
analytics := features.GetBool("analytics", true)
betaFeatures := features.GetBool("beta_access", false)

// User settings
user := data.Get("user")
emailNotifications := user.GetBool("settings", "email_notifications", true)
smsNotifications := user.GetBool("settings", "sms_notifications", false)
publicProfile := user.GetBool("settings", "public_profile", false)
```

### GetFloat

Get float64 values with optional default:

```go
func (jv *JSONValue) GetFloat(keys ...interface{}) float64
```

**Usage:**
```go
// Single key with default
price := data.GetFloat("price", 0.0)

// Nested path with default
rating := data.GetFloat("product", "rating", 0.0)
discount := data.GetFloat("pricing", "discount_percent", 0.0)

// Without default (returns 0.0 if not found)
weight := data.GetFloat("weight")
```

**Examples:**
```go
// E-commerce product
product := data.Get("product")

price := product.GetFloat("price", 0.0)
discount := product.GetFloat("discount", 0.0)
tax := product.GetFloat("tax_rate", 0.0)
weight := product.GetFloat("shipping", "weight_kg", 0.0)

finalPrice := price * (1 - discount/100) * (1 + tax/100)

// Financial calculations
transaction := data.Get("transaction")

amount := transaction.GetFloat("amount", 0.0)
fee := transaction.GetFloat("fee", 0.0)
exchangeRate := transaction.GetFloat("exchange_rate", 1.0)

total := (amount + fee) * exchangeRate
```

## GetOr Method

Generic getter that matches the type of the default value:

```go
func (jv *JSONValue) GetOr(keys ...interface{}) interface{}
```

**Smart Type Matching:**
```go
// Automatically matches type to default
str := data.GetOr("name", "Default").(string)
num := data.GetOr("count", 100).(int)
flag := data.GetOr("enabled", true).(bool)
rate := data.GetOr("rate", 0.5).(float64)
```

**Use Cases:**
```go
// When you need flexibility
func getConfigValue(data *easyjson.JSONValue, key string, defaultVal interface{}) interface{} {
    return data.GetOr(key, defaultVal)
}

// String config
timeout := getConfigValue(config, "timeout", "30s").(string)

// Numeric config
maxRetries := getConfigValue(config, "max_retries", 3).(int)

// Boolean config
debugMode := getConfigValue(config, "debug", false).(bool)
```

## Safety Helpers

### IsEmptyOrNull

Check if value is effectively empty:

```go
func (jv *JSONValue) IsEmptyOrNull() bool
```

Returns true for:
- Null values
- Empty strings
- Empty arrays
- Empty objects

**Usage:**
```go
if data.Get("name").IsEmptyOrNull() {
    fmt.Println("Name is required")
}

// Check multiple fields
requiredFields := []string{"name", "email", "phone"}
for _, field := range requiredFields {
    if data.Get(field).IsEmptyOrNull() {
        fmt.Printf("%s is required\n", field)
    }
}
```

### Type-Specific Or Methods

Convenience methods for simple cases:

```go
// StringOrEmpty - returns "" if null
name := data.Get("name").StringOrEmpty()

// IntOrZero - returns 0 if null
count := data.Get("count").IntOrZero()

// BoolOrFalse - returns false if null
active := data.Get("active").BoolOrFalse()

// FloatOrZero - returns 0.0 if null
price := data.Get("price").FloatOrZero()
```

**When to Use:**
```go
// Good: Already have JSONValue
userValue := data.Get("user")
name := userValue.Get("name").StringOrEmpty()
age := userValue.Get("age").IntOrZero()

// Better: Use smart getter directly
name := data.GetString("user", "name", "")
age := data.GetInt("user", "age", 0)
```

## Comparison with Other Methods

### Smart Getters vs Traditional Get

```go
// Traditional Get - requires nil checking
name := data.Get("user").Get("name").AsString()
if name == "" {
    name = "Anonymous"
}

// Smart Getter - one line, safe
name := data.GetString("user", "name", "Anonymous")
```

### Smart Getters vs Q

```go
// Q - requires separate default handling
name := data.Q("user", "name").AsString()
if name == "" {
    name = "Anonymous"
}

// Smart Getter - default built-in
name := data.GetString("user", "name", "Anonymous")
```

### Smart Getters vs Path

```go
// Path - requires separate default handling
age := data.Path("user.age").AsInt()
if age == 0 {
    age = 18
}

// Smart Getter - cleaner
age := data.GetInt("user", "age", 18)
```

### When to Use Each

**Use Smart Getters when:**
- You need default values
- You want maximum safety
- You're working with optional fields
- You're processing user input
- You want minimal code

**Use Q when:**
- You don't need defaults
- You'll check IsNull() separately
- You're chaining many operations
- Default would be misleading (like 0 for count)

**Use Get when:**
- You need to inspect intermediate values
- Debugging access chains
- Performance is critical (slight overhead)

**Use Path when:**
- Paths are stored as strings
- Working with configuration keys
- Paths are user-provided

## Best Practices

### 1. Always Provide Meaningful Defaults

```go
// Good: Meaningful default
port := config.GetInt("port", 8080)
timeout := config.GetInt("timeout", 30)

// Avoid: Using 0 when it might be valid
count := data.GetInt("count", 0) // Is 0 missing or actual count?

// Better: Check separately if 0 is valid
countVal := data.Get("count")
if countVal.IsNull() {
    // Handle missing count
} else {
    count := countVal.AsInt() // Could be 0
}
```

### 2. Use Descriptive Variable Names

```go
// Good: Clear what default means
maxRetries := config.GetInt("max_retries", 3)
defaultPageSize := config.GetInt("page_size", 10)

// Avoid: Ambiguous
retries := config.GetInt("max_retries", 3) // Max or current?
```

### 3. Group Related Defaults

```go
// Good: Related defaults together
const (
    DefaultPort = 8080
    DefaultTimeout = 30
    DefaultMaxConns = 100
)

port := config.GetInt("port", DefaultPort)
timeout := config.GetInt("timeout", DefaultTimeout)
maxConns := config.GetInt("max_connections", DefaultMaxConns)
```

### 4. Document Why Defaults Are Chosen

```go
// Good: Explain default choice
// Use port 8080 to avoid requiring root privileges
port := config.GetInt("server", "port", 8080)

// Default to 30s timeout based on average API response time
timeout := config.GetInt("timeout", 30)
```

### 5. Consider Configuration Hierarchy

```go
// Application defaults < Environment defaults < User config
func getPort(config *easyjson.JSONValue) int {
    // Try user config first
    if port := config.GetInt("port"); port != 0 {
        return port
    }
    
    // Try environment variable
    if envPort := os.Getenv("PORT"); envPort != "" {
        if port, err := strconv.Atoi(envPort); err == nil {
            return port
        }
    }
    
    // Fall back to application default
    return 8080
}
```

## Examples

### Example 1: User Profile Processing

```go
func processUserProfile(data *easyjson.JSONValue) UserProfile {
    // Extract with sensible defaults
    return UserProfile{
        ID:           data.GetInt("id", 0),
        Username:     data.GetString("username", "guest"),
        DisplayName:  data.GetString("display_name", "Guest User"),
        Email:        data.GetString("email", ""),
        Bio:          data.GetString("bio", "No bio provided"),
        AvatarURL:    data.GetString("avatar_url", "/default-avatar.png"),
        IsVerified:   data.GetBool("verified", false),
        IsActive:     data.GetBool("active", true),
        FollowerCount: data.GetInt("follower_count", 0),
        FollowingCount: data.GetInt("following_count", 0),
        PostCount:    data.GetInt("post_count", 0),
        Rating:       data.GetFloat("rating", 0.0),
        JoinedDate:   data.GetString("joined_date", ""),
        LastActive:   data.GetString("last_active", ""),
    }
}
```

### Example 2: Configuration Loading

```go
type ServerConfig struct {
    Host           string
    Port           int
    ReadTimeout    int
    WriteTimeout   int
    MaxConnections int
    EnableTLS      bool
    TLSCert        string
    TLSKey         string
    LogLevel       string
    EnableCORS     bool
}

func LoadServerConfig(configFile string) (*ServerConfig, error) {
    data, err := easyjson.LoadFile(configFile)
    if err != nil {
        // Return defaults on error
        return DefaultServerConfig(), nil
    }
    
    server := data.Get("server")
    
    return &ServerConfig{
        Host:           server.GetString("host", "0.0.0.0"),
        Port:           server.GetInt("port", 8080),
        ReadTimeout:    server.GetInt("read_timeout", 30),
        WriteTimeout:   server.GetInt("write_timeout", 30),
        MaxConnections: server.GetInt("max_connections", 100),
        EnableTLS:      server.GetBool("enable_tls", false),
        TLSCert:        server.GetString("tls_cert", ""),
        TLSKey:         server.GetString("tls_key", ""),
        LogLevel:       server.GetString("log_level", "info"),
        EnableCORS:     server.GetBool("enable_cors", true),
    }, nil
}
```

### Example 3: API Request Validation

```go
func validateAndExtractRequest(data *easyjson.JSONValue) (*Request, error) {
    // Extract with defaults
    req := &Request{
        Page:     data.GetInt("page", 1),
        Limit:    data.GetInt("limit", 10),
        SortBy:   data.GetString("sort_by", "created_at"),
        Order:    data.GetString("order", "desc"),
        Search:   data.GetString("search", ""),
        Filters:  data.Get("filters"),
    }
    
    // Validate extracted values
    if req.Page < 1 {
        return nil, fmt.Errorf("page must be >= 1")
    }
    
    if req.Limit < 1 || req.Limit > 100 {
        return nil, fmt.Errorf("limit must be between 1 and 100")
    }
    
    validOrders := []string{"asc", "desc"}
    if !contains(validOrders, req.Order) {
        return nil, fmt.Errorf("order must be 'asc' or 'desc'")
    }
    
    return req, nil
}
```

### Example 4: Multi-Source Configuration

```go
type AppConfig struct {
    config *easyjson.JSONValue
}

func (c *AppConfig) GetString(key string, defaultVal string) string {
    // Try user config
    if val := c.config.GetString(key); val != "" {
        return val
    }
    
    // Try environment variable
    envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
    if envVal := os.Getenv(envKey); envVal != "" {
        return envVal
    }
    
    // Use default
    return defaultVal
}

func (c *AppConfig) GetInt(key string, defaultVal int) int {
    // Try user config
    if val := c.config.GetInt(key); val != 0 {
        return val
    }
    
    // Try environment variable
    envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
    if envVal := os.Getenv(envKey); envVal != "" {
        if num, err := strconv.Atoi(envVal); err == nil {
            return num
        }
    }
    
    // Use default
    return defaultVal
}

// Usage
config := &AppConfig{config: configData}
port := config.GetString("server.port", "8080")
timeout := config.GetInt("server.timeout", 30)
```

### Example 5: Feature Flags

```go
type FeatureFlags struct {
    data *easyjson.JSONValue
}

func NewFeatureFlags(data *easyjson.JSONValue) *FeatureFlags {
    return &FeatureFlags{data: data}
}

func (f *FeatureFlags) IsEnabled(feature string) bool {
    return f.data.GetBool("features", feature, false)
}

func (f *FeatureFlags) GetConfig(feature, key string, defaultVal interface{}) interface{} {
    return f.data.GetOr("features", feature, "config", key, defaultVal)
}

// Usage
flags := NewFeatureFlags(configData)

// Simple flag check
if flags.IsEnabled("dark_mode") {
    enableDarkMode()
}

// Feature with configuration
if flags.IsEnabled("rate_limiting") {
    maxRequests := flags.GetConfig("rate_limiting", "max_requests", 100).(int)
    window := flags.GetConfig("rate_limiting", "window_seconds", 60).(int)
    setupRateLimiting(maxRequests, window)
}
```

## Summary

Smart getters are the recommended way to access JSON data in EasyJSON when:
- You need default values
- You want maximum safety
- You're processing optional fields
- You want clean, readable code

They eliminate boilerplate, prevent nil panics, and make your code more maintainable.

**Key Takeaways:**
- Use `GetString`, `GetInt`, `GetBool`, `GetFloat` for type-safe access
- Always provide meaningful defaults
- Document why defaults are chosen
- Consider configuration hierarchies
- Use `GetOr` for generic access with type inference

## Next Steps

- Learn about [Multi-Path Access](multi-path-access.md) for handling varying API formats
- Explore [Pattern Extractors](pattern-extractors.md) for common data patterns
- Review [Best Practices](best-practices.md) for production use
