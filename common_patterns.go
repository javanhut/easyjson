package easyjson

import (
	"fmt"
	"strings"
	"time"
)

// common_patterns.go - Common JSON pattern extractors

// GetUserInfo extracts user information from various JSON structures
// Usage: userInfo := data.GetUserInfo()
func (jv *JSONValue) GetUserInfo() map[string]string {
	return map[string]string{
		"id": jv.TryPaths("id", "user_id", "userId", "ID").AsString(),
		"name": jv.TryPaths("name", "full_name", "fullName", "username", "display_name").
			AsString(),
		"email":    jv.TryPaths("email", "email_address", "emailAddress", "mail").AsString(),
		"role":     jv.TryPaths("role", "user_role", "userRole", "type", "account_type").AsString(),
		"username": jv.TryPaths("username", "user_name", "userName", "login", "handle").AsString(),
		"phone":    jv.TryPaths("phone", "phone_number", "phoneNumber", "mobile", "tel").AsString(),
	}
}

// GetPaginationInfo extracts pagination information from various structures
// Usage: pagination := data.GetPaginationInfo()
func (jv *JSONValue) GetPaginationInfo() map[string]int {
	return map[string]int{
		"page": jv.TryPaths("page", "current_page", "pageNumber", "pageNum").AsInt(),
		"total": jv.TryPaths("total", "total_count", "totalCount", "count", "total_items").
			AsInt(),
		"limit":       jv.TryPaths("limit", "page_size", "pageSize", "per_page", "size").AsInt(),
		"offset":      jv.TryPaths("offset", "start", "skip", "from").AsInt(),
		"total_pages": jv.TryPaths("total_pages", "totalPages", "page_count", "pages").AsInt(),
	}
}

// GetTimestamps extracts timestamp information
// Usage: timestamps := data.GetTimestamps()
func (jv *JSONValue) GetTimestamps() map[string]string {
	return map[string]string{
		"created": jv.TryPaths("created_at", "createdAt", "created", "date_created", "creation_date").
			AsString(),
		"updated": jv.TryPaths("updated_at", "updatedAt", "updated", "date_updated", "modification_date").
			AsString(),
		"deleted": jv.TryPaths("deleted_at", "deletedAt", "deleted", "date_deleted").AsString(),
		"published": jv.TryPaths("published_at", "publishedAt", "published", "date_published").
			AsString(),
	}
}

// GetAPIResponseInfo extracts standard API response information
// Usage: response := data.GetAPIResponseInfo()
func (jv *JSONValue) GetAPIResponseInfo() map[string]string {
	return map[string]string{
		"status":  jv.TryPaths("status", "state", "result", "success").AsString(),
		"message": jv.TryPaths("message", "msg", "description", "detail", "info").AsString(),
		"error":   jv.TryPaths("error", "error_message", "errorMessage", "err").AsString(),
		"code":    jv.TryPaths("code", "error_code", "errorCode", "status_code").AsString(),
	}
}

// GetContactInfo extracts contact information
// Usage: contact := data.GetContactInfo()
func (jv *JSONValue) GetContactInfo() map[string]string {
	return map[string]string{
		"email":   jv.TryPaths("email", "email_address", "emailAddress", "mail").AsString(),
		"phone":   jv.TryPaths("phone", "phone_number", "phoneNumber", "mobile", "tel").AsString(),
		"address": jv.TryPaths("address", "street_address", "streetAddress", "location").AsString(),
		"city":    jv.TryPaths("city", "town", "locality").AsString(),
		"state":   jv.TryPaths("state", "province", "region").AsString(),
		"country": jv.TryPaths("country", "nation").AsString(),
		"zip":     jv.TryPaths("zip", "postal_code", "postalCode", "zipcode").AsString(),
	}
}

