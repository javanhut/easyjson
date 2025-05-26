package easyjson

import (
	"sort"
	"strings"
)

// suggestions.go - AI-like intelligent path suggestions

// JSONSuggester provides intelligent path suggestions and completion
type JSONSuggester struct {
	data        *JSONValue
	commonPaths map[string]int // Track access frequency
	history     []string       // Track recent access patterns
}

// WithSuggestions creates a suggester for intelligent assistance
// Usage: smart := easyjson.WithSuggestions(data)
func WithSuggestions(jv *JSONValue) *JSONSuggester {
	return &JSONSuggester{
		data:        jv,
		commonPaths: make(map[string]int),
		history:     []string{},
	}
}

// SuggestPaths returns likely paths based on JSON structure and common patterns
// Usage: suggestions := smart.SuggestPaths()
func (js *JSONSuggester) SuggestPaths() []string {
	suggestions := []string{}

	// Common API response patterns
	commonPatterns := []string{
		"data", "result", "results", "items", "list", "payload",
		"user", "users", "profile", "account",
		"user.name", "user.email", "user.id", "user.profile",
		"status", "message", "error", "errors", "success",
		"page", "total", "limit", "count", "offset",
		"created_at", "updated_at", "timestamp", "date",
		"settings", "config", "preferences",
		"settings.theme", "settings.language", "settings.notifications",
		"meta", "metadata", "info",
		"id", "name", "title", "description", "type",
	}

	// Check which patterns exist in the data
	for _, pattern := range commonPatterns {
		if !js.data.Path(pattern).IsNull() {
			suggestions = append(suggestions, pattern)
		}
	}

	// Add top-level keys
	if js.data.IsObject() {
		for _, key := range js.data.Keys() {
			// Avoid duplicates
			found := false
			for _, existing := range suggestions {
				if existing == key {
					found = true
					break
				}
			}
			if !found {
				suggestions = append(suggestions, key)
			}
		}
	}

	// Add nested keys for objects
	js.addNestedSuggestions(js.data, "", &suggestions, 2) // Max depth 2

	// Sort by frequency if we have history
	if len(js.commonPaths) > 0 {
		sort.Slice(suggestions, func(i, j int) bool {
			freqI := js.commonPaths[suggestions[i]]
			freqJ := js.commonPaths[suggestions[j]]
			return freqI > freqJ
		})
	}

	return suggestions
}

// addNestedSuggestions recursively adds nested path suggestions
func (js *JSONSuggester) addNestedSuggestions(
	jv *JSONValue,
	prefix string,
	suggestions *[]string,
	maxDepth int,
) {
	if maxDepth <= 0 {
		return
	}

	if jv.IsObject() {
		for _, key := range jv.Keys() {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}

			// Add this path if not already present
			found := false
			for _, existing := range *suggestions {
				if existing == path {
					found = true
					break
				}
			}
			if !found {
				*suggestions = append(*suggestions, path)
			}

			// Recurse into nested objects
			child := jv.Get(key)
			js.addNestedSuggestions(child, path, suggestions, maxDepth-1)
		}
	}
}

// CompletePartial suggests completions for partial paths
// Usage: completions := smart.CompletePartial("user.")
func (js *JSONSuggester) CompletePartial(partial string) []string {
	var completions []string

	if partial == "" {
		return js.SuggestPaths()
	}

	// Split path and try to complete the last part
	parts := strings.Split(partial, ".")
	if len(parts) == 1 {
		// Complete top-level key
		if js.data.IsObject() {
			for _, key := range js.data.Keys() {
				if strings.HasPrefix(strings.ToLower(key), strings.ToLower(partial)) {
					completions = append(completions, key)
				}
			}
		}
	} else {
		// Navigate to parent and complete child key
		parentPath := strings.Join(parts[:len(parts)-1], ".")
		lastPart := parts[len(parts)-1]

		parent := js.data.Path(parentPath)
		if parent.IsObject() {
			for _, key := range parent.Keys() {
				if strings.HasPrefix(strings.ToLower(key), strings.ToLower(lastPart)) {
					completions = append(completions, parentPath+"."+key)
				}
			}
		}
	}

	return completions
}

// ValidatePathWithSuggestions checks path and suggests alternatives if invalid
// Usage: valid, suggestions := smart.ValidatePathWithSuggestions("user.nam")
func (js *JSONSuggester) ValidatePathWithSuggestions(path string) (bool, []string) {
	if !js.data.Path(path).IsNull() {
		js.trackAccess(path) // Track successful access
		return true, []string{}
	}

	// Find similar paths
	suggestions := []string{}
	allPaths := js.SuggestPaths()

	for _, validPath := range allPaths {
		if similarity(path, validPath) > 0.6 { // 60% similarity threshold
			suggestions = append(suggestions, validPath)
		}
	}

	// Sort by similarity
	sort.Slice(suggestions, func(i, j int) bool {
		simI := similarity(path, suggestions[i])
		simJ := similarity(path, suggestions[j])
		return simI > simJ
	})

	// Limit to top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return false, suggestions
}

