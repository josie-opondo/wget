package mirror

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ParseArgs parses the command-line arguments and returns a UrlArgs struct
func ParseArgs() (*MirrorState, error) {
	mirrorMode := false
	track := false

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-O=") {
			state.UrlArgs.File = arg[len("-O="):]
		} else if strings.HasPrefix(arg, "-P=") {
			state.UrlArgs.Path = arg[len("-P="):]
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			state.UrlArgs.RateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "--mirror") {
			state.UrlArgs.Mirroring = true
			mirrorMode = true
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				fmt.Println("Error: --convert-links can only be used with --mirror.")
				os.Exit(1)
			}
			state.UrlArgs.ConvertLinksFlag = true
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				fmt.Println("Error: --reject can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-R=") {
				state.UrlArgs.RejectFlag = arg[len("-R="):]
			} else {
				state.UrlArgs.RejectFlag = arg[len("--reject="):]
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				fmt.Println("Error: --exclude can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-X=") {
				state.UrlArgs.ExcludeFlag = arg[len("-X="):]
			} else {
				state.UrlArgs.ExcludeFlag = arg[len("--exclude="):]
			}
		} else if strings.HasPrefix(arg, "-B") {
			state.UrlArgs.WorkInBackground = true
		} else if strings.HasPrefix(arg, "-i=") {
			state.UrlArgs.Sourcefile = arg[len("-i="):]
			track = true
		} else if strings.HasPrefix(arg, "http") {
			state.UrlArgs.URL = arg
		} else {
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			os.Exit(1)
		}
	}
	if state.UrlArgs.RateLimit != "" {
		if strings.ToLower(string(state.UrlArgs.RateLimit[len(state.UrlArgs.RateLimit)-1])) != "k" &&
			strings.ToLower(string(state.UrlArgs.RateLimit[len(state.UrlArgs.RateLimit)-1])) != "m" {
			fmt.Println("Invalid RateLimit")
			os.Exit(1)
		}
	}
	if state.UrlArgs.WorkInBackground {
		if state.UrlArgs.Sourcefile != "" || state.UrlArgs.Path != "" {
			fmt.Println("-B flag shpuld not be used with -i or -P flags")
			os.Exit(1)
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if state.UrlArgs.Mirroring {
		if state.UrlArgs.File != "" || state.UrlArgs.Path != "" || state.UrlArgs.RateLimit != "" || state.UrlArgs.Sourcefile != "" || state.UrlArgs.WorkInBackground {
			fmt.Println("Error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed.")
			os.Exit(1)
		}
	} else {
		if state.UrlArgs.ConvertLinksFlag || state.UrlArgs.RejectFlag != "" || state.UrlArgs.ExcludeFlag != "" {
			fmt.Println("Error: --convert-links, --reject, and --exclude can only be used with --mirror.")
			os.Exit(1)
		}
	}

	// Ensure URL is provided
	if state.UrlArgs.URL == "" && !track {
		fmt.Println("Error: URL not provided.")
		os.Exit(1)

		// Validate the URL
		err := validateURL(state.UrlArgs.URL)
		if err != nil {
			fmt.Println("Error: invalid URL provided")
			os.Exit(1)
		}
	}

	return state, nil
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
