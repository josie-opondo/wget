package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type progressRecoder struct {
	reader           io.Reader
	total            int64
	progress         int64
	progressFunction func(int64, int64)
}

func (pr *progressRecoder) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	checkError(err)
	pr.progress += int64(n)
	if pr.progressFunction != nil {
		pr.progressFunction(pr.progress, pr.total)
	}
	return
}

func showProgress(progress int64, total int64) {
	const length = 50
	if total <= 0 {
		fmt.Printf("\rDownloading...")
		return
	}
	percent := float64(progress) / float64(total) * 100
	numBars := int((percent / 100) * length)
	out := fmt.Sprintf("%.2f KiB / %.2f KiB [%s%s] %.0f%%",float64(progress)/1024, float64(total)/1024, strings.Repeat("=", numBars), strings.Repeat(" ", length-numBars), percent)
	fmt.Printf("\r%s", out)
}

func checkError(err error) {
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}

func main() {
	var (
		flagOVal string
		flagBVal bool
		url      string
		fileName string
	)

	// Define flags
	flag.StringVar(&flagOVal, "O", "", "Specify the output file name")
	flag.BoolVar(&flagBVal, "B", false, "Enable verbose logging")

	// Parse flags
	flag.Parse()

	// Get the remaining arguments after flags
	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("URL is required")
	}
	url = args[0] // Assume the first argument after flags is the URL

	// Perform the HTTP GET request
	res, err := http.Get(url)
	checkError(err)
	defer res.Body.Close() // Ensure the body is closed after reading

	// Determine the file name
	if flagOVal != "" {
		// Use the `-O` flag value for the file name, preserving the extension if possible
		prefIndex := strings.LastIndex(url, ".")
		if prefIndex != -1 && len(url)-prefIndex <= 4 && !strings.Contains(flagOVal, "."){ // Check if the extension is likely valid
			fileName = flagOVal + url[prefIndex:]
		} else {
			fileName = flagOVal
		}
	} else {
		// Extract file name from URL if `-O` is not provided
		parts := strings.Split(url, "/")
		fileName = parts[len(parts)-1]
	}

	// Handle the `-B` flag for logging time started
	if flagBVal {
		log.Println("Download started at:")
	}
	timeStarted := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Started at: ", timeStarted)
	fmt.Printf("sending request, awaiting response... status %s\n", res.Status)
	fmt.Printf("content size: %d [~%.2fMB]\n", res.ContentLength, float64(float64(res.ContentLength)/1048576))
	fmt.Printf("saving file to: ./%s\n", fileName)

	pr := &progressRecoder{
		reader:           res.Body,
		total:            res.ContentLength,
		progressFunction: showProgress,
	}

	file, err := os.Create(fileName)
	checkError(err)
	defer file.Close()

	// Read the response body
	_, err = io.Copy(file, pr)
	checkError(err)

	fmt.Println() // New line after progress bar
	if flagBVal {
		log.Println("Download completed at:", time.Now().Format("2006-01-02 15:04:05"))
	}
}
