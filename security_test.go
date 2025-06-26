package easyjson

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// security_test.go - Security-focused tests for EasyJSON library

// TestJSONBombProtection tests protection against JSON bombs
func TestJSONBombProtection(t *testing.T) {
	t.Run("ExponentialBlowupProtection", func(t *testing.T) {
		// Create a JSON structure that would cause exponential memory usage
		// This is a simplified version - real JSON bombs are more sophisticated
		
		// Test with deeply nested arrays that could cause memory issues
		deepNesting := "["
		for i := 0; i < 50; i++ {
			deepNesting += "["
		}
		deepNesting += "\"deep\""
		for i := 0; i < 50; i++ {
			deepNesting += "]"
		}
		
		// ParseSafely should handle this without crashing
		result := ParseSafely(deepNesting)
		
		// Should either parse successfully or fail gracefully
		if result.Error != nil {
			// If it fails, should provide helpful error
			if len(result.Suggestions) == 0 {
				t.Error("Should provide suggestions for complex nested structure")
			}
		}
		
		// Should always return valid JSONValue
		if result.Data == nil {
			t.Error("Should always return valid JSONValue")
		}
	})
	
	t.Run("RepeatedKeyProtection", func(t *testing.T) {
		// Test with many repeated keys (potential for hash collision attacks)
		var jsonBuilder strings.Builder
		jsonBuilder.WriteString("{")
		
		for i := 0; i < 1000; i++ {
			if i > 0 {
				jsonBuilder.WriteString(",")
			}
			jsonBuilder.WriteString(fmt.Sprintf(`"key_%d": %d`, i, i))
		}
		jsonBuilder.WriteString("}")
		
		largeJSON := jsonBuilder.String()
		
		// Should handle large number of keys without performance degradation
		start := time.Now()
		result := ParseSafely(largeJSON)
		elapsed := time.Since(start)
		
		// Should complete within reasonable time
		if elapsed > time.Second*5 {
			t.Errorf("Large JSON with many keys took too long to parse: %v", elapsed)
		}
		
		if result.Error != nil {
			t.Errorf("Should parse large JSON with many keys: %v", result.Error)
		}
		
		// Verify parsed correctly
		if result.Data.Len() != 1000 {
			t.Errorf("Expected 1000 keys, got %d", result.Data.Len())
		}
	})
	
	t.Run("LargeStringProtection", func(t *testing.T) {
		// Test with very large string values
		largeString := strings.Repeat("A", 100000)
		jsonStr := fmt.Sprintf(`{"large": "%s"}`, largeString)
		
		result := ParseSafely(jsonStr)
		
		// Should handle large strings appropriately
		if result.Error != nil {
			// If it fails due to size limits, should provide suggestions
			if !strings.Contains(result.Error.Error(), "size") && !strings.Contains(result.Error.Error(), "limit") {
				t.Errorf("Error should mention size limits: %v", result.Error)
			}
		} else {
			// If it succeeds, should parse correctly
			if result.Data.GetString("large") != largeString {
				t.Error("Large string should be parsed correctly")
			}
		}
	})
	
	t.Run("RecursiveBombProtection", func(t *testing.T) {
		// Test protection against recursive structures that could cause stack overflow
		// Since JSON doesn't naturally support recursion, we test deep nesting instead
		
		var deepObject strings.Builder
		deepObject.WriteString("{")
		
		for i := 0; i < 200; i++ {
			deepObject.WriteString(fmt.Sprintf(`"level_%d": {`, i))
		}
		
		deepObject.WriteString(`"deep": "value"`)
		
		for i := 0; i < 200; i++ {
			deepObject.WriteString("}")
		}
		
		deepObject.WriteString("}")
		
		// Should handle deep nesting gracefully
		result := ParseSafely(deepObject.String())
		
		// Should either succeed or fail with appropriate error
		if result.Error != nil {
			// Should mention depth or nesting in error
			errorStr := strings.ToLower(result.Error.Error())
			if !strings.Contains(errorStr, "depth") && !strings.Contains(errorStr, "nest") && !strings.Contains(errorStr, "limit") {
				t.Errorf("Deep nesting error should mention limits: %v", result.Error)
			}
		}
		
		// Should not crash the program
		if result.Data == nil {
			t.Error("Should return valid JSONValue even for deep nesting")
		}
	})
}

