# Array Operations

EasyJSON provides powerful array operations inspired by JavaScript and Python, making it easy to filter, map, and transform JSON arrays.

## Table of Contents

- [Basic Array Operations](#basic-array-operations)
- [Filtering and Finding](#filtering-and-finding)
- [Transforming Arrays](#transforming-arrays)
- [Extracting Data](#extracting-data)
- [Grouping and Sorting](#grouping-and-sorting)
- [Array Utilities](#array-utilities)
- [Iteration](#iteration)
- [Predicates](#predicates)
- [Examples](#examples)

## Basic Array Operations

### Creating Arrays

```go
// Empty array
arr := easyjson.NewArray()

// From slice
items := easyjson.NewArrayFrom([]interface{}{"a", "b", "c"})

// Quick builder
arr := easyjson.QuickArray("item1", "item2", "item3")
```

### Adding Items

```go
arr := easyjson.NewArray()

// Append single item
arr.Append("item1")
arr.Append("item2")

// Extend with multiple items
arr.Extend([]interface{}{"item3", "item4", "item5"})

// Set by index
arr.Set(0, "new first item")
```

### Accessing Items

```go
// By index
first := arr.Get(0)
second := arr.Get(1)

// First and last
first := arr.First()
last := arr.Last()

// Length
count := arr.Len()
```

## Filtering and Finding

### Filter

Returns a new array with items that match a predicate:

```go
users := data.Get("users")

// Filter active users
activeUsers := users.FilterArray(func(user *easyjson.JSONValue) bool {
    return user.GetBool("active", false)
})

// Filter by age
adults := users.FilterArray(func(user *easyjson.JSONValue) bool {
    return user.GetInt("age") >= 18
})

// Complex filter
premiumActiveUsers := users.FilterArray(func(user *easyjson.JSONValue) bool {
    return user.GetBool("active") &&
           user.GetString("plan") == "premium" &&
           user.GetInt("age") >= 18
})
```

### Find

Find the first item that matches a condition:

```go
// Find with custom predicate
admin := users.FindInArray(func(user *easyjson.JSONValue) bool {
    return user.GetString("role") == "admin"
})

if !admin.IsNull() {
    fmt.Printf("Admin: %s\n", admin.GetString("name"))
}
```

### Find By Field

Convenient methods for finding by field value:

```go
// Find first user with specific ID
user := users.FindByField("id", 123)

// Find all admins
admins := users.FindAllByField("role", "admin")

// Works with different types
activeUser := users.FindByField("active", true)
moderator := users.FindByField("role", "moderator")
```

## Transforming Arrays

### Map

Transform each item in an array:

```go
users := data.Get("users")

// Extract names
names := users.MapArray(func(user *easyjson.JSONValue) interface{} {
    return user.GetString("name")
})

// Transform to display format
displayUsers := users.MapArray(func(user *easyjson.JSONValue) interface{} {
    return map[string]interface{}{
        "id":   user.GetInt("id"),
        "name": user.GetString("name"),
        "email": user.GetString("email"),
    }
})

// Calculate ages
ages := users.MapArray(func(user *easyjson.JSONValue) interface{} {
    birthYear := user.GetInt("birth_year")
    return 2024 - birthYear
})
```

### Reduce

Reduce array to a single value:

```go
numbers := data.Get("numbers")

// Sum all numbers
total := numbers.ReduceArray(0, func(sum interface{}, item *easyjson.JSONValue) interface{} {
    return sum.(int) + item.AsInt()
})

// Find maximum
max := numbers.ReduceArray(0, func(max interface{}, item *easyjson.JSONValue) interface{} {
    value := item.AsInt()
    if value > max.(int) {
        return value
    }
    return max
})

// Build a map
users := data.Get("users")
userMap := users.ReduceArray(
    make(map[int]string),
    func(m interface{}, user *easyjson.JSONValue) interface{} {
        userMap := m.(map[int]string)
        id := user.GetInt("id")
        name := user.GetString("name")
        userMap[id] = name
        return userMap
    },
})
```

## Extracting Data

### Pluck

Extract a specific field from all items:

```go
users := data.Get("users")

// Pluck returns JSONValue array
names := users.Pluck("name")

// PluckStrings returns []string
nameStrings := users.PluckStrings("name")
// Result: ["Alice", "Bob", "Charlie"]

// PluckInts returns []int
ages := users.PluckInts("age")
// Result: [25, 30, 28]
```

### Extract Array Field

Alias for PluckStrings:

```go
emails := users.ExtractArrayField("email")
// Same as: users.PluckStrings("email")
```

## Grouping and Sorting

### Group By

Group array items by field value:

```go
users := data.Get("users")

// Group by role
roleGroups := users.GroupBy("role")
// Returns: map[string][]*JSONValue

// Access groups
admins := roleGroups["admin"]
regularUsers := roleGroups["user"]

fmt.Printf("Admins: %d\n", len(admins))
fmt.Printf("Users: %d\n", len(regularUsers))

// Process each group
for role, usersInRole := range roleGroups {
    fmt.Printf("%s: %d users\n", role, len(usersInRole))
}
```

### Count By Field

Count occurrences of each value:

```go
users := data.Get("users")

// Count by role
roleCounts := users.CountByField("role")
// Returns: map[string]int
// Example: {"admin": 2, "user": 10, "moderator": 3}

for role, count := range roleCounts {
    fmt.Printf("%s: %d\n", role, count)
}
```

### Sort By

Sort array by field value:

```go
users := data.Get("users")

// Sort by name (case-insensitive)
sortedUsers := users.SortBy("name")

// Sort is stable
sortedByAge := users.SortBy("age")
```

## Array Utilities

### Unique

Remove duplicate items:

```go
tags := data.Get("tags")

// Remove duplicates
uniqueTags := tags.Unique()

// Example:
// ["go", "json", "go", "api", "json"] 
// becomes: ["go", "json", "api"]
```

### Take

Get first N items:

```go
users := data.Get("users")

// Get first 5 users
firstFive := users.Take(5)

// Safe if array is smaller
firstThree := users.Take(3) // Works even if array has only 2 items
```

### Skip

Skip first N items:

```go
users := data.Get("users")

// Skip first 10 users (pagination)
remaining := users.Skip(10)

// Combine with Take for pagination
page2 := users.Skip(10).Take(10) // Items 11-20
```

### Slice

Combine Skip and Take for pagination:

```go
func getPage(data *easyjson.JSONValue, page, pageSize int) *easyjson.JSONValue {
    offset := (page - 1) * pageSize
    return data.Skip(offset).Take(pageSize)
}

// Get page 3 with 20 items per page
page3 := getPage(users, 3, 20)
```

## Iteration

### For Each

Iterate with index and value:

```go
users := data.Get("users")

users.ForEach(func(i int, user *easyjson.JSONValue) {
    fmt.Printf("%d: %s (%s)\n", 
        i, 
        user.GetString("name"),
        user.GetString("email"))
})
```

### Traditional Loop

```go
for i := 0; i < users.Len(); i++ {
    user := users.Get(i)
    // Process user
}
```

### AsArray

Convert to Go slice for more control:

```go
userSlice := users.AsArray() // []*JSONValue

for i, user := range userSlice {
    fmt.Printf("%d: %s\n", i, user.GetString("name"))
}
```

## Predicates

### Some

Check if at least one item matches:

```go
users := data.Get("users")

// Check if any user is an admin
hasAdmin := users.Some(func(user *easyjson.JSONValue) bool {
    return user.GetString("role") == "admin"
})

if hasAdmin {
    fmt.Println("At least one admin found")
}

// Check if any user is underage
hasMinor := users.Some(func(user *easyjson.JSONValue) bool {
    return user.GetInt("age") < 18
})
```

### Every

Check if all items match:

```go
users := data.Get("users")

// Check if all users are active
allActive := users.Every(func(user *easyjson.JSONValue) bool {
    return user.GetBool("active", false)
})

// Check if all users are adults
allAdults := users.Every(func(user *easyjson.JSONValue) bool {
    return user.GetInt("age") >= 18
})

if allActive && allAdults {
    fmt.Println("All users are active adults")
}
```

## Examples

### Example 1: User Processing Pipeline

```go
func processUsers(data *easyjson.JSONValue) {
    users := data.Get("users")

    // Filter active premium users
    premium := users.FilterArray(func(user *easyjson.JSONValue) bool {
        return user.GetBool("active") &&
               user.GetString("plan") == "premium"
    })

    // Sort by name
    sorted := premium.SortBy("name")

    // Extract emails
    emails := sorted.PluckStrings("email")

    // Send notifications
    for _, email := range emails {
        sendEmail(email, "Premium User Update")
    }
}
```

### Example 2: Analytics Aggregation

```go
func analyzeOrders(data *easyjson.JSONValue) {
    orders := data.Get("orders")

    // Calculate total revenue
    totalRevenue := orders.ReduceArray(0.0, 
        func(sum interface{}, order *easyjson.JSONValue) interface{} {
            return sum.(float64) + order.GetFloat("total")
        },
    )

    // Group by status
    statusGroups := orders.GroupBy("status")

    // Count by payment method
    paymentCounts := orders.CountByField("payment_method")

    fmt.Printf("Total Revenue: $%.2f\n", totalRevenue)
    fmt.Printf("Pending: %d\n", len(statusGroups["pending"]))
    fmt.Printf("Completed: %d\n", len(statusGroups["completed"]))
    
    for method, count := range paymentCounts {
        fmt.Printf("%s payments: %d\n", method, count)
    }
}
```

### Example 3: Data Transformation

```go
func transformAPIResponse(data *easyjson.JSONValue) *easyjson.JSONValue {
    users := data.Get("data", "users")

    // Transform to simplified format
    simplified := users.MapArray(func(user *easyjson.JSONValue) interface{} {
        return map[string]interface{}{
            "id":       user.GetInt("id"),
            "name":     user.GetString("full_name"),
            "email":    user.GetString("email_address"),
            "role":     user.GetString("user_role"),
            "active":   user.GetBool("is_active"),
            "created":  user.GetString("created_at"),
        }
    })

    // Build response
    return easyjson.NewBuilder().
        AddField("users", simplified.Raw()).
        AddField("count", simplified.Len()).
        AddTimestamp("generated_at").
        ToJSON()
}
```

### Example 4: Search and Filter

```go
func searchUsers(data *easyjson.JSONValue, query string) []*easyjson.JSONValue {
    users := data.Get("users")
    query = strings.ToLower(query)

    // Find users matching query
    results := users.FilterArray(func(user *easyjson.JSONValue) bool {
        name := strings.ToLower(user.GetString("name"))
        email := strings.ToLower(user.GetString("email"))
        role := strings.ToLower(user.GetString("role"))
        
        return strings.Contains(name, query) ||
               strings.Contains(email, query) ||
               strings.Contains(role, query)
    })

    return results.AsArray()
}
```

### Example 5: Validation

```go
func validateUsers(data *easyjson.JSONValue) []string {
    users := data.Get("users")
    var errors []string

    users.ForEach(func(i int, user *easyjson.JSONValue) {
        // Check required fields
        if user.GetString("email") == "" {
            errors = append(errors, 
                fmt.Sprintf("User %d: missing email", i))
        }

        // Validate email format
        email := user.Get("email")
        if !email.IsValidEmail() {
            errors = append(errors, 
                fmt.Sprintf("User %d: invalid email format", i))
        }

        // Check age range
        age := user.GetInt("age")
        if age < 18 || age > 120 {
            errors = append(errors, 
                fmt.Sprintf("User %d: invalid age %d", i, age))
        }
    })

    return errors
}
```

## Performance Tips

### Efficient Chains

```go
// Efficient: Chain operations
result := users.
    FilterArray(isActive).
    SortBy("name").
    Take(10)

// Less efficient: Multiple intermediate variables
filtered := users.FilterArray(isActive)
sorted := filtered.SortBy("name")
result := sorted.Take(10)
```

### Reuse Predicates

```go
// Define predicate once
isActive := func(user *easyjson.JSONValue) bool {
    return user.GetBool("active", false)
}

// Reuse multiple times
activeUsers := users.FilterArray(isActive)
hasActive := users.Some(isActive)
allActive := users.Every(isActive)
```

### Avoid Unnecessary Conversions

```go
// Good: Work with JSONValue
result := users.FilterArray(predicate).SortBy("name")

// Less efficient: Convert unnecessarily
userSlice := users.AsArray()
// ... complex Go operations ...
result := easyjson.NewArrayFrom(userSlice)
```

## Next Steps

- Learn about [Smart Getters](smart-getters.md) for safe field access
- Explore [Pattern Extractors](pattern-extractors.md) for common operations
- See more [Examples](examples.md) for real-world use cases
