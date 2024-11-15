package utils

import (
	"regexp"
)

// ValidateURL checks if a given string is a valid URL
// It accepts URLs with or without the protocol (http://, https://, etc.).
func (wg WgetValues) ValidateURL(url string) bool {
	// Define the updated regular expression pattern for URL validation
	pattern := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$|^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the URL matches the pattern
	return re.MatchString(url)
}
