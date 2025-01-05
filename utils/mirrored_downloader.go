package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DownloadAndMirror downloads the given URL and its assets recursively.
func (w *WgetValues) DownloadAndMirror() {
	visited := make(map[string]bool) // Track visited URLs
	queue := []string{w.Url}         // Start with the root URL

	for len(queue) > 0 {
		currentURL := queue[0]
		queue = queue[1:] // Dequeue the first URL

		// Skip if already visited
		if visited[currentURL] {
			continue
		}
		visited[currentURL] = true

		// Fetch the page
		res, err := http.Get(currentURL)
		if err != nil {
			fmt.Printf("Failed to fetch %s: %v\n", currentURL, err)
			continue
		}
		defer res.Body.Close()

		// Parse HTML if content is text/html
		contentType := res.Header.Get("Content-Type")
		if strings.Contains(contentType, "text/html") {
			htmlData, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("Failed to read response body: %v\n", err)
				continue
			}

			// Save the HTML file
			basePath := extractFileName(currentURL)
			err = saveFile(basePath, "index.html", htmlData)
			if err != nil {
				fmt.Printf("Failed to save file: %v\n", err)
				continue
			}

			// Parse and extract asset links
			assets, newLinks := parseHTMLForAssets(currentURL, htmlData)
			queue = append(queue, newLinks...)

			// Download each asset
			for _, asset := range assets {
				if err := downloadAsset(asset, w.OutPutDirectory); err != nil {
					fmt.Printf("Failed to download asset %s: %v\n", asset, err)
				}
			}
		}
	}
}