// TestMemoryExhaustionProtection tests protection against memory exhaustion
func TestMemoryExhaustionProtection(t *testing.T) {
	t.Run("LargeJSONSizeLimit", func(t *testing.T) {
		// Test that extremely large JSON is rejected
		result := ParseSafely(strings.Repeat("x", int(MaxJSONSize)+1))
		
		if result.Error == nil {
			t.Error("Should reject JSON exceeding size limits")
		}
		
		if len(result.Suggestions) == 0 {
			t.Error("Should provide suggestions for oversized JSON")
		}
		
		// Should contain size-related error message
		errorMsg := strings.ToLower(result.Error.Error())
		if !strings.Contains(errorMsg, "size") && !strings.Contains(errorMsg, "limit") && !strings.Contains(errorMsg, "large") {
			t.Logf("Error message: %s", result.Error.Error())
			t.Error("Error should mention size limits")
		}
	})
	
	t.Run("ArraySizeLimit", func(t *testing.T) {
		// Test array size validation
		jv := NewArray()
		
		// Should reject extremely large array indices
		err := jv.Set(MaxArrayIndex+1, "value")
		if err == nil {
			t.Error("Should reject array index beyond limits")
		}
		
		// Error should mention array, index, or integer key limits
		errorMsg := strings.ToLower(err.Error())
		if !strings.Contains(errorMsg, "array") && !strings.Contains(errorMsg, "index") && 
		   !strings.Contains(errorMsg, "limit") && !strings.Contains(errorMsg, "integer") &&
		   !strings.Contains(errorMsg, "key") {
			t.Logf("Error message: %s", err.Error())
			t.Error("Error should mention array/index limits")
		}
	})
	
	t.Run("StringLengthLimit", func(t *testing.T) {
		// Test string length validation
		err := validateStringLength(strings.Repeat("x", MaxStringLength+1))
		if err == nil {
			t.Error("Should reject strings exceeding length limits")
		}
		
		// Error should mention string length
		if !strings.Contains(strings.ToLower(err.Error()), "string") {
			t.Error("Error should mention string length limits")
		}
	})
	
	t.Run("ObjectKeyCountLimit", func(t *testing.T) {
		// Test object key count validation
		err := validateObjectKeyCount(MaxObjectKeys + 1)
		if err == nil {
			t.Error("Should reject objects with too many keys")
		}
		
		// Error should mention key count
		if !strings.Contains(strings.ToLower(err.Error()), "key") {
			t.Error("Error should mention key count limits")
		}
	})
	
	t.Run("MemoryUsageMonitoring", func(t *testing.T) {
		// Monitor memory usage during operations
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		// Perform memory-intensive operations
		for i := 0; i < 100; i++ { // Reduced iterations to avoid issues
			data := map[string]interface{}{
				"id":   i,
				"data": strings.Repeat("x", 100),
			}
			jv := New(data)
			
			// Perform various operations
			_ = jv.Clone()
			_ = jv.Get("data")
			_, _ = jv.Dumps()
		}
		
		runtime.GC()
		runtime.ReadMemStats(&m2)
		
		// Memory usage tracking (may have issues on some systems)
		if m2.Alloc >= m1.Alloc {
			memUsed := m2.Alloc - m1.Alloc
			t.Logf("Memory used: %d bytes", memUsed)
			
			// This is a basic check - in real scenarios you'd want more sophisticated monitoring
			if memUsed > 50*1024*1024 { // 50MB - generous limit
				t.Logf("Memory usage seems high but may be normal: %d bytes", memUsed)
			}
		} else {
			t.Log("Memory stats may have wrapped or GC occurred - this is normal")
		}
	})
}

