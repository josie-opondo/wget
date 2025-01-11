package mirror

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func mirrorAsyncDownload(outputFileName, urlStr,  directory string) {
	app.ProcessedURLs.Lock()
	if processed, exists := app.ProcessedURLs.URLs[urlStr]; exists && processed {
		app.ProcessedURLs.Unlock()
		fmt.Printf("URL already processed: %s\n", urlStr)
		return
	}
	app.ProcessedURLs.Unlock()

	// Parse the URL to get the path components
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Create the necessary directories based on the URL path
	rootPath := expandPath(directory)
	pathComponents := strings.Split(strings.Trim(u.Path, "/"), "/")
	relativeDirPath := filepath.Join(pathComponents[:len(pathComponents)-1]...)
	fullDirPath := filepath.Join(rootPath, relativeDirPath)
	fileName := pathComponents[len(pathComponents)-1]

	resp, err := HttpRequest(urlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s url: %s\n", resp.Status, urlStr)
		return
	}

	contentType := resp.Header.Get("Content-Type")

	if outputFileName == "" {
		if fileName == "" || strings.HasSuffix(urlStr, "/") {
			fileName = "index.html"
		} else if contentType == "text/html" && !strings.HasSuffix(fileName, ".html") {
			fileName += ".html"
		}
		outputFileName = filepath.Join(fullDirPath, fileName)
	} else {
		if contentType == "text/html" && !strings.HasSuffix(outputFileName, ".html") {
			outputFileName += ".html"
		}
		outputFileName = filepath.Join(fullDirPath, outputFileName)
	}

	if fullDirPath != "" {
		if _, err := os.Stat(fullDirPath); os.IsNotExist(err) {
			err = os.MkdirAll(fullDirPath, 0o755)
			if err != nil {
				fmt.Println("Error creating path:", err)
				return
			}
		}
	}
	if fileExists(outputFileName) {
		return
	}

	var out *os.File
	out, err = os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer out.Close()

	var reader io.Reader = resp.Body

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading response body:", err)
			return
		}

		if n > 0 {
			if _, err := out.Write(buffer[:n]); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			downloaded += int64(n)
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("\033[32mDownloaded [%s]\033[0m\n", urlStr)

	// Mark the URL as processed
	app.ProcessedURLs.Lock()
	app.ProcessedURLs.URLs[urlStr] = true
	app.ProcessedURLs.Unlock()
}