// GetProductInfo extracts product/item information
// Usage: product := data.GetProductInfo()
func (jv *JSONValue) GetProductInfo() map[string]string {
	return map[string]string{
		"id":          jv.TryPaths("id", "product_id", "item_id", "sku").AsString(),
		"name":        jv.TryPaths("name", "title", "product_name", "item_name").AsString(),
		"description": jv.TryPaths("description", "desc", "summary", "details").AsString(),
		"price":       jv.TryPaths("price", "cost", "amount", "value").AsString(),
		"currency":    jv.TryPaths("currency", "currency_code", "symbol").AsString(),
		"category":    jv.TryPaths("category", "type", "classification").AsString(),
		"brand":       jv.TryPaths("brand", "manufacturer", "company").AsString(),
	}
}

// IsAPISuccess checks if API response indicates success
// Usage: if data.IsAPISuccess() { ... }
func (jv *JSONValue) IsAPISuccess() bool {
	status := jv.TryPaths("status", "state", "result", "success").AsString()
	status = strings.ToLower(status)

	successValues := []string{"success", "ok", "true", "1", "complete", "completed", "done"}
	for _, val := range successValues {
		if status == val {
			return true
		}
	}

	// Check boolean success field
	if jv.TryPaths("success", "ok", "is_success").AsBool() {
		return true
	}

	return false
}

// IsAPIError checks if API response indicates error
// Usage: if data.IsAPIError() { ... }
func (jv *JSONValue) IsAPIError() bool {
	status := jv.TryPaths("status", "state", "result").AsString()
	status = strings.ToLower(status)

	errorValues := []string{"error", "fail", "failed", "failure", "false", "0"}
	for _, val := range errorValues {
		if status == val {
			return true
		}
	}

	// Check if error field exists and is not empty
	errorMsg := jv.TryPaths("error", "error_message", "errorMessage").AsString()
	if errorMsg != "" {
		return true
	}

	return false
}

// HasPagination checks if data contains pagination information
// Usage: if data.HasPagination() { ... }
func (jv *JSONValue) HasPagination() bool {
	paginationFields := []string{"page", "total", "limit", "offset", "pagination", "meta"}
	for _, field := range paginationFields {
		if !jv.DeepSearch(field).IsNull() {
			return true
		}
	}
	return false
}

// GetNestedValue safely gets nested value with multiple fallback paths
// Usage: value := data.GetNestedValue("user.profile.name", "user.name", "name")
func (jv *JSONValue) GetNestedValue(paths ...string) *JSONValue {
	return jv.TryPaths(paths...)
}

// ExtractArrayField extracts a specific field from all items in an array
// Usage: names := data.Get("users").ExtractArrayField("name")
func (jv *JSONValue) ExtractArrayField(fieldName string) []string {
	return jv.PluckStrings(fieldName)
}

// CountByField counts occurrences of each value for a field in an array
// Usage: counts := data.Get("users").CountByField("role")
func (jv *JSONValue) CountByField(fieldName string) map[string]int {
	counts := make(map[string]int)

	if !jv.IsArray() {
		return counts
	}

	for _, item := range jv.AsArray() {
		value := item.GetString(fieldName)
		if value != "" {
			counts[value]++
		}
	}

	return counts
}

// GetMetadata extracts common metadata patterns
// Usage: metadata := data.GetMetadata()
func (jv *JSONValue) GetMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})

	// Common metadata fields
	metaFields := map[string][]string{
		"version":     {"version", "v", "api_version", "schema_version"},
		"timestamp":   {"timestamp", "time", "date", "created_at"},
		"source":      {"source", "origin", "from", "provider"},
		"environment": {"environment", "env", "stage", "tier"},
		"server":      {"server", "host", "hostname", "instance"},
		"request_id":  {"request_id", "requestId", "correlation_id", "trace_id"},
	}

	for key, paths := range metaFields {
		value := jv.TryPaths(paths...)
		if !value.IsNull() {
			metadata[key] = value.Raw()
		}
	}

	// Check for explicit metadata/meta objects
	if metaObj := jv.TryPaths("metadata", "meta", "_meta"); !metaObj.IsNull() {
		if metaObj.IsObject() {
			for _, key := range metaObj.Keys() {
				metadata[key] = metaObj.Get(key).Raw()
			}
		}
	}

	return metadata
}

