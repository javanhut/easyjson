# Pattern Extractors

Pattern extractors provide convenient methods for extracting common data patterns from JSON structures.

## User Information

Extract user data from any structure:

```go
userInfo := data.GetUserInfo()
// Returns: map[string]string{
//   "id": "123",
//   "name": "John Doe",
//   "email": "john@example.com",
//   "role": "admin",
//   "username": "johndoe",
//   "phone": "555-0123"
// }
```

## Pagination Information

Extract pagination from any API format:

```go
pagination := data.GetPaginationInfo()
// Returns: map[string]int{
//   "page": 1,
//   "total": 100,
//   "limit": 10,
//   "offset": 0,
//   "total_pages": 10
// }
```

## Timestamps

Extract timestamp fields:

```go
timestamps := data.GetTimestamps()
// Returns: map[string]string{
//   "created": "2023-01-01T00:00:00Z",
//   "updated": "2023-01-02T00:00:00Z",
//   "deleted": "",
//   "published": ""
// }
```

## API Response Info

Extract standard API response fields:

```go
response := data.GetAPIResponseInfo()
// Returns: map[string]string{
//   "status": "success",
//   "message": "OK",
//   "error": "",
//   "code": "200"
// }
```

## Status Checking

```go
if data.IsAPISuccess() {
    // Handle success
}

if data.IsAPIError() {
    error := data.GetAPIResponseInfo()["error"]
    // Handle error
}
```

## Specialized Extractors

### Contact Information

```go
contact := data.GetContactInfo()
// email, phone, address, city, state, country, zip
```

### Product Information

```go
product := data.GetProductInfo()
// id, name, description, price, currency, category, brand
```

### Location Information

```go
location := data.GetLocationInfo()
// latitude, longitude, address, city, state, country, zip, timezone
```

### Financial Information

```go
financial := data.GetFinancialInfo()
// amount, currency, tax, discount, subtotal, transaction_id, status
```

## Validation Helpers

### HasRequiredFields

```go
if data.HasRequiredFields("name", "email", "phone") {
    // All fields present
}
```

### GetMissingFields

```go
missing := data.GetMissingFields("name", "email", "phone")
if len(missing) > 0 {
    fmt.Printf("Missing: %v\n", missing)
}
```

### IsComplete

```go
if data.IsComplete("user") {
    // Has all required user fields (id, name, email)
}
```

### GetCompletionScore

```go
score := data.GetCompletionScore("user") // 0.0 to 1.0
fmt.Printf("Profile is %.0f%% complete\n", score*100)
```

## Data Validation

### IsValidEmail

```go
email := data.Get("email")
if email.IsValidEmail() {
    sendEmail(email.AsString())
}
```

### IsValidURL

```go
website := data.Get("website")
if website.IsValidURL() {
    openURL(website.AsString())
}
```

### IsValidDate

```go
date := data.Get("created_at")
if date.IsValidDate() {
    formatted := date.GetFormattedDate("2006-01-02")
    relative := date.GetRelativeTime() // "2 hours ago"
}
```

## Data Sanitization

Remove sensitive fields for safe output:

```go
publicData := data.SanitizeForOutput()
// Removes: password, secret, token, key, ssn, credit_card, etc.
```

## Summary Methods

### GetSummary

Get overview of JSON structure:

```go
summary := data.GetSummary()
// Returns: map[string]interface{}{
//   "type": "object",
//   "size": 1234,
//   "keys": 5,
//   "key_names": []string{...},
//   "value_types": map[string]int{...}
// }
```

### TypeString

```go
typeStr := data.TypeString()
// Returns: "null", "string", "number", "boolean", "array", or "object"
```

## Examples

### Extract User from Varying Formats

```go
// Works with any of these structures:
// {"user_id": 123, "user_name": "John", "user_email": "..."}
// {"id": 123, "name": "John", "email": "..."}
// {"userId": 123, "username": "John", "emailAddress": "..."}

userInfo := data.GetUserInfo()
fmt.Printf("User: %s (%s)\n", userInfo["name"], userInfo["email"])
```

### Handle API Pagination

```go
func processPaginatedData(response *easyjson.JSONValue) {
    if response.HasPagination() {
        pagination := response.GetPaginationInfo()
        fmt.Printf("Page %d of %d\n", 
            pagination["page"], 
            pagination["total_pages"])
    }
}
```

### Validate Complete User Profile

```go
func validateProfile(data *easyjson.JSONValue) error {
    if !data.HasRequiredFields("name", "email", "phone") {
        missing := data.GetMissingFields("name", "email", "phone")
        return fmt.Errorf("missing fields: %v", missing)
    }
    
    email := data.Get("email")
    if !email.IsValidEmail() {
        return fmt.Errorf("invalid email")
    }
    
    score := data.GetCompletionScore("user")
    if score < 0.8 {
        return fmt.Errorf("profile incomplete (%.0f%%)", score*100)
    }
    
    return nil
}
```

See [API Reference](api-reference.md) for complete method signatures.
