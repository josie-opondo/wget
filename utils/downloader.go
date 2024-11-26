package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func extractFileName(inputUrl string) string {
	// Parse the URL to extract its components
	parsedURL, err := url.Parse(inputUrl)
	if err != nil || parsedURL.Host == "" {
		return "default"
	}

	// Extract the domain (host)
	domain := parsedURL.Host

	// Replace invalid characters (if any, unlikely in domains)
	domain = regexp.MustCompile(`[<>:"/\\|?*]`).ReplaceAllString(domain, "_")

	// Enforce a length limit for safety
	if len(domain) > 100 {
		domain = domain[:100]
	}

	return domain
}

// expandPath expands ~ to the user's home directory if it's present in the path.
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}
		return filepath.Join(homeDir, path[1:]), nil
	}
	return path, nil
}

func pathEnding(path string) string {
	if strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}
	return path
}

// CreateIncrementalFile ensures the directory exists, creates a unique file, and returns the file pointer, its path, and an error if any.
func CreateIncrementalFile(dir, filename string) (*os.File, string, error) {
	// Resolve the directory to an absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to resolve directory path: %v", err)
	}

	// Ensure the directory exists, create it if necessary
	if err := os.MkdirAll(absDir, os.ModePerm); err != nil {
		return nil, "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Split the filename into name and extension
	base := strings.TrimSuffix(filename, filepath.Ext(filename)) // Base name without extension
	ext := filepath.Ext(filename)                                // File extension

	var fullPath string
	newFilename := filename

	// Try to create the file with incremental names
	for i := 1; ; i++ {
		fullPath = filepath.Join(absDir, newFilename)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// File does not exist, attempt to create it
			file, err := os.Create(fullPath)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create file: %v", err)
			}
			return file, fullPath, nil
		}
		// Increment the filename in the correct format
		newFilename = fmt.Sprintf("%s%s(%d)", base, ext, i)
	}
}

func (w *WgetValues) Downloader() {

	/// Output logging if Background mode is false
	var log_file string = ""

	// Create HTTP client and set User-Agent
	client := &http.Client{}

	req, err := http.NewRequest("GET", w.Url, nil)
	CheckError(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")

	// Perform the HTTP GET request
	res, err := client.Do(req)
	CheckError(err)
	defer res.Body.Close() // Ensure the body is closed after reading

	// Check if the file extension is in the reject list
	fileExtension := strings.ToLower(filepath.Ext(res.Request.URL.Path))
	for _, suffix := range w.RejectSuffixes {
		if strings.EqualFold(fileExtension, "."+suffix) {
			fmt.Printf("Rejected file with extension: %s\n", fileExtension)
			return
		}
	}

	/// If Output file is not given , extract the given filename from the metadata.
	if w.OutputFile == "" {
		w.OutputFile = extractFileName(w.Url)
	}

	// Create the file and return its address
	file, filename, err := CreateIncrementalFile(w.OutPutDirectory, w.OutputFile)

	time_stated := time.Now().Format("2006-01-02 15:04:05")
	cmd_output := fmt.Sprintf("started at: %s\nsending request, awaiting response... status %s\ncontent size: %d [~%.2fMB]\nsaving file to: %s", time_stated, res.Status, res.ContentLength, float64(float64(res.ContentLength)/1048576), (pathEnding(w.OutPutDirectory) + "/" + filename))
	// Check Background Mode if given
	if !w.BackgroudMode {
		fmt.Println(cmd_output)
	}

	if w.BackgroudMode {
		// create the file and get filename
		_, name, err := CreateIncrementalFile(".", "wget-log")
		log_file = name

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf(`Output will be written to "%s".`, log_file)
		// log_file
		if err := os.WriteFile(log_file, []byte(cmd_output), os.ModeAppend); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Create the output file
	CheckError(err)
	defer file.Close()

	// Reader
	var reader io.Reader = res.Body

	// Check if rate-limiter was passed & value is greater than 0
	if w.RateLimitValue > 0 {
		limiter := &RateLimitedReader{
			Reader: res.Body,
			Rate:   int64(w.RateLimitValue),
			Ticker: time.NewTicker(time.Second),
		}
		defer limiter.Ticker.Stop()
		reader = limiter
	}

	pr := &ProgressRecoder{
		Reader:           reader,
		Total:            res.ContentLength,
		startTime:        time.Now(),
		ProgressFunction: w.ShowProgress,
	}

	// Read the response body
	_, err = io.Copy(file, pr)
	CheckError(err)

	// completed time
	completed_at := time.Now().Format("2006-01-02 15:04:05")
	completed_str := fmt.Sprintf("\nDownload completed at: %s", completed_at)
	// Completed downloading the file
	if !w.BackgroudMode {
		fmt.Println(completed_str)
		os.Exit(0)
	} else {
		if err := os.WriteFile(log_file, []byte(completed_str), os.ModeAppend); err != nil {
			fmt.Println(err)
			return
		}
	}
	// return
}