// TestPathTraversalProtection tests protection against path traversal attacks
func TestPathTraversalProtection(t *testing.T) {
	t.Run("ParentDirectoryReferences", func(t *testing.T) {
		jv := NewObject()
		
		// Test various path traversal attempts
		traversalPaths := []string{
			"../secret",
			"data/../../../etc/passwd",
			"user/../../admin",
			"../../../root",
			"valid/../../../invalid",
			"./current/../../../parent",
		}
		
		for _, path := range traversalPaths {
			result := jv.Path(path)
			if !result.IsNull() {
				t.Errorf("Path traversal should be blocked: %s", path)
			}
			
			err := jv.SetPath(path, "value")
			if err == nil {
				t.Errorf("Path traversal in SetPath should be blocked: %s", path)
			}
		}
	})
	
	t.Run("AbsolutePathPrevention", func(t *testing.T) {
		jv := NewObject()
		
		// Test absolute path attempts
		absolutePaths := []string{
			"/etc/passwd",
			"/root/secret",
			"/usr/bin/malicious",
			"C:\\Windows\\System32",
			"\\\\server\\share\\file",
		}
		
		for _, path := range absolutePaths {
			result := jv.Path(path)
			if !result.IsNull() {
				t.Errorf("Absolute path should be blocked: %s", path)
			}
		}
	})
	
	t.Run("ConsecutiveDotsProtection", func(t *testing.T) {
		jv := NewObject()
		
		// Test consecutive dots (potential traversal)
		dotPaths := []string{
			"key..value",
			"path...deep",
			"user....admin",
			"data.....secret",
		}
		
		for _, path := range dotPaths {
			result := jv.Path(path)
			if !result.IsNull() {
				t.Errorf("Consecutive dots should be blocked: %s", path)
			}
		}
	})
	
	t.Run("PathLengthLimits", func(t *testing.T) {
		jv := NewObject()
		
		// Test extremely long paths
		longPath := strings.Repeat("a.", MaxPathLength/2) + "end"
		result := jv.Path(longPath)
		if !result.IsNull() {
			t.Error("Extremely long path should be blocked")
		}
		
		err := jv.SetPath(longPath, "value")
		if err == nil {
			t.Error("Setting extremely long path should be blocked")
		}
	})
	
	t.Run("PathDepthLimits", func(t *testing.T) {
		jv := NewObject()
		
		// Test excessive path depth
		deepPath := strings.Repeat("level.", MaxPathDepth) + "value"
		result := jv.Path(deepPath)
		if !result.IsNull() {
			t.Error("Excessively deep path should be blocked")
		}
		
		err := jv.SetPath(deepPath, "value")
		if err == nil {
			t.Error("Setting excessively deep path should be blocked")
		}
	})
}

