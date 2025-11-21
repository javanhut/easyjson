# Examples and Recipes

Real-world examples and common use cases for EasyJSON.

## REST API Handler

```go
func UserAPIHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request body safely
    body, _ := io.ReadAll(r.Body)
    result := easyjson.ParseSafely(string(body))
    
    if result.Error != nil {
        respondError(w, 400, "Invalid JSON")
        return
    }
    
    req := result.Data
    
    // Extract and validate
    name := req.GetString("name", "")
    email := req.GetString("email", "")
    
    if name == "" || email == "" {
        respondError(w, 400, "Name and email required")
        return
    }
    
    // Create user
    user := createUser(name, email)
    
    // Build response
    response := easyjson.QuickAPIResponse("success", 
        easyjson.QuickObject(
            "id", user.ID,
            "name", user.Name,
            "email", user.Email,
        ).Raw(),
        "User created successfully",
    )
    
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(response.ToPrettyString()))
}
```

## Configuration Manager

```go
type AppConfig struct {
    data *easyjson.JSONValue
}

func LoadConfig(path string) (*AppConfig, error) {
    data, err := easyjson.LoadFile(path)
    if err != nil {
        return nil, err
    }
    return &AppConfig{data: data}, nil
}

func (c *AppConfig) Database() DatabaseConfig {
    db := c.data.Get("database")
    return DatabaseConfig{
        Host:     db.GetString("host", "localhost"),
        Port:     db.GetInt("port", 5432),
        Database: db.GetString("database", "app"),
        User:     db.GetString("user", "postgres"),
        Password: db.GetString("password", ""),
    }
}

func (c *AppConfig) Server() ServerConfig {
    srv := c.data.Get("server")
    return ServerConfig{
        Host: srv.GetString("host", "0.0.0.0"),
        Port: srv.GetInt("port", 8080),
        TLS:  srv.GetBool("tls", false),
    }
}
```

## Data Transformation Pipeline

```go
func TransformAPIData(input string) (string, error) {
    data, err := easyjson.Loads(input)
    if err != nil {
        return "", err
    }
    
    users := data.Get("users")
    
    // Filter active users
    active := users.FilterArray(func(u *easyjson.JSONValue) bool {
        return u.GetBool("active", false)
    })
    
    // Sort by name
    sorted := active.SortBy("name")
    
    // Extract relevant fields
    simplified := sorted.MapArray(func(u *easyjson.JSONValue) interface{} {
        return map[string]interface{}{
            "id":    u.GetInt("id"),
            "name":  u.GetString("name"),
            "email": u.GetString("email"),
        }
    })
    
    // Build output
    output := easyjson.NewBuilder().
        AddField("users", simplified.Raw()).
        AddField("count", simplified.Len()).
        AddTimestamp("generated_at").
        ToJSON()
    
    return output.DumpsIndent("  ")
}
```

## Pagination Handler

```go
func ListUsers(page, limit int) *easyjson.JSONValue {
    // Get data
    allUsers := getUsersFromDB()
    total := len(allUsers)
    
    // Paginate
    start := (page - 1) * limit
    end := start + limit
    if end > total {
        end = total
    }
    pageUsers := allUsers[start:end]
    
    // Build response
    return easyjson.QuickPaginatedResponse(
        pageUsers,
        page,
        total,
        limit,
    )
}
```

## Error Response Builder

```go
func BuildErrorResponse(code, message string, details interface{}) string {
    response := easyjson.QuickErrorResponse(code, message, details)
    json, _ := response.DumpsIndent("  ")
    return json
}

// Usage
if err := validateInput(data); err != nil {
    response := BuildErrorResponse(
        "VALIDATION_ERROR",
        "Invalid input",
        err.Details,
    )
    w.WriteHeader(400)
    w.Write([]byte(response))
}
```

## Feature Flags

```go
type FeatureManager struct {
    config *easyjson.JSONValue
}

func (f *FeatureManager) IsEnabled(feature string) bool {
    return f.config.GetBool("features", feature, "enabled", false)
}

func (f *FeatureManager) GetValue(feature, key string) interface{} {
    return f.config.Get("features", feature, key).Raw()
}

// Usage
features := &FeatureManager{config: configData}
if features.IsEnabled("new_ui") {
    enableNewUI()
}
```

## Multi-Format API Client

```go
func ParseUserResponse(response string) (User, error) {
    data, err := easyjson.Loads(response)
    if err != nil {
        return User{}, err
    }
    
    // Handle various API formats
    return User{
        ID:    data.TryPaths("id", "user_id", "userId").AsInt(),
        Name:  data.TryPaths("name", "username", "display_name").AsString(),
        Email: data.TryPaths("email", "email_address", "emailAddress").AsString(),
    }, nil
}
```

## Validation Middleware

```go
func ValidateJSON(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" && r.Method != "DELETE" {
            body, _ := io.ReadAll(r.Body)
            r.Body = io.NopCloser(bytes.NewBuffer(body))
            
            result := easyjson.ParseSafely(string(body))
            if result.Error != nil {
                w.WriteHeader(400)
                w.Write([]byte(`{"error": "Invalid JSON"}`))
                return
            }
        }
        next.ServeHTTP(w, r)
    })
}
```

See [Best Practices](best-practices.md) for more patterns.
