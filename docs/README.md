# EasyJSON Documentation

Welcome to the comprehensive documentation for EasyJSON, a Python-like JSON manipulation library for Go.

## Table of Contents

1. [Getting Started](getting-started.md)
2. [Core Concepts](core-concepts.md)
3. [API Reference](api-reference.md)
4. [Array Operations](array-operations.md)
5. [Smart Getters and Safe Access](smart-getters.md)
6. [JSON Building](json-building.md)
7. [Multi-Path Access](multi-path-access.md)
8. [Safe Parsing](safe-parsing.md)
9. [Pattern Extractors](pattern-extractors.md)
10. [Validation and Security](validation-security.md)
11. [Best Practices](best-practices.md)
12. [Examples and Recipes](examples.md)
13. [Migration Guide](migration-guide.md)
14. [Performance Tips](performance.md)
15. [Troubleshooting](troubleshooting.md)

## Quick Links

- [Installation Instructions](getting-started.md#installation)
- [Quick Start Guide](getting-started.md#quick-start)
- [Common Use Cases](examples.md#common-use-cases)
- [API Cheat Sheet](api-reference.md#quick-reference)

## What is EasyJSON?

EasyJSON is a Go library that brings Python-like simplicity to JSON manipulation while maintaining Go's performance and type safety. It provides:

- Intuitive, chainable API for accessing nested JSON data
- Safe operations that never panic
- Multiple access patterns (traditional, fluent, path notation)
- Powerful array operations inspired by JavaScript and Python
- Built-in validation and security features
- Zero external dependencies

## Why EasyJSON?

### Before EasyJSON:
```go
// Complex type assertions and error-prone
userMap, ok := data["user"].(map[string]interface{})
if !ok {
    return errors.New("user not found")
}
profileMap, ok := userMap["profile"].(map[string]interface{})
if !ok {
    return errors.New("profile not found")
}
name, ok := profileMap["name"].(string)
if !ok {
    name = "Unknown"
}
```

### With EasyJSON:
```go
// Simple, safe, and elegant
name := data.GetString("user", "profile", "name", "Unknown")
```

## Getting Help

- Check the [Troubleshooting Guide](troubleshooting.md) for common issues
- Review [Examples](examples.md) for real-world use cases
- Read [Best Practices](best-practices.md) for optimal usage
- Open an issue on [GitHub](https://github.com/javanhut/easyjson/issues)

## Documentation Structure

Each documentation file is self-contained and covers a specific aspect of EasyJSON:

- **Getting Started**: Installation, basic usage, and first steps
- **Core Concepts**: Understanding JSONValue, access patterns, and philosophy
- **API Reference**: Complete method documentation with examples
- **Specialized Guides**: Deep dives into specific features
- **Practical Guides**: Examples, recipes, and best practices

## Contributing to Documentation

Found an error or want to improve the docs? Contributions are welcome! Please see our contributing guidelines.

## Version

This documentation is for EasyJSON v2.0.0 and later.
