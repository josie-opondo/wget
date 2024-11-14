package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"wget/utils"
	"strings"
	"time"
)


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
	url = args[0] 

	// Perform the HTTP GET request
	res, err := http.Get(url)
	utils.CheckError(err)
	defer res.Body.Close() // Ensure the body is closed after reading

	// Determine the file name
	if flagOVal != "" {
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

	pr := &utils.ProgressRecoder{
		Reader:           res.Body,
		Total:            res.ContentLength,
		ProgressFunction: utils.ShowProgress,
	}

	file, err := os.Create(fileName)
	utils.CheckError(err)
	defer file.Close()

	// Read the response body
	_, err = io.Copy(file, pr)
	utils.CheckError(err)

	fmt.Println() // New line after progress bar
	if flagBVal {
		log.Println("Download completed at:", time.Now().Format("2006-01-02 15:04:05"))
	}
}
