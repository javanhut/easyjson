# API Reference

Complete reference for all EasyJSON methods and functions.

## Quick Reference Table

| Category | Methods |
|----------|---------|
| **Parsing** | `Loads`, `Load`, `LoadFile`, `ParseSafely`, `ParseSafelyFrom`, `TryParse`, `ParseOrDefault` |
| **Creation** | `New`, `NewObject`, `NewArray`, `NewObjectFrom`, `NewArrayFrom` |
| **Access** | `Get`, `Q`, `Path`, `GetString`, `GetInt`, `GetBool`, `GetFloat`, `GetOr` |
| **Type Check** | `IsNull`, `IsString`, `IsNumber`, `IsBool`, `IsArray`, `IsObject` |
| **Conversion** | `AsString`, `AsInt`, `AsFloat`, `AsBool`, `AsArray`, `AsObject`, `Raw` |
| **Modification** | `Set`, `SetPath`, `Delete`, `Append`, `Extend`, `Update` |
| **Array Ops** | `FilterArray`, `MapArray`, `ReduceArray`, `FindByField`, `PluckStrings`, `GroupBy` |
| **Multi-Path** | `TryPaths`, `TryKeys`, `DeepSearch`, `FindPath`, `HasAnyKey` |
| **Serialization** | `Dumps`, `DumpsIndent`, `SaveFile`, `SaveFileIndent` |
| **Utility** | `Clone`, `Keys`, `Values`, `Items`, `Len`, `Has` |

## Parsing Functions

### Loads

```go
func Loads(jsonStr string) (*JSONValue, error)
```

Parse JSON from string.

**Parameters:**
- `jsonStr`: JSON string to parse

**Returns:**
- `*JSONValue`: Parsed JSON data
- `error`: Parse error if invalid

**Example:**
```go
data, err := easyjson.Loads(`{"name": "John", "age": 30}`)
if err != nil {
    log.Fatal(err)
}
```

### Load

```go
func Load(jsonBytes []byte) (*JSONValue, error)
```

Parse JSON from byte slice.

**Parameters:**
- `jsonBytes`: JSON bytes to parse

**Returns:**
- `*JSONValue`: Parsed JSON data
- `error`: Parse error if invalid

**Example:**
```go
jsonBytes := []byte(`{"name": "John"}`)
data, err := easyjson.Load(jsonBytes)
```

### LoadFile

```go
func LoadFile(filename string) (*JSONValue, error)
```

Read and parse JSON from file.

**Parameters:**
- `filename`: Path to JSON file

**Returns:**
- `*JSONValue`: Parsed JSON data
- `error`: File or parse error

**Example:**
```go
data, err := easyjson.LoadFile("config.json")
if err != nil {
    log.Fatal(err)
}
```

### ParseSafely

```go
func ParseSafely(jsonStr string) *ParseResult
```

Parse JSON safely, never fails, always returns valid JSONValue.

**Parameters:**
- `jsonStr`: JSON string to parse

**Returns:**
- `*ParseResult`: Contains Data, Error, and Suggestions

**Example:**
```go
result := easyjson.ParseSafely(userInput)
if result.Error != nil {
    for _, suggestion := range result.Suggestions {
        fmt.Println("Suggestion:", suggestion)
    }
}
data := result.Data // Always valid
```

### TryParse

```go
func TryParse(jsonStr string) (*JSONValue, bool)
```

Try to parse, returns success boolean.

**Parameters:**
- `jsonStr`: JSON string to parse

**Returns:**
- `*JSONValue`: Parsed data or empty object
- `bool`: true if parse succeeded

**Example:**
```go
if data, ok := easyjson.TryParse(jsonString); ok {
    // Successfully parsed
}
```

### ParseOrDefault

```go
func ParseOrDefault(jsonStr string, defaultValue *JSONValue) *JSONValue
```

Parse JSON or return default on error.

**Parameters:**
- `jsonStr`: JSON string to parse
- `defaultValue`: Value to return on error