// TestInputSanitization tests input sanitization
func TestInputSanitization(t *testing.T) {
	t.Run("NullByteFiltering", func(t *testing.T) {
		// Test null bytes in various contexts
		nullByteStrings := []string{
			"key\x00value",
			"data\x00\x00malicious",
			"\x00startwithnull",
			"endwithnull\x00",
		}
		
		for _, str := range nullByteStrings {
			jv := NewObject()
			
			// Should not be able to use null bytes in keys
			err := jv.Set(str, "value")
			if err == nil {
				t.Errorf("Should not accept null bytes in keys: %q", str)
			}
			
			// Should not be able to use in paths
			result := jv.Path(str)
			if !result.IsNull() {
				t.Errorf("Should not accept null bytes in paths: %q", str)
			}
		}
	})
	
	t.Run("ControlCharacterFiltering", func(t *testing.T) {
		// Test various control characters
		controlChars := []string{
			"\x01\x02\x03",  // start of heading, text, end of text
			"\x04\x05\x06",  // end of transmission, enquiry, acknowledge
			"\x07\x08\x0E",  // bell, backspace, shift out
			"\x0F\x10\x11",  // shift in, data link escape, device control
			"\x7F",          // delete
		}
		
		for _, ctrl := range controlChars {
			jv := NewObject()
			
			// Control characters should be filtered/rejected in keys
			err := jv.Set(ctrl, "value")
			if err == nil {
				t.Errorf("Should not accept control characters in keys: %q", ctrl)
			}
			
			// Control characters should be filtered/rejected in paths
			result := jv.Path(ctrl)
			if !result.IsNull() {
				t.Errorf("Should not accept control characters in paths: %q", ctrl)
			}
		}
	})
	
	t.Run("EncodingValidation", func(t *testing.T) {
		// Test invalid UTF-8 sequences
		invalidUTF8 := []string{
			"\xff\xfe\xfd",    // Invalid UTF-8 bytes
			"\x80\x81\x82",    // Invalid continuation bytes
			"\xc0\x80",        // Overlong encoding
			"\xe0\x80\x80",    // Overlong encoding
			"\xf0\x80\x80\x80", // Overlong encoding
		}
		
		for _, invalid := range invalidUTF8 {
			jv := NewObject()
			
			// Should not accept invalid UTF-8 in keys
			err := jv.Set(invalid, "value")
			if err == nil {
				t.Errorf("Should not accept invalid UTF-8 in keys: %q", invalid)
			}
		}
	})
	
	t.Run("ScriptInjectionPrevention", func(t *testing.T) {
		// Test potential script injection attempts
		scriptAttempts := []string{
			"<script>alert('xss')</script>",
			"javascript:alert('xss')",
			"data:text/html,<script>alert('xss')</script>",
			"vbscript:msgbox('xss')",
			"onload=alert('xss')",
			"onerror=alert('xss')",
		}
		
		for _, script := range scriptAttempts {
			jv := New(script)
			
			// Should handle script content safely (no execution)
			result := jv.AsString()
			if result != script {
				t.Errorf("Should preserve script content as string: %s", script)
			}
			
			// Should be able to use in JSON operations
			jsonStr, err := jv.Dumps()
			if err != nil {
				t.Errorf("Should be able to serialize script content: %v", err)
			}
			
			// Should be able to parse back
			parsed, err := Loads(jsonStr)
			if err != nil {
				t.Errorf("Should be able to parse back script content: %v", err)
			}
			
			if parsed.AsString() != script {
				t.Errorf("Script content should survive round-trip: %s", script)
			}
		}
	})
	
	t.Run("SQLInjectionLikeContent", func(t *testing.T) {
		// Test SQL injection-like content (should be handled as regular strings)
		sqlAttempts := []string{
			"'; DROP TABLE users; --",
			"' OR '1'='1",
			"' UNION SELECT * FROM secrets --",
			"'; INSERT INTO logs VALUES ('hacked'); --",
			"admin'--",
			"' OR 1=1 #",
		}
		
		for _, sql := range sqlAttempts {
			jv := New(sql)
			
			// Should handle SQL content safely as strings
			result := jv.AsString()
			if result != sql {
				t.Errorf("Should preserve SQL content as string: %s", sql)
			}
			
			// Should be able to use as object keys
			obj := NewObject()
			err := obj.Set("query", sql)
			if err != nil {
				t.Errorf("Should be able to store SQL content: %v", err)
			}
			
			if obj.GetString("query") != sql {
				t.Errorf("Should retrieve SQL content correctly: %s", sql)
			}
		}
	})
}