// IsValidEmail checks if a string field contains a valid email
// Usage: if data.Get("email").IsValidEmail() { ... }
func (jv *JSONValue) IsValidEmail() bool {
	email := jv.AsString()
	if email == "" {
		return false
	}

	// Basic email validation
	return strings.Contains(email, "@") &&
		strings.Contains(email, ".") &&
		len(email) > 5 &&
		!strings.HasPrefix(email, "@") &&
		!strings.HasSuffix(email, "@")
}

// IsValidURL checks if a string field contains a valid URL
// Usage: if data.Get("website").IsValidURL() { ... }
func (jv *JSONValue) IsValidURL() bool {
	url := jv.AsString()
	if url == "" {
		return false
	}

	// Basic URL validation
	return strings.HasPrefix(strings.ToLower(url), "http://") ||
		strings.HasPrefix(strings.ToLower(url), "https://") ||
		strings.HasPrefix(strings.ToLower(url), "ftp://")
}

// IsValidDate checks if a string field contains a valid date
// Usage: if data.Get("created_at").IsValidDate() { ... }
func (jv *JSONValue) IsValidDate() bool {
	dateStr := jv.AsString()
	if dateStr == "" {
		return false
	}

	// Try common date formats
	formats := []string{
		time.RFC3339,
		time.RFC822,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"01-02-2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, dateStr); err == nil {
			return true
		}
	}

	return false
}

// GetFormattedDate returns formatted date string
// Usage: formatted := data.Get("created_at").GetFormattedDate("2006-01-02")
func (jv *JSONValue) GetFormattedDate(format string) string {
	dateStr := jv.AsString()
	if dateStr == "" {
		return ""
	}

	// Try to parse with common formats
	inputFormats := []string{
		time.RFC3339,
		time.RFC822,
		"2006-01-02",
		"2006-01-02 15:04:05",
		"01/02/2006",
		"01-02-2006",
		"2006/01/02",
	}

	for _, inputFormat := range inputFormats {
		if parsedTime, err := time.Parse(inputFormat, dateStr); err == nil {
			return parsedTime.Format(format)
		}
	}

	return dateStr // Return original if can't parse
}

// GetRelativeTime returns relative time string (e.g., "2 hours ago")
// Usage: relative := data.Get("created_at").GetRelativeTime()
func (jv *JSONValue) GetRelativeTime() string {
	dateStr := jv.AsString()
	if dateStr == "" {
		return ""
	}

	// Try to parse the date
	var parsedTime time.Time
	var err error

	inputFormats := []string{
		time.RFC3339,
		time.RFC822,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range inputFormats {
		if parsedTime, err = time.Parse(format, dateStr); err == nil {
			break
		}
	}

	if err != nil {
		return dateStr
	}

	// Calculate relative time
	now := time.Now()
	diff := now.Sub(parsedTime)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return parsedTime.Format("2006-01-02")
	}
}

// GetImageInfo extracts image/media information
// Usage: imageInfo := data.GetImageInfo()
func (jv *JSONValue) GetImageInfo() map[string]string {
	return map[string]string{
		"url":      jv.TryPaths("url", "image_url", "imageUrl", "src", "source").AsString(),
		"alt":      jv.TryPaths("alt", "alt_text", "altText", "description", "caption").AsString(),
		"width":    jv.TryPaths("width", "w", "image_width").AsString(),
		"height":   jv.TryPaths("height", "h", "image_height").AsString(),
		"format":   jv.TryPaths("format", "type", "mime_type", "content_type").AsString(),
		"size":     jv.TryPaths("size", "file_size", "fileSize", "bytes").AsString(),
		"filename": jv.TryPaths("filename", "file_name", "fileName", "name").AsString(),
	}
}