**Returns:**
- `*JSONValue`: Parsed data or default

**Example:**
```go
data := easyjson.ParseOrDefault(jsonString, easyjson.NewObject())
```

## Creation Functions

### New

```go
func New(data interface{}) *JSONValue
```

Create JSONValue from any Go value.

**Example:**
```go
data := easyjson.New(map[string]interface{}{
    "name": "John",
    "age":  30,
})
```

### NewObject

```go
func NewObject() *JSONValue
```

Create empty JSON object.

**Example:**
```go
obj := easyjson.NewObject()
obj.Set("key", "value")
```

### NewArray

```go
func NewArray() *JSONValue
```

Create empty JSON array.

**Example:**
```go
arr := easyjson.NewArray()
arr.Append("item")
```

### NewObjectFrom

```go
func NewObjectFrom(obj map[string]interface{}) *JSONValue
```

Create JSONValue from map.

**Example:**
```go
obj := easyjson.NewObjectFrom(map[string]interface{}{
    "name": "John",
    "age":  30,
})
```

### NewArrayFrom

```go
func NewArrayFrom(items []interface{}) *JSONValue
```

Create JSONValue from slice.

**Example:**
```go
arr := easyjson.NewArrayFrom([]interface{}{"a", "b", "c"})
```

## Access Methods

### Get

```go
func (jv *JSONValue) Get(key interface{}) *JSONValue
```

Get value by key (object) or index (array).

**Parameters:**
- `key`: string for objects, int for arrays

**Returns:**
- `*JSONValue`: Value or null if not found

**Example:**
```go
name := data.Get("name")
first := data.Get(0)
```

### Q

```go
func (jv *JSONValue) Q(keys ...interface{}) *JSONValue
```

Fluent query with mixed key types.

**Parameters:**
- `keys`: Variable number of string or int keys

**Returns:**
- `*JSONValue`: Value or null if not found

**Example:**
```go
email := data.Q("users", 0, "profile", "email")
```

### Path

```go
func (jv *JSONValue) Path(path string) *JSONValue
```

Access with dot-separated path.

**Parameters:**
- `path`: Dot-separated path (e.g., "user.profile.email")

**Returns:**
- `*JSONValue`: Value or null if not found

**Example:**
```go
email := data.Path("user.profile.email")
age := data.Path("users.0.age")
```

### GetString

```go
func (jv *JSONValue) GetString(keys ...interface{}) string
```

Safe string access with default.

**Parameters:**
- `keys`: Path keys, last parameter is default if string

**Returns:**
- `string`: Value or default

**Example:**
```go
name := data.GetString("user", "name", "Anonymous")
```

### GetInt

```go
func (jv *JSONValue) GetInt(keys ...interface{}) int
```

Safe integer access with default.

**Example:**
```go
age := data.GetInt("user", "age", 0)
```

### GetBool

```go
func (jv *JSONValue) GetBool(keys ...interface{}) bool
```

Safe boolean access with default.

**Example:**
```go
active := data.GetBool("user", "active", false)
```

### GetFloat

```go
func (jv *JSONValue) GetFloat(keys ...interface{}) float64
```

Safe float access with default.

**Example:**
```go
rating := data.GetFloat("product", "rating", 0.0)
```

### GetOr

```go
func (jv *JSONValue) GetOr(keys ...interface{}) interface{}
```

Get with smart type-matched default.

**Example:**
```go
value := data.GetOr("user", "name", "Default")
```

## Type Checking Methods

### IsNull

```go
func (jv *JSONValue) IsNull() bool
```

Check if value is null or missing.

### IsString

```go
func (jv *JSONValue) IsString() bool
```

Check if value is a string.

### IsNumber

```go
func (jv *JSONValue) IsNumber() bool
```

Check if value is a number.

### IsBool

```go
func (jv *JSONValue) IsBool() bool
```

Check if value is a boolean.

### IsArray

