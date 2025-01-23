package appState

import (
	"fmt"
	"os"
	"strings"

	"wget/utils"
)

// taskManager calls to action methods depending on the passed flags
func (app *AppState) taskManager(err error) error {
	if err != nil {
		return err
	}

	// Mirror website handling
	if app.urlArgs.mirroring {
		err := app.downloadAndMirror(app.urlArgs.url, app.urlArgs.rejectFlag, app.urlArgs.convertLinksFlag, app.urlArgs.excludeFlag)
		if err != nil {
			return err
		}
		return nil
	}

	// If no file name is provided, derive it from the url
	if app.urlArgs.file == "" && app.urlArgs.url != "" {
		urlParts := strings.Split(app.urlArgs.url, "/")
		app.urlArgs.file = urlParts[len(urlParts)-1]
	}

	// Handle the work-in-background flag
	if app.urlArgs.workInBackground {
		err := app.downloadInBackground(app.urlArgs.file, app.urlArgs.url, app.urlArgs.rateLimit)
		if err != nil {
			return err
		}
		return nil
	}

	// Handle multiple file downloads from sourceFile
	if app.urlArgs.sourceFile != "" {
		err := app.downloadMultipleFiles(app.urlArgs.sourceFile, app.urlArgs.file, app.urlArgs.rateLimit, app.urlArgs.path)
		if err != nil {
			return err
		}
		return nil
	}

	// Ensure url is provided
	if app.urlArgs.url == "" {
		return fmt.Errorf("error: url not provided")
	}

	// Start downloading the file
	err = app.singleDownloader(app.urlArgs.file, app.urlArgs.url, app.urlArgs.rateLimit, app.urlArgs.path)
	if err != nil {
		return err
	}
	return nil
}

// ParseArgs parses the command-line arguments and returns a urlArgs struct
func (app *AppState) parseArgs() error {
	mirrorMode := false
	track := false

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-O=") {
			app.urlArgs.file = arg[len("-O="):]
		} else if strings.HasPrefix(arg, "-P=") {
			app.urlArgs.path = arg[len("-P="):]
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			if err := utils.RateLimitValidator(arg); err != nil {
				return err
			}

			app.urlArgs.rateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "--mirror") {
			app.urlArgs.mirroring = true
			mirrorMode = true
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				return fmt.Errorf("error: --convert-links can only be used with --mirror")
			}
			app.urlArgs.convertLinksFlag = true
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				return fmt.Errorf("error: --reject can only be used with --mirror")
			}
			if strings.HasPrefix(arg, "-R=") {
				app.urlArgs.rejectFlag = arg[len("-R="):]
			} else {
				app.urlArgs.rejectFlag = arg[len("--reject="):]
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				return fmt.Errorf("error: --exclude can only be used with --mirror")
			}
			if strings.HasPrefix(arg, "-X=") {
				app.urlArgs.excludeFlag = arg[len("-X="):]
			} else {
				app.urlArgs.excludeFlag = arg[len("--exclude="):]
			}
		} else if strings.HasPrefix(arg, "-B") {
			app.urlArgs.workInBackground = true
		} else if strings.HasPrefix(arg, "-i=") {
			app.urlArgs.sourceFile = arg[len("-i="):]
			track = true
		} else if strings.HasPrefix(arg, "http") {
			app.urlArgs.url = arg
		} else {
			return fmt.Errorf("error: Unrecognized argument'%s'", arg)
		}
	}

	if app.urlArgs.rateLimit != "" {
		if strings.ToLower(string(app.urlArgs.rateLimit[len(app.urlArgs.rateLimit)-1])) != "k" &&
			strings.ToLower(string(app.urlArgs.rateLimit[len(app.urlArgs.rateLimit)-1])) != "m" {
			return fmt.Errorf("invalid rateLimit")
		}
	}

	if app.urlArgs.workInBackground {
		if app.urlArgs.sourceFile != "" || app.urlArgs.path != "" {
			return fmt.Errorf("-B flag shpuld not be used with -i or -P flags")
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if app.urlArgs.mirroring {
		if app.urlArgs.file != "" || app.urlArgs.path != "" || app.urlArgs.rateLimit != "" || app.urlArgs.sourceFile != "" || app.urlArgs.workInBackground {
			return fmt.Errorf("error: --mirror can only be used with --convert-links, --reject, --exclude, and a url. No other flags are allowed")
		}
	} else {
		if app.urlArgs.convertLinksFlag || app.urlArgs.rejectFlag != "" || app.urlArgs.excludeFlag != "" {
			return fmt.Errorf("error: --convert-links, --reject, and --exclude can only be used with --mirror")
		}
	}

	// Ensure url is provided
	if app.urlArgs.url == "" && !track {
		return fmt.Errorf("error: url not provided")
	}

	// Validate the url
	err := utils.Validateurl(app.urlArgs.url)
	if err != nil {
		return fmt.Errorf("error: invalid url provided")
	}

	return nil
}