// GetLocationInfo extracts location/geographic information
// Usage: location := data.GetLocationInfo()
func (jv *JSONValue) GetLocationInfo() map[string]string {
	return map[string]string{
		"latitude":  jv.TryPaths("latitude", "lat", "y").AsString(),
		"longitude": jv.TryPaths("longitude", "lng", "lon", "x").AsString(),
		"address":   jv.TryPaths("address", "full_address", "street_address").AsString(),
		"city":      jv.TryPaths("city", "town", "locality", "municipality").AsString(),
		"state":     jv.TryPaths("state", "province", "region", "admin_area").AsString(),
		"country":   jv.TryPaths("country", "nation", "country_code").AsString(),
		"zip":       jv.TryPaths("zip", "postal_code", "postcode", "zipcode").AsString(),
		"timezone":  jv.TryPaths("timezone", "tz", "time_zone").AsString(),
	}
}

// GetSocialMediaInfo extracts social media links/handles
// Usage: social := data.GetSocialMediaInfo()
func (jv *JSONValue) GetSocialMediaInfo() map[string]string {
	return map[string]string{
		"twitter":   jv.TryPaths("twitter", "twitter_handle", "twitter_url").AsString(),
		"facebook":  jv.TryPaths("facebook", "facebook_url", "fb_url").AsString(),
		"instagram": jv.TryPaths("instagram", "instagram_handle", "ig_handle").AsString(),
		"linkedin":  jv.TryPaths("linkedin", "linkedin_url").AsString(),
		"github":    jv.TryPaths("github", "github_username", "github_url").AsString(),
		"website":   jv.TryPaths("website", "homepage", "url", "site").AsString(),
		"blog":      jv.TryPaths("blog", "blog_url").AsString(),
	}
}

// GetFinancialInfo extracts financial/monetary information
// Usage: financial := data.GetFinancialInfo()
func (jv *JSONValue) GetFinancialInfo() map[string]string {
	return map[string]string{
		"amount":   jv.TryPaths("amount", "value", "price", "cost", "total").AsString(),
		"currency": jv.TryPaths("currency", "currency_code", "symbol").AsString(),
		"tax":      jv.TryPaths("tax", "tax_amount", "vat").AsString(),
		"discount": jv.TryPaths("discount", "discount_amount", "sale_price").AsString(),
		"subtotal": jv.TryPaths("subtotal", "sub_total", "net_amount").AsString(),
		"transaction_id": jv.TryPaths("transaction_id", "txn_id", "payment_id", "order_id").
			AsString(),
		"status": jv.TryPaths("status", "payment_status", "transaction_status").AsString(),
	}
}

// HasRequiredFields checks if all specified fields exist and are not empty
// Usage: if data.HasRequiredFields("name", "email", "phone") { ... }
func (jv *JSONValue) HasRequiredFields(fields ...string) bool {
	for _, field := range fields {
		value := jv.TryPaths(field)
		if value.IsEmptyOrNull() {
			return false
		}
	}
	return true
}

// GetMissingFields returns list of missing required fields
// Usage: missing := data.GetMissingFields("name", "email", "phone")
func (jv *JSONValue) GetMissingFields(fields ...string) []string {
	var missing []string
	for _, field := range fields {
		value := jv.TryPaths(field)
		if value.IsEmptyOrNull() {
			missing = append(missing, field)
		}
	}
	return missing
}

// IsComplete checks if object has all expected fields for a given type
// Usage: if data.IsComplete("user") { ... }
func (jv *JSONValue) IsComplete(objectType string) bool {
	requiredFields := map[string][]string{
		"user":    {"id", "name", "email"},
		"product": {"id", "name", "price"},
		"order":   {"id", "user_id", "total", "status"},
		"address": {"street", "city", "state", "zip"},
		"contact": {"name", "email"},
		"event":   {"name", "date", "location"},
	}

	if fields, exists := requiredFields[strings.ToLower(objectType)]; exists {
		return jv.HasRequiredFields(fields...)
	}

	return true // Unknown type, assume complete
}

// GetCompletionScore returns completion percentage (0.0 to 1.0)
// Usage: score := data.GetCompletionScore("user")
func (jv *JSONValue) GetCompletionScore(objectType string) float64 {
	requiredFields := map[string][]string{
		"user":    {"id", "name", "email", "phone", "address"},
		"product": {"id", "name", "description", "price", "category", "image"},
		"order":   {"id", "user_id", "items", "total", "status", "date"},
		"profile": {"name", "bio", "avatar", "location", "website"},
	}

	fields, exists := requiredFields[strings.ToLower(objectType)]
	if !exists {
		return 1.0 // Unknown type, assume complete
	}

	presentCount := 0
	for _, field := range fields {
		if !jv.TryPaths(field).IsEmptyOrNull() {
			presentCount++
		}
	}

	return float64(presentCount) / float64(len(fields))
}

