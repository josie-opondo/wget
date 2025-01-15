package appState

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ParseArgs parses the command-line arguments and returns a UrlArgs struct
func (app *AppState) ParseArgs() error {
	mirrorMode := false
	track := false

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-O=") {
			app.UrlArgs.File = arg[len("-O="):]
		} else if strings.HasPrefix(arg, "-P=") {
			app.UrlArgs.Path = arg[len("-P="):]
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			app.UrlArgs.RateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "--mirror") {
			app.UrlArgs.Mirroring = true
			mirrorMode = true
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				fmt.Println("Error: --convert-links can only be used with --mirror.")
				os.Exit(1)
			}
			app.UrlArgs.ConvertLinksFlag = true
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				fmt.Println("Error: --reject can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-R=") {
				app.UrlArgs.RejectFlag = arg[len("-R="):]
			} else {
				app.UrlArgs.RejectFlag = arg[len("--reject="):]
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				fmt.Println("Error: --exclude can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-X=") {
				app.UrlArgs.ExcludeFlag = arg[len("-X="):]
			} else {
				app.UrlArgs.ExcludeFlag = arg[len("--exclude="):]
			}
		} else if strings.HasPrefix(arg, "-B") {
			app.UrlArgs.WorkInBackground = true
		} else if strings.HasPrefix(arg, "-i=") {
			app.UrlArgs.Sourcefile = arg[len("-i="):]
			track = true
		} else if strings.HasPrefix(arg, "http") {
			app.UrlArgs.URL = arg
		} else {
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			os.Exit(1)
		}
	}
	if app.UrlArgs.RateLimit != "" {
		if strings.ToLower(string(app.UrlArgs.RateLimit[len(app.UrlArgs.RateLimit)-1])) != "k" &&
			strings.ToLower(string(app.UrlArgs.RateLimit[len(app.UrlArgs.RateLimit)-1])) != "m" {
			fmt.Println("Invalid RateLimit")
			os.Exit(1)
		}
	}
	if app.UrlArgs.WorkInBackground {
		if app.UrlArgs.Sourcefile != "" || app.UrlArgs.Path != "" {
			fmt.Println("-B flag shpuld not be used with -i or -P flags")
			os.Exit(1)
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if app.UrlArgs.Mirroring {
		if app.UrlArgs.File != "" || app.UrlArgs.Path != "" || app.UrlArgs.RateLimit != "" || app.UrlArgs.Sourcefile != "" || app.UrlArgs.WorkInBackground {
			fmt.Println("Error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed.")
			os.Exit(1)
		}
	} else {
		if app.UrlArgs.ConvertLinksFlag || app.UrlArgs.RejectFlag != "" || app.UrlArgs.ExcludeFlag != "" {
			fmt.Println("Error: --convert-links, --reject, and --exclude can only be used with --mirror.")
			os.Exit(1)
		}
	}

	// Ensure URL is provided
	if app.UrlArgs.URL == "" && !track {
		fmt.Println("Error: URL not provided.")
		os.Exit(1)

		// Validate the URL
		err := validateURL(app.UrlArgs.URL)
		if err != nil {
			fmt.Println("Error: invalid URL provided")
			os.Exit(1)
		}
	}

	return nil
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