```go
func (jv *JSONValue) IsArray() bool
```

Check if value is an array.

### IsObject

```go
func (jv *JSONValue) IsObject() bool
```

Check if value is an object.

## Type Conversion Methods

### AsString

```go
func (jv *JSONValue) AsString() string
```

Convert to string, returns "" if not string.

### AsInt

```go
func (jv *JSONValue) AsInt() int
```

Convert to int, returns 0 if not number.

### AsFloat

```go
func (jv *JSONValue) AsFloat() float64
```

Convert to float64, returns 0.0 if not number.

### AsBool

```go
func (jv *JSONValue) AsBool() bool
```

Convert to bool, returns false if not boolean.

### AsArray

```go
func (jv *JSONValue) AsArray() []*JSONValue
```

Convert to slice, returns empty slice if not array.

### AsObject

```go
func (jv *JSONValue) AsObject() map[string]*JSONValue
```

Convert to map, returns empty map if not object.

### Raw

```go
func (jv *JSONValue) Raw() interface{}
```

Get underlying Go value.

## Modification Methods

### Set

```go
func (jv *JSONValue) Set(key interface{}, value interface{}) error
```

Set value by key or index.

**Example:**
```go
data.Set("name", "John")
data.Set(0, "first item")
```

### SetPath

```go
func (jv *JSONValue) SetPath(path string, value interface{}) error
```

Set nested value, creates intermediate objects.

**Example:**
```go
data.SetPath("user.profile.email", "john@example.com")
```

### Delete

```go
func (jv *JSONValue) Delete(key interface{}) error
```

Remove key or index.

**Example:**
```go
data.Delete("oldField")
data.Delete(0)
```

### Append

```go
func (jv *JSONValue) Append(value interface{}) error
```

Append to array.

**Example:**
```go
arr.Append("new item")
```

### Extend

```go
func (jv *JSONValue) Extend(values []interface{}) error
```

Append multiple items to array.

**Example:**
```go
arr.Extend([]interface{}{"item1", "item2"})
```

### Update

```go
func (jv *JSONValue) Update(other *JSONValue) error
```

Merge another object into this one.

**Example:**
```go
data.Update(newData)
```

## Serialization Methods

### Dumps

```go
func (jv *JSONValue) Dumps() (string, error)
```

Convert to JSON string.

**Example:**
```go
jsonStr, err := data.Dumps()
```

### DumpsIndent

```go
func (jv *JSONValue) DumpsIndent(indent string) (string, error)
```

Convert to pretty JSON string.

**Example:**
```go
prettyJSON, err := data.DumpsIndent("  ")
```

### SaveFile

```go
func (jv *JSONValue) SaveFile(filename string) error
```

Save to file.

**Example:**
```go
err := data.SaveFile("output.json")
```

### SaveFileIndent

```go
func (jv *JSONValue) SaveFileIndent(filename string, indent string) error
```

Save pretty-printed to file.

**Example:**
```go
err := data.SaveFileIndent("output.json", "  ")
```

## Utility Methods

### Clone

```go
func (jv *JSONValue) Clone() *JSONValue
```

Create deep copy.

### Keys

```go
func (jv *JSONValue) Keys() []string
```

Get all object keys.

### Values

```go
func (jv *JSONValue) Values() []*JSONValue
```

Get all values.

### Items

```go
func (jv *JSONValue) Items() map[string]*JSONValue
```

Get key-value pairs.

### Len

```go
func (jv *JSONValue) Len() int
```

Get length of array or object.

### Has

```go
func (jv *JSONValue) Has(key interface{}) bool
```

Check if key or index exists.

## Array Operations

See [Array Operations](array-operations.md) for detailed documentation.

## Multi-Path Access

See [Multi-Path Access](multi-path-access.md) for detailed documentation.

## Pattern Extractors

See [Pattern Extractors](pattern-extractors.md) for detailed documentation.

## Builder API

See [JSON Building](json-building.md) for detailed documentation.