// SanitizeForOutput cleans data for safe output (removes sensitive fields)
// Usage: safe := data.SanitizeForOutput()
func (jv *JSONValue) SanitizeForOutput() *JSONValue {
	if !jv.IsObject() {
		return jv
	}

	sensitiveFields := []string{
		"password", "secret", "token", "key", "private",
		"ssn", "social_security", "credit_card", "cvv",
		"api_key", "access_token", "refresh_token",
		"private_key", "certificate", "hash", "salt",
	}

	cleaned := jv.Clone()

	// Remove sensitive fields
	for _, field := range sensitiveFields {
		for _, key := range cleaned.Keys() {
			if strings.Contains(strings.ToLower(key), field) {
				cleaned.Delete(key)
			}
		}
	}

	// Recursively clean nested objects
	for _, key := range cleaned.Keys() {
		child := cleaned.Get(key)
		if child.IsObject() {
			cleaned.Set(key, child.SanitizeForOutput().Raw())
		} else if child.IsArray() {
			// Clean array items if they're objects
			cleanedArray := make([]interface{}, 0)
			for i := 0; i < child.Len(); i++ {
				item := child.Get(i)
				if item.IsObject() {
					cleanedArray = append(cleanedArray, item.SanitizeForOutput().Raw())
				} else {
					cleanedArray = append(cleanedArray, item.Raw())
				}
			}
			cleaned.Set(key, cleanedArray)
		}
	}

	return cleaned
}

// GetSummary returns a summary of the JSON structure
// Usage: summary := data.GetSummary()
func (jv *JSONValue) GetSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"type": jv.TypeString(),
		"size": jv.calculateByteSize(),
	}

	if jv.IsObject() {
		summary["keys"] = len(jv.Keys())
		summary["key_names"] = jv.Keys()

		// Analyze value types
		typeCount := make(map[string]int)
		for _, key := range jv.Keys() {
			valueType := jv.Get(key).TypeString()
			typeCount[valueType]++
		}
		summary["value_types"] = typeCount

	} else if jv.IsArray() {
		summary["length"] = jv.Len()

		if jv.Len() > 0 {
			// Analyze item types
			typeCount := make(map[string]int)
			for i := 0; i < jv.Len(); i++ {
				itemType := jv.Get(i).TypeString()
				typeCount[itemType]++
			}
			summary["item_types"] = typeCount
			summary["first_item_type"] = jv.Get(0).TypeString()
		}

	} else if jv.IsString() {
		str := jv.AsString()
		summary["length"] = len(str)
		summary["is_email"] = jv.IsValidEmail()
		summary["is_url"] = jv.IsValidURL()
		summary["is_date"] = jv.IsValidDate()
	}

	return summary
}

// TypeString returns human-readable type name
func (jv *JSONValue) TypeString() string {
	switch {
	case jv.IsNull():
		return "null"
	case jv.IsString():
		return "string"
	case jv.IsNumber():
		return "number"
	case jv.IsBool():
		return "boolean"
	case jv.IsArray():
		return "array"
	case jv.IsObject():
		return "object"
	default:
		return "unknown"
	}
}

// calculateByteSize estimates the size in bytes
func (jv *JSONValue) calculateByteSize() int {
	if bytes, err := jv.Dump(); err == nil {
		return len(bytes)
	}
	return 0
}

// Pretty returns nicely formatted string representation
func (jv *JSONValue) Pretty() string {
	if result, err := jv.DumpsIndent("  "); err == nil {
		return result
	}
	return jv.String()
}

// Compact returns minified JSON string
func (jv *JSONValue) Compact() string {
	if result, err := jv.Dumps(); err == nil {
		return result
	}
	return jv.String()
}
