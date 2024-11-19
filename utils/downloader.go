package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func extractFileName(url string) string {
	idx := strings.LastIndex(url, "/")
	return url[idx+1:]
}

func (w *WgetValues) Downloader() {

	// Create HTTP client and set User-Agent
	client := &http.Client{}

	req, err := http.NewRequest("GET", w.Url, nil)
	CheckError(err)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:132.0) Gecko/20100101 Firefox/132.0")

	// Perform the HTTP GET request
	res, err := client.Do(req)
	CheckError(err)
	defer res.Body.Close() // Ensure the body is closed after reading

	/// If Output file is not given , extract the given filename from the metadata.
	if w.OutputFile == "" {
		w.OutputFile = extractFileName(w.Url)
	}

	// Check Background Mode if given
	if !w.BackgroudMode {
		timeStarted := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("Started at: ", timeStarted)
		fmt.Printf("sending request, awaiting response... status %s\n", res.Status)
		fmt.Printf("content size: %d [~%.2fMB]\n", res.ContentLength, float64(float64(res.ContentLength)/1048576))
		fmt.Printf("saving file to: ./%s\n", w.OutputFile)
	}

	// Create the output file
	file, err := os.Create(w.OutputFile)
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
		ProgressFunction: ShowProgress,
	}

	// Read the response body
	_, err = io.Copy(file, pr)
	CheckError(err)

	// Completed downloading the file
	if !w.BackgroudMode {
		fmt.Println("\nDownload completed at:", time.Now().Format("2006-01-02 15:04:05"))
		os.Exit(0)
	}
	// return
}