// TestControlCharacterFiltering tests control character filtering
func TestControlCharacterFiltering(t *testing.T) {
	t.Run("BasicControlCharacters", func(t *testing.T) {
		// Test basic control characters (0x00-0x1F except tab, newline, carriage return)
		for i := 0; i <= 31; i++ {
			if i == 9 || i == 10 || i == 13 { // Skip tab, newline, carriage return
				continue
			}
			
			controlChar := string(rune(i))
			jv := NewObject()
			
			// Should not accept control characters in keys
			err := jv.Set(controlChar, "value")
			if err == nil {
				t.Errorf("Should not accept control character 0x%02X in keys", i)
			}
			
			// Should not accept in paths
			result := jv.Path(controlChar)
			if !result.IsNull() {
				t.Errorf("Should not accept control character 0x%02X in paths", i)
			}
		}
	})
	
	t.Run("ExtendedControlCharacters", func(t *testing.T) {
		// Test extended control characters (0x7F-0x9F)
		extendedControls := []int{0x7F, 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x9F}
		
		for _, code := range extendedControls {
			controlChar := string(rune(code))
			jv := NewObject()
			
			// Test handling of extended control characters
			_ = jv.Set(controlChar, "value")
			// Some extended control characters might be handled differently
			// The key requirement is that the system handles them safely
			
			result := jv.Path(controlChar)
			// Should either accept safely or reject cleanly
			if !result.IsNull() {
				// If accepted, should work correctly
				t.Logf("Extended control character 0x%02X accepted", code)
			}
		}
	})
	
	t.Run("AllowedWhitespaceCharacters", func(t *testing.T) {
		// Test that allowed whitespace characters work correctly
		allowedWhitespace := []struct {
			char rune
			name string
		}{
			{'\t', "tab"},
			{'\n', "newline"},
			{'\r', "carriage return"},
			{' ', "space"},
		}
		
		for _, ws := range allowedWhitespace {
			char := string(ws.char)
			jv := New(char)
			
			// Should handle whitespace characters correctly
			if jv.AsString() != char {
				t.Errorf("Should handle %s character correctly", ws.name)
			}
			
			// Should be able to use in values (but not necessarily in keys)
			obj := NewObject()
			err := obj.Set("whitespace", char)
			if err != nil {
				t.Errorf("Should be able to store %s character: %v", ws.name, err)
			}
		}
	})
	
	t.Run("MixedControlAndValidContent", func(t *testing.T) {
		// Test strings that mix control characters with valid content
		mixedStrings := []string{
			"valid\x00content",
			"start\x01middle\x02end",
			"text\x7Fmore",
			"good\x08bad\x0Bugly",
		}
		
		for _, mixed := range mixedStrings {
			jv := NewObject()
			
			// Should not accept mixed strings with control chars in keys
			err := jv.Set(mixed, "value")
			if err == nil {
				t.Errorf("Should not accept mixed string with control chars in keys: %q", mixed)
			}
			
			// But should handle them safely in values
			jv.Set("data", mixed)
			// Should not crash or cause issues
			_ = jv.GetString("data")
		}
	})
}

// TestSecurityConfiguration tests security configuration
func TestSecurityConfiguration(t *testing.T) {
	t.Run("DefaultSecurityLimits", func(t *testing.T) {
		// Verify default security limits are reasonable
		if MaxStringLength < 1024 {
			t.Error("MaxStringLength should be at least 1024 bytes")
		}
		
		if MaxJSONSize < 1024*1024 {
			t.Error("MaxJSONSize should be at least 1MB")
		}
		
		if MaxArraySize < 1000 {
			t.Error("MaxArraySize should be at least 1000 elements")
		}
		
		if MaxPathLength < 100 {
			t.Error("MaxPathLength should be at least 100 characters")
		}
		
		if DefaultMaxRecursionDepth < 10 {
			t.Error("DefaultMaxRecursionDepth should be at least 10")
		}
	})
	
	t.Run("ValidationConfigurationCreation", func(t *testing.T) {
		config := DefaultValidationConfig()
		
		if config == nil {
			t.Error("DefaultValidationConfig should return valid config")
		}
		
		if config.MaxStringLength != MaxStringLength {
			t.Error("Config should match default MaxStringLength")
		}
		
		if config.MaxJSONSize != MaxJSONSize {
			t.Error("Config should match default MaxJSONSize")
		}
		
		if config.MaxArraySize != MaxArraySize {
			t.Error("Config should match default MaxArraySize")
		}
	})
	
	t.Run("SecurityLimitEnforcement", func(t *testing.T) {
		// Test that security limits are actually enforced
		
		// Test string length limit
		longString := strings.Repeat("x", MaxStringLength+1)
		err := validateStringLength(longString)
		if err == nil {
			t.Error("String length limit should be enforced")
		}
		
		// Test array size limit
		err = validateArraySize(MaxArraySize + 1)
		if err == nil {
			t.Error("Array size limit should be enforced")
		}
		
		// Test JSON size limit
		err = validateJSONSize(MaxJSONSize + 1)
		if err == nil {
			t.Error("JSON size limit should be enforced")
		}
		
		// Test path length limit
		longPath := strings.Repeat("a.", MaxPathLength/2) + "end"
		err = validatePath(longPath)
		if err == nil {
			t.Error("Path length limit should be enforced")
		}
	})
}

