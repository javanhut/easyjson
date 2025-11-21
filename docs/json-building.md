# JSON Building

EasyJSON provides a fluent, chainable API for building JSON structures programmatically. This guide covers all building techniques from simple to complex.

## Table of Contents

- [Basic Building](#basic-building)
- [Fluent Builder API](#fluent-builder-api)
- [Conditional Building](#conditional-building)
- [Quick Builders](#quick-builders)
- [Building API Responses](#building-api-responses)
- [Advanced Patterns](#advanced-patterns)
- [Best Practices](#best-practices)
- [Examples](#examples)

## Basic Building

### Creating Empty Structures

```go
// Empty object
obj := easyjson.NewObject()

// Empty array
arr := easyjson.NewArray()

// From existing data
obj := easyjson.NewObjectFrom(map[string]interface{}{
    "name": "John",
    "age":  30,
})

arr := easyjson.NewArrayFrom([]interface{}{"a", "b", "c"})
```

### Manual Building

```go
// Build object step by step
user := easyjson.NewObject()
user.Set("id", 123)
user.Set("name", "John Doe")
user.Set("email", "john@example.com")
user.Set("active", true)

// Build array step by step
tags := easyjson.NewArray()
tags.Append("golang")
tags.Append("json")
tags.Append("api")

// Nest structures
user.Set("tags", tags.Raw())

// Set nested paths (creates intermediate objects)
user.SetPath("profile.bio", "Software developer")
user.SetPath("profile.website", "https://example.com")
```

## Fluent Builder API

### NewBuilder

Create a fluent builder for objects:

```go
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddField("message", "User created").
    AddField("timestamp", time.Now().Unix()).
    ToJSON()
```

### AddField

Add a single field to the object:

```go
builder := easyjson.NewBuilder().
    AddField("id", 123).
    AddField("name", "John Doe").
    AddField("email", "john@example.com").
    AddField("verified", true)
```

### AddFields

Add multiple fields at once:

```go
builder := easyjson.NewBuilder().
    AddFields(map[string]interface{}{
        "id":    123,
        "name":  "John",
        "email": "john@example.com",
        "age":   30,
    })
```

### AddObject

Add a nested object using a builder function:

```go
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddObject("user", func(user *easyjson.JSONBuilder) {
        user.AddField("id", 123).
            AddField("name", "John Doe").
            AddField("email", "john@example.com")
    }).
    ToJSON()
```

### AddArray

Add a nested array using a builder function:

```go
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddArray("items", func(items *easyjson.JSONBuilder) {
        items.AddItem("item1").
             AddItem("item2").
             AddItem("item3")
    }).
    ToJSON()
```

### NewArrayBuilder

Create a fluent builder for arrays:

```go
arr := easyjson.NewArrayBuilder().
    AddItem("first").
    AddItem("second").
    AddItem("third").
    ToJSON()
```

### AddItem

Add a single item to an array:

```go
builder := easyjson.NewArrayBuilder().
    AddItem(map[string]interface{}{
        "id":   1,
        "name": "Alice",
    }).
    AddItem(map[string]interface{}{
        "id":   2,
        "name": "Bob",
    })
```

### AddItems

Add multiple items at once:

```go
builder := easyjson.NewArrayBuilder().
    AddItems("item1", "item2", "item3")
```

## Conditional Building

### AddIf

Add field only if condition is true:

```go
builder := easyjson.NewBuilder().
    AddField("name", "John").
    AddIf(isAdmin, "admin_panel_url", "/admin").
    AddIf(isPremium, "premium_features", premiumFeatures)
```

### AddIfNotEmpty

Add field only if value is not empty:

```go
builder := easyjson.NewBuilder().
    AddField("name", "John").
    AddIfNotEmpty("bio", bio).           // Only if bio != ""
    AddIfNotEmpty("website", website).   // Only if website != ""
    AddIfNotEmpty("company", company)    // Only if company != ""
```

### When

Execute builder function if condition is true:

```go
builder := easyjson.NewBuilder().
    AddField("name", "John").
    When(includeDetails, func(b *easyjson.JSONBuilder) {
        b.AddField("email", email).
          AddField("phone", phone).
          AddField("address", address)
    }).
    When(includeStats, func(b *easyjson.JSONBuilder) {
        b.AddField("login_count", loginCount).
          AddField("last_login", lastLogin)
    })
```

### Unless

Execute builder function if condition is false (opposite of When):

```go
builder := easyjson.NewBuilder().
    AddField("name", "John").
    Unless(isGuest, func(b *easyjson.JSONBuilder) {
        b.AddField("email", email).
          AddField("member_since", memberSince)
    }).
    Unless(isPublicView, func(b *easyjson.JSONBuilder) {
        b.AddField("internal_id", internalID).
          AddField("permissions", permissions)
    })
```

## Quick Builders

### QuickObject

Create an object with inline key-value pairs:

```go
user := easyjson.QuickObject(
    "id", 123,
    "name", "John",
    "age", 30,
    "active", true,
)
```

### QuickArray

Create an array with inline items:

```go
tags := easyjson.QuickArray("golang", "json", "api", "backend")
```

### QuickAPIResponse

Create a standard API response:

```go
response := easyjson.QuickAPIResponse(
    "success",      // status
    userData,       // data
    "User created", // message
)

// Produces:
// {
//   "status": "success",
//   "data": {...userData...},
//   "message": "User created",
//   "timestamp": 1234567890
// }
```

### QuickErrorResponse

Create a standard error response:

```go
response := easyjson.QuickErrorResponse(
    "VALIDATION_ERROR",           // error code
    "Invalid input data",         // message
    validationErrors,             // details
)

// Produces:
// {
//   "status": "error",
//   "timestamp": 1234567890,
//   "error": {
//     "code": "VALIDATION_ERROR",
//     "message": "Invalid input data",
//     "details": {...validationErrors...}
//   }
// }
```

### QuickPaginatedResponse

Create a paginated response:

```go
response := easyjson.QuickPaginatedResponse(
    users,  // data
    2,      // current page
    150,    // total items
    10,     // items per page
)

// Produces:
// {
//   "status": "success",
//   "data": [...users...],
//   "pagination": {
//     "page": 2,
//     "total": 150,
//     "limit": 10,
//     "has_next": true,
//     "has_prev": true
//   },
//   "timestamp": 1234567890
// }
```

## Building API Responses

### AddTimestamp

Add current Unix timestamp:

```go
builder := easyjson.NewBuilder().
    AddField("status", "success").
    AddTimestamp("created_at")
```

### AddISO8601Timestamp

Add current timestamp in ISO8601 format:

```go
builder := easyjson.NewBuilder().
    AddField("status", "success").
    AddISO8601Timestamp("created_at")
// Adds: "created_at": "2024-01-01T12:00:00Z"
```

### AddAPIStatus

Add standard API status fields:

```go
builder := easyjson.NewBuilder().
    AddAPIStatus("success", "Operation completed successfully")
// Adds: status, message, timestamp
```

### AddPaginationInfo

Add pagination metadata:

```go
builder := easyjson.NewBuilder().
    AddField("users", users).
    AddPaginationInfo(2, 150, 10)
// Adds: pagination object with page, total, limit, has_next, has_prev
```

### AddError

Add error information:

```go
builder := easyjson.NewBuilder().
    AddField("status", "error").
    AddError("VALIDATION_ERROR", "Invalid input", validationDetails)
```

### AddUserInfo

Add user information from struct or map:

```go
builder := easyjson.NewBuilder().
    AddUserInfo(userStruct)
// Automatically extracts: id, name, email, username, role, etc.
```

## Advanced Patterns

### Building Complex Responses

```go
response := easyjson.NewBuilder().
    AddAPIStatus("success", "Data retrieved").
    AddObject("data", func(data *easyjson.JSONBuilder) {
        data.AddObject("user", func(user *easyjson.JSONBuilder) {
            user.AddField("id", userID).
                AddField("name", userName).
                AddField("email", userEmail).
                AddObject("profile", func(profile *easyjson.JSONBuilder) {
                    profile.AddField("bio", bio).
                           AddField("avatar", avatarURL).
                           AddField("location", location)
                }).
                AddArray("permissions", func(perms *easyjson.JSONBuilder) {
                    for _, perm := range permissions {
                        perms.AddItem(perm)
                    }
                })
        }).
        AddArray("recent_activity", func(activity *easyjson.JSONBuilder) {
            for _, item := range recentItems {
                activity.AddItem(map[string]interface{}{
                    "type":      item.Type,
                    "timestamp": item.Timestamp,
                    "details":   item.Details,
                })
            }
        })
    }).
    AddPaginationInfo(page, total, limit).
    AddObject("metadata", func(meta *easyjson.JSONBuilder) {
        meta.AddField("version", "2.0").
            AddField("server", serverID).
            AddISO8601Timestamp("generated_at")
    }).
    ToJSON()
```

### Building from Structs

```go
type User struct {
    ID       int
    Name     string
    Email    string
    Role     string
    Active   bool
    Settings map[string]interface{}
}

func BuildUserResponse(user User) *easyjson.JSONValue {
    return easyjson.NewBuilder().
        AddField("id", user.ID).
        AddField("name", user.Name).
        AddField("email", user.Email).
        AddField("role", user.Role).
        AddField("active", user.Active).
        AddField("settings", user.Settings).
        AddTimestamp("retrieved_at").
        ToJSON()
}
```

### Builder Utilities

```go
// Get size of built JSON
size := builder.Size() // bytes

// Validate built JSON
if !builder.Validate() {
    log.Fatal("Invalid JSON structure")
}

// Clone builder
builderCopy := builder.Clone()

// Reset builder
builder.Reset() // Clears all fields

// Merge with existing JSONValue
builder.Merge(existingData)
```

### Output Methods

```go
// Get as JSONValue
jsonValue := builder.ToJSON()

// Get as JSON string
jsonString := builder.ToJSONString()

// Get as pretty JSON string
prettyJSON := builder.ToPrettyString()

// Get as bytes
jsonBytes := builder.ToBytes()
```

## Best Practices

### 1. Use Builders for Complex Structures

```go
// Good: Use builder for nested structures
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddObject("user", func(user *easyjson.JSONBuilder) {
        user.AddField("id", 123).
            AddField("name", "John")
    }).
    ToJSON()

// Avoid: Manual nesting gets messy
user := easyjson.NewObject()
user.Set("id", 123)
user.Set("name", "John")
response := easyjson.NewObject()
response.Set("status", "success")
response.Set("user", user.Raw())
```

### 2. Use Quick Builders for Simple Cases

```go
// Good: Quick builder for simple objects
config := easyjson.QuickObject("port", 8080, "host", "localhost")

// Overkill: Full builder for simple case
config := easyjson.NewBuilder().
    AddField("port", 8080).
    AddField("host", "localhost").
    ToJSON()
```

### 3. Use Conditional Methods

```go
// Good: Use AddIf for optional fields
builder := easyjson.NewBuilder().
    AddField("name", "John").
    AddIf(user.IsAdmin, "admin_data", adminInfo).
    AddIfNotEmpty("bio", user.Bio)

// Avoid: Manual conditionals
builder := easyjson.NewBuilder().
    AddField("name", "John")
if user.IsAdmin {
    builder.AddField("admin_data", adminInfo)
}
if user.Bio != "" {
    builder.AddField("bio", user.Bio)
}
```

### 4. Extract to Helper Functions

```go
// Good: Extract repeated patterns
func buildUserObject(user User) *easyjson.JSONValue {
    return easyjson.NewBuilder().
        AddField("id", user.ID).
        AddField("name", user.Name).
        AddField("email", user.Email).
        ToJSON()
}

// Use in multiple places
response := easyjson.NewBuilder().
    AddField("user", buildUserObject(user).Raw()).
    AddField("friend", buildUserObject(friend).Raw()).
    ToJSON()
```

### 5. Use Standard Response Builders

```go
// Good: Use standard response builders
response := easyjson.QuickAPIResponse("success", data, "OK")

// Avoid: Building standard responses manually
response := easyjson.NewBuilder().
    AddField("status", "success").
    AddField("data", data).
    AddField("message", "OK").
    AddTimestamp("timestamp").
    ToJSON()
```

## Examples

### Example 1: REST API Response

```go
func BuildUserListResponse(users []User, page, total, limit int) string {
    response := easyjson.NewBuilder().
        AddAPIStatus("success", "Users retrieved").
        AddArray("users", func(arr *easyjson.JSONBuilder) {
            for _, user := range users {
                arr.AddItem(map[string]interface{}{
                    "id":       user.ID,
                    "name":     user.Name,
                    "email":    user.Email,
                    "role":     user.Role,
                    "active":   user.Active,
                    "joined":   user.JoinedDate,
                })
            }
        }).
        AddPaginationInfo(page, total, limit).
        AddObject("metadata", func(meta *easyjson.JSONBuilder) {
            meta.AddField("api_version", "v2").
                AddField("server", "api-01").
                AddISO8601Timestamp("generated_at")
        })
    
    return response.ToPrettyString()
}
```

### Example 2: Configuration File Generation

```go
func GenerateConfig(settings Settings) error {
    config := easyjson.NewBuilder().
        AddObject("server", func(server *easyjson.JSONBuilder) {
            server.AddField("host", settings.Host).
                  AddField("port", settings.Port).
                  AddField("read_timeout", settings.ReadTimeout).
                  AddField("write_timeout", settings.WriteTimeout)
        }).
        AddObject("database", func(db *easyjson.JSONBuilder) {
            db.AddField("driver", settings.DBDriver).
              AddField("host", settings.DBHost).
              AddField("port", settings.DBPort).
              AddField("database", settings.DBName).
              AddField("max_connections", settings.MaxConns)
        }).
        AddObject("logging", func(log *easyjson.JSONBuilder) {
            log.AddField("level", settings.LogLevel).
               AddField("file", settings.LogFile).
               AddField("max_size_mb", settings.LogMaxSize)
        }).
        AddObject("features", func(feat *easyjson.JSONBuilder) {
            feat.AddField("enable_cors", settings.EnableCORS).
                AddField("enable_tls", settings.EnableTLS).
                AddField("enable_metrics", settings.EnableMetrics)
        })
    
    return config.ToJSON().SaveFileIndent("config.json", "  ")
}
```

### Example 3: Dynamic Form Response

```go
func BuildFormResponse(form Form, options FormOptions) *easyjson.JSONValue {
    return easyjson.NewBuilder().
        AddField("form_id", form.ID).
        AddField("title", form.Title).
        AddField("description", form.Description).
        AddArray("fields", func(fields *easyjson.JSONBuilder) {
            for _, field := range form.Fields {
                fields.AddItem(map[string]interface{}{
                    "id":          field.ID,
                    "type":        field.Type,
                    "label":       field.Label,
                    "placeholder": field.Placeholder,
                    "required":    field.Required,
                    "validation":  field.ValidationRules,
                })
            }
        }).
        When(options.IncludeSubmissions, func(b *easyjson.JSONBuilder) {
            b.AddField("submission_count", form.SubmissionCount).
              AddField("last_submission", form.LastSubmission)
        }).
        Unless(options.IsPublic, func(b *easyjson.JSONBuilder) {
            b.AddField("owner_id", form.OwnerID).
              AddField("created_at", form.CreatedAt).
              AddField("updated_at", form.UpdatedAt)
        }).
        ToJSON()
}
```

### Example 4: Error Response with Details

```go
func BuildValidationError(errors map[string][]string) *easyjson.JSONValue {
    return easyjson.NewBuilder().
        AddField("status", "error").
        AddObject("error", func(err *easyjson.JSONBuilder) {
            err.AddField("code", "VALIDATION_ERROR").
               AddField("message", "Request validation failed").
               AddObject("details", func(details *easyjson.JSONBuilder) {
                   for field, fieldErrors := range errors {
                       details.AddArray(field, func(arr *easyjson.JSONBuilder) {
                           for _, errMsg := range fieldErrors {
                               arr.AddItem(errMsg)
                           }
                       })
                   }
               })
        }).
        AddISO8601Timestamp("timestamp").
        ToJSON()
}
```

### Example 5: GraphQL-Style Response

```go
func BuildGraphQLResponse(data interface{}, errors []Error) *easyjson.JSONValue {
    builder := easyjson.NewBuilder()
    
    if data != nil {
        builder.AddField("data", data)
    } else {
        builder.AddField("data", nil)
    }
    
    if len(errors) > 0 {
        builder.AddArray("errors", func(arr *easyjson.JSONBuilder) {
            for _, err := range errors {
                arr.AddItem(map[string]interface{}{
                    "message": err.Message,
                    "path":    err.Path,
                    "extensions": map[string]interface{}{
                        "code":      err.Code,
                        "timestamp": err.Timestamp,
                    },
                })
            }
        })
    }
    
    return builder.ToJSON()
}
```

## Summary

The builder API provides:
- Fluent, chainable syntax
- Conditional building
- Quick builders for common patterns
- Type safety and validation
- Clean, readable code

Choose the right approach:
- **Simple objects**: Use `QuickObject`
- **Standard responses**: Use `QuickAPIResponse`, `QuickErrorResponse`
- **Complex structures**: Use `NewBuilder` with nested functions
- **Optional fields**: Use `AddIf`, `When`, `Unless`
- **Repeated patterns**: Extract to helper functions

## Next Steps

- Learn about [Pattern Extractors](pattern-extractors.md) for extracting data
- Review [Best Practices](best-practices.md) for production use
- See more [Examples](examples.md) for real-world patterns
