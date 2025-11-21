# Multi-Path Access

Multi-path access helps you handle varying API formats and missing data gracefully.

## TryPaths - Try Multiple Paths

Try multiple paths until one works:

```go
// Handle different API response formats
title := data.TryPaths("title", "name", "label", "header").AsString()

// Different user ID formats
userID := data.TryPaths("user_id", "userId", "id", "ID").AsString()

// Deep paths with variations
email := data.TryPaths(
    "user.profile.email",
    "user.email_address",
    "email",
    "contact.email",
).AsString()
```

## DeepSearch - Find Keys Anywhere

Search for a key at any depth in the structure:

```go
// Find first occurrence of "email" anywhere
email := data.DeepSearch("email").AsString()

// Find all occurrences
allEmails := data.DeepSearchAll("email")
for _, email := range allEmails {
    fmt.Println(email.AsString())
}
```

## FindPath - Get Path to Key

Find the path to a key:

```go
// Returns: "user.profile.email"
emailPath := data.FindPath("email")

// Use the path
value := data.Path(emailPath)
```

## Key Checking

```go
// Check if any key exists
if data.HasAnyKey("name", "title", "label") {
    // At least one exists
}

// Check if all keys exist
if data.HasAllKeys("name", "email", "id") {
    // All exist
}
```

## Use Cases

### Handling Different API Versions

```go
func extractUserData(response *easyjson.JSONValue) User {
    return User{
        ID:    response.TryPaths("id", "user_id", "userId").AsInt(),
        Name:  response.TryPaths("name", "username", "display_name").AsString(),
        Email: response.TryPaths("email", "email_address", "emailAddress").AsString(),
    }
}
```

### Robust Data Extraction

```go
// Works with any of these structures:
// {"user": {"email": "..."}}
// {"profile": {"contact": {"email": "..."}}}
// {"email": "..."}
email := data.DeepSearch("email").AsString()
```

See [API Reference](api-reference.md) for complete method documentation.