// PredictNext predicts what the user might want to access next
// Usage: predictions := smart.PredictNext()
func (js *JSONSuggester) PredictNext() []string {
	if len(js.history) == 0 {
		return js.SuggestPaths()
	}

	predictions := []string{}

	// Look at recent access patterns
	recent := js.history
	if len(recent) > 5 {
		recent = recent[len(recent)-5:] // Last 5 accesses
	}

	// Common follow-up patterns
	patterns := map[string][]string{
		"user":     {"user.name", "user.email", "user.id", "user.profile"},
		"profile":  {"profile.name", "profile.age", "profile.avatar"},
		"settings": {"settings.theme", "settings.language", "settings.notifications"},
		"data":     {"data.items", "data.total", "data.page"},
		"result":   {"result.status", "result.message", "result.data"},
		"error":    {"error.message", "error.code", "error.details"},
	}

	// Check patterns for recent accesses
	for _, recentPath := range recent {
		parts := strings.Split(recentPath, ".")
		for _, part := range parts {
			if followUps, exists := patterns[part]; exists {
				for _, followUp := range followUps {
					if !js.data.Path(followUp).IsNull() {
						predictions = append(predictions, followUp)
					}
				}
			}
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	unique := []string{}
	for _, pred := range predictions {
		if !seen[pred] {
			seen[pred] = true
			unique = append(unique, pred)
		}
	}

	return unique
}

// GetSmartRecommendations provides contextual recommendations
// Usage: recommendations := smart.GetSmartRecommendations()
func (js *JSONSuggester) GetSmartRecommendations() map[string][]string {
	recommendations := make(map[string][]string)

	// User-related recommendations
	if !js.data.DeepSearch("user").IsNull() || !js.data.DeepSearch("users").IsNull() {
		userPaths := []string{}
		for _, path := range js.SuggestPaths() {
			if strings.Contains(path, "user") {
				userPaths = append(userPaths, path)
			}
		}
		if len(userPaths) > 0 {
			recommendations["User Data"] = userPaths
		}
	}

	// API response recommendations
	apiPaths := []string{}
	apiKeys := []string{"status", "message", "data", "result", "error", "success"}
	for _, key := range apiKeys {
		if !js.data.Get(key).IsNull() {
			apiPaths = append(apiPaths, key)
		}
	}
	if len(apiPaths) > 0 {
		recommendations["API Response"] = apiPaths
	}

	// Pagination recommendations
	paginationPaths := []string{}
	paginationKeys := []string{"page", "total", "limit", "offset", "count"}
	for _, key := range paginationKeys {
		path := js.data.FindPath(key)
		if path != "" {
			paginationPaths = append(paginationPaths, path)
		}
	}
	if len(paginationPaths) > 0 {
		recommendations["Pagination"] = paginationPaths
	}

	// Settings recommendations
	settingsPaths := []string{}
	settingsKeys := []string{"settings", "config", "preferences", "theme", "language"}
	for _, key := range settingsKeys {
		path := js.data.FindPath(key)
		if path != "" {
			settingsPaths = append(settingsPaths, path)
		}
	}
	if len(settingsPaths) > 0 {
		recommendations["Settings"] = settingsPaths
	}

	return recommendations
}

// TrackAccess records path access for learning
func (js *JSONSuggester) trackAccess(path string) {
	js.commonPaths[path]++
	js.history = append(js.history, path)

	// Keep history limited
	if len(js.history) > 100 {
		js.history = js.history[50:] // Keep last 50
	}
}

// GetAccessStats returns statistics about path usage
func (js *JSONSuggester) GetAccessStats() map[string]int {
	return js.commonPaths
}

// ResetStats clears access statistics
func (js *JSONSuggester) ResetStats() {
	js.commonPaths = make(map[string]int)
	js.history = []string{}
}

// similarity calculates string similarity (0.0 to 1.0)
func similarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	longer, shorter := s1, s2
	if len(s1) < len(s2) {
		longer, shorter = s2, s1
	}

	if len(longer) == 0 {
		return 1.0
	}

	distance := levenshteinDistance(longer, shorter)
	return (float64(len(longer)) - float64(distance)) / float64(len(longer))
}

// levenshteinDistance calculates edit distance between strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	// Initialize first row and column
	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
