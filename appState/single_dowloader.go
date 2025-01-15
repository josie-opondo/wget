package appState

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wget/utils"
)

func (app *AppState) singleDownloader(file, url, limit, directory string) {
	path := utils.ExpandPath(directory)
	fileURL := url
	startTime := time.Now()
	toDisplay, err := utils.LoadShowProgressState(app.TempConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("started at %s\n", startTime.Format("2006-01-02 15:04:05"))

	resp, err := utils.HttpRequest(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s url: [%s]\n", resp.Status, url)
		return
	}
	fmt.Printf("sending request, awaiting response... status %s\n", resp.Status)

	contentLength := resp.ContentLength
	fmt.Printf("content size: %d bytes [~%.2fMB]\n", contentLength, float64(contentLength)/1000000)

	// Set the output file name
	var outputFile string
	if file == "" {
		urlParts := strings.Split(fileURL, "/")
		file = urlParts[len(urlParts)-1]
		outputFile = filepath.Join(path, file)
	} else {
		outputFile = filepath.Join(path, file)
	}
	// Create the path if it doesn't exist
	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			fmt.Println("Error creating path:", err)
			return
		}
	}
	temp := ""
	if file != "" && directory != "" {
		fmt.Printf("saving file to: %s%s\n", directory, file)
	} else if path == "" && file != "" {
		temp = "./"
		fmt.Printf("saving file to: %s%s\n", temp, file)
	} else {
		temp = "./"
		fmt.Printf("saving file to: %s%s\n", temp, file)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	var reader io.Reader
	if limit != "" {
		reader = utils.NewRateLimitedReader(resp.Body, limit)
	} else {
		reader = resp.Body
	}

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64
	startDownload := time.Now()

	if toDisplay {
		fmt.Print("Downloading... ")
	}
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
			// Update the downloaded size
			downloaded += int64(n)

			if toDisplay {
				// Calculate and display the progress
				progress := float64(downloaded) / float64(contentLength) * 50
				speed := float64(downloaded) / time.Since(startDownload).Seconds()
				timeRemaining := time.Duration(float64(contentLength-downloaded)/speed) * time.Second

				// Update the same line with progress
				fmt.Printf("\r %.2f KiB / %.2f KiB [", float64(downloaded)/1024, float64(contentLength)/1024)
				for i := 0; i < 50; i++ {
					if i < int(progress) {
						fmt.Print("=")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Printf("] %.2f%% %.2f KiB/s %s", (float64(downloaded)*100)/float64(contentLength), speed/1024, timeRemaining.String())

			}

		}

		if downloaded >= contentLength {
			break
		}
	}
	if toDisplay {
		fmt.Println() // Move to the next line after download completes
		fmt.Println()
	}

	endTime := time.Now()
	fmt.Printf("Downloaded [%s]\n", fileURL)
	fmt.Printf("finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
	if !toDisplay {
		fmt.Println()
	}
}
