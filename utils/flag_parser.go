package utils

import (
	"fmt"
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

func (w *WgetValues) FlagsParser(args []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-B":
			w.BackgroudMode = true
		case "-O":
			if i+1 < len(args) {
				w.OutputFile = args[i+1]
				i++ // Skip the next element as it's the value for -O
			}
		case "-P":
			if i+1 < len(args) {
				w.OutPutDirectory = args[i+1]
				i++
			}
		case "--rate-limit":
			if i+1 < len(args) {
				w.RateLimitValue = args[i+1]
				i++
			}
		case "--reject":
			w.Reject = true
		case "--exclude", "-X":
			if i+1 < len(args) {
				w.Exclude = args[i+1]
				i++
			}
		case "--convert-links":
			w.ConvertLinks = true
		case "--mirror":
			w.Mirror = true
		default:
			fmt.Printf("Unknown argument: %s\n", arg)
			return
		}
	}
}