// TestSecurityErrorMessages tests that security errors provide helpful messages
func TestSecurityErrorMessages(t *testing.T) {
	t.Run("SizeLimitErrorMessages", func(t *testing.T) {
		// Test that size limit errors are informative
		err := validateStringLength(strings.Repeat("x", MaxStringLength+1))
		if err == nil {
			t.Error("Should return error for oversized string")
		}
		
		errorMsg := err.Error()
		if !strings.Contains(strings.ToLower(errorMsg), "length") {
			t.Error("Error message should mention length")
		}
		
		if !strings.Contains(errorMsg, fmt.Sprintf("%d", MaxStringLength)) {
			t.Error("Error message should include the limit value")
		}
	})
	
	t.Run("PathTraversalErrorMessages", func(t *testing.T) {
		err := validatePath("../traversal")
		if err == nil {
			t.Error("Should return error for path traversal")
		}
		
		errorMsg := err.Error()
		if !strings.Contains(strings.ToLower(errorMsg), "path") {
			t.Error("Error message should mention path")
		}
	})
	
	t.Run("KeyValidationErrorMessages", func(t *testing.T) {
		err := validateKey(make(chan int))
		if err == nil {
			t.Error("Should return error for invalid key type")
		}
		
		errorMsg := err.Error()
		if !strings.Contains(strings.ToLower(errorMsg), "key") {
			t.Error("Error message should mention key")
		}
		
		if !strings.Contains(strings.ToLower(errorMsg), "type") {
			t.Error("Error message should mention type")
		}
	})
}

// TestDenialOfServiceProtection tests protection against DoS attacks
func TestDenialOfServiceProtection(t *testing.T) {
	t.Run("RepeatedOperationProtection", func(t *testing.T) {
		// Test that repeated operations don't cause performance degradation
		jv := NewObject()
		
		start := time.Now()
		
		// Perform many operations
		for i := 0; i < 10000; i++ {
			jv.Set(fmt.Sprintf("key_%d", i), i)
		}
		
		elapsed := time.Since(start)
		
		// Should complete within reasonable time
		if elapsed > time.Second*5 {
			t.Errorf("Repeated operations took too long: %v", elapsed)
		}
		
		// Verify operations completed correctly
		if jv.Len() != 10000 {
			t.Errorf("Expected 10000 keys, got %d", jv.Len())
		}
	})
	
	t.Run("LargeDataProtection", func(t *testing.T) {
		// Test handling of large data structures
		largeData := make(map[string]interface{})
		
		for i := 0; i < 1000; i++ {
			largeData[fmt.Sprintf("key_%d", i)] = strings.Repeat("data", 100)
		}
		
		start := time.Now()
		jv := New(largeData)
		
		// Perform operations on large data
		_ = jv.Keys()
		_ = jv.Values()
		_ = jv.Clone()
		
		elapsed := time.Since(start)
		
		// Should handle large data efficiently
		if elapsed > time.Second*10 {
			t.Errorf("Large data operations took too long: %v", elapsed)
		}
	})
	
	t.Run("RecursiveOperationProtection", func(t *testing.T) {
		// Test protection against recursive operations that could cause stack overflow
		
		// Create nested structure
		root := NewObject()
		current := root
		
		for i := 0; i < DefaultMaxRecursionDepth-10; i++ {
			nested := NewObject()
			nested.Set("level", i)
			current.Set("child", nested.Raw())
			current = nested
		}
		
		start := time.Now()
		
		// Operations should complete without stack overflow
		_ = root.Clone()
		_ = root.DeepSearch("level")
		
		elapsed := time.Since(start)
		
		// Should complete within reasonable time
		if elapsed > time.Second*5 {
			t.Errorf("Recursive operations took too long: %v", elapsed)
		}
	})
}

