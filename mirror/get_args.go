package mirror

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ParseArgs parses the command-line arguments and returns a UrlArgs struct
func ParseArgs() UrlArgs {
	args := &UrlArgs{}
	mirrorMode := false // Flag to track if --mirror is set
	track := false

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		// Enforce flags with the '=' sign
		if strings.HasPrefix(arg, "-O=") {
			args.File = arg[len("-O="):] // Capture the file name
		} else if strings.HasPrefix(arg, "-P=") {
			args.Path = arg[len("-P="):] // Capture the path
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			args.RateLimit = arg[len("--rate-limit="):] // Capture the rate limit
		} else if strings.HasPrefix(arg, "--mirror") {
			args.Mirroring = true // Enable mirroring
			mirrorMode = true      // Track mirror mode
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				fmt.Println("Error: --convert-links can only be used with --mirror.")
				os.Exit(1)
			}
			args.ConvertLinksFlag = true // Enable link conversion
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				fmt.Println("Error: --reject can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-R=") {
				args.RejectFlag = arg[len("-R="):] // Capture reject flag for -R
			} else {
				args.RejectFlag = arg[len("--reject="):] // Capture reject flag for --reject
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				fmt.Println("Error: --exclude can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-X=") {
				args.ExcludeFlag = arg[len("-X="):] // Capture exclude flag for -X
			} else {
				args.ExcludeFlag = arg[len("--exclude="):] // Capture exclude flag for --exclude
			}
		} else if strings.HasPrefix(arg, "-B") {
			args.WorkInBackground = true // Enable background downloading
		} else if strings.HasPrefix(arg, "-i=") {
			args.Sourcefile = arg[len("-i="):] // Capture source file
			track = true
		} else if strings.HasPrefix(arg, "http") {
			// This must be the URL
			args.URL = arg
		} else {
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			os.Exit(1)
		}
	}
	if args.RateLimit != "" {
		if strings.ToLower(string(args.RateLimit[len(args.RateLimit)-1])) != "k" &&
			strings.ToLower(string(args.RateLimit[len(args.RateLimit)-1])) != "m" {
			fmt.Println("Invalid RateLimit")
			os.Exit(1)
		}
	}
	if args.WorkInBackground {
		if args.Sourcefile != "" || args.Path != "" {
			fmt.Println("-B flag shpuld not be used with -i or -P flags")
			os.Exit(1)
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if args.Mirroring {
		// Only allow --convert-links, --reject, and --exclude with --mirror
		if args.File != "" || args.Path != "" || args.RateLimit != "" || args.Sourcefile != "" || args.WorkInBackground {
			fmt.Println("Error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed.")
			os.Exit(1)
		}
	} else {
		// If --mirror is not provided, reject the use of --convert-links, --reject, and --exclude
		if args.ConvertLinksFlag || args.RejectFlag != "" || args.ExcludeFlag != "" {
			fmt.Println("Error: --convert-links, --reject, and --exclude can only be used with --mirror.")
			os.Exit(1)
		}
	}

	// Ensure URL is provided
	if args.URL == "" && !track {
		fmt.Println("Error: URL not provided.")
		os.Exit(1)

		// Validate the URL
		err := validateURL(args.URL)
		if err != nil {
			fmt.Println("Error: invalid URL provided")
			os.Exit(1)
		}
	}

	return *args
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
