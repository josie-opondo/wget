package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ValidateURL checks if a given string is a valid URL
// It accepts URLs with or without the protocol (http://, https://, etc.).
func ValidateURL(url string) bool {
	// Define the updated regular expression pattern for URL validation
	pattern := `^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$|^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+$`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Check if the URL matches the pattern
	return re.MatchString(url)
}

func RateLimitValue(s string) int {
	if !strings.Contains(s, "k") && !strings.Contains(s, "M") {
		fmt.Println("Invalid rate limit value.\nUsage: --rate-limit=400k || --rate-limit=2M")
		os.Exit(0)
	}

	if strings.Contains(s, "k") {
		ln := len(s) - 1
		// int value string
		val := s[:ln]
		// convert the value to int
		num, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Invalid rate limit value here")
			os.Exit(0)
		}
		return num * 1024
	}

	if strings.Contains(s, "M") {
		ln := len(s) - 1
		// int value string
		val := s[:ln]
		// convert the value to int
		num, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Invalid rate limit value there")
			os.Exit(0)
		}
		return num * 1024 * 1024
	}
	return 0
}

func (w *WgetValues) FlagsParser(args []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if ValidateURL(arg) {
			w.Url = arg
		} else if strings.Contains(arg, "--rate-limit=") {
			idx := strings.Index(arg, "=")
			value := strings.Trim(arg[idx+1:], " ")
			if value == "" {
				fmt.Println("Usage: --rate-limit=400k || --rate-limit=1M")
				os.Exit(0)
			}
			w.RateLimitValue = RateLimitValue(value)
		} else if strings.Contains(arg, "-P") {
			idx := strings.Index(arg, "=")
			value := strings.Trim(arg[idx+1:], " ")

			if value == "" {
				fmt.Println("Invalid download path")
				os.Exit(0)
			}
			w.OutPutDirectory = value
		} else if strings.Contains(arg, "--reject=") || strings.Contains(arg, "-R=") {
			// Parse the --reject flag
			var rejectValue string
			if strings.Contains(arg, "--reject=") {
				rejectValue = strings.TrimPrefix(arg, "--reject=")
			} else {
				rejectValue = strings.TrimPrefix(arg, "-R=")
			}
			w.RejectSuffixes = strings.Split(rejectValue, ",")
		} else {
			switch arg {
			case "-B":
				w.BackgroudMode = true
			case "-O":
				if i+1 < len(args) {
					w.OutputFile = args[i+1]
					i++ // Skip the next element as it's the value for -O
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
}