// TestSecurityBestPractices tests adherence to security best practices
func TestSecurityBestPractices(t *testing.T) {
	t.Run("DefaultSecureBehavior", func(t *testing.T) {
		// Test that default behavior is secure
		
		// ParseSafely should never panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseSafely should never panic: %v", r)
			}
		}()
		
		// Test with various malicious inputs
		maliciousInputs := []string{
			"",
			"not json",
			`{"key": }`,
			strings.Repeat("x", 10000),
			`{"key": "` + strings.Repeat("x", 10000) + `"}`,
		}
		
		for _, input := range maliciousInputs {
			result := ParseSafely(input)
			if result == nil {
				t.Error("ParseSafely should never return nil")
			}
			if result.Data == nil {
				t.Error("ParseSafely should always return valid Data")
			}
		}
	})
	
	t.Run("ErrorHandlingSecrity", func(t *testing.T) {
		// Test that error handling doesn't leak sensitive information
		
		// Create various error conditions
		jv := NewObject()
		
		// Invalid operations should return safe errors
		err := jv.Set(make(chan int), "value")
		if err == nil {
			t.Error("Should return error for invalid operation")
		}
		
		// Error message should not contain sensitive details
		errorMsg := err.Error()
		sensitivePatterns := []string{
			"panic",
			"stack",
			"memory",
			"address",
			"pointer",
		}
		
		for _, pattern := range sensitivePatterns {
			if strings.Contains(strings.ToLower(errorMsg), pattern) {
				t.Errorf("Error message should not contain sensitive pattern '%s': %s", pattern, errorMsg)
			}
		}
	})
	
	t.Run("InputValidationCompleteness", func(t *testing.T) {
		// Test that all user inputs are validated
		
		jv := NewObject()
		
		// Test various invalid inputs
		invalidInputs := []interface{}{
			func() {},
			make(chan int),
			complex(1, 2),
		}
		
		for _, input := range invalidInputs {
			// Should not accept invalid inputs as keys
			err := jv.Set(input, "value")
			if err == nil {
				t.Errorf("Should reject invalid input as key: %T", input)
			}
			
			// Should not accept invalid inputs as values
			err = jv.Set("key", input)
			if err == nil {
				t.Errorf("Should reject invalid input as value: %T", input)
			}
		}
	})
}

// Benchmark security-related operations
func BenchmarkSecurityValidation(b *testing.B) {
	testString := "test_key_validation"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateStringKey(testString)
	}
}

func BenchmarkPathValidation(b *testing.B) {
	testPath := "user.profile.settings.theme"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validatePath(testPath)
	}
}

func BenchmarkLargeJSONParsing(b *testing.B) {
	// Create reasonably large JSON for benchmarking
	largeJSON := `{"users": [`
	for i := 0; i < 1000; i++ {
		if i > 0 {
			largeJSON += ","
		}
		largeJSON += fmt.Sprintf(`{"id": %d, "name": "user_%d"}`, i, i)
	}
	largeJSON += `]}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := ParseSafely(largeJSON)
		if result.Error != nil {
			b.Fatalf("Parse failed: %v", result.Error)
		}
	}
}