package appState

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"wget/utils"
)

func (app *AppState) DownloadInBackground(file, urlStr, rateLimit string) {
	// Parse the URL to derive the output name
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Invalid URL:", err)
		return
	}
	outputName := filepath.Base(parsedURL.Path) // Get the file name from the URL
	if file != "" {
		outputName = file
	}

	path := "." // Default path to save the file
	// Create the wget-log file to log output
	logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer logFile.Close()

	// Ensure the output directory exists
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}
	cmd := exec.Command(os.Args[0], "-O="+outputName, "-P="+path, "--rate-limit="+rateLimit, urlStr)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	fmt.Println("Output will be written to \"wget-log\".")

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting download:", err)
		return
	}
	if err := utils.SaveShowProgressState(app.TempConfigFile, false); err != nil {
		fmt.Println(err)
		return
	}

	// Wait for the command to complete in the background
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("Error during download:", err)
			return
		}
	}()
}
