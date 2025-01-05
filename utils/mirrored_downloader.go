package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
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

// parseHTMLForAssets parses HTML content and extracts asset URLs and links.
func parseHTMLForAssets(baseURL string, htmlData []byte) (assets []string, links []string) {
	doc, err := html.Parse(strings.NewReader(string(htmlData)))
	if err != nil {
		fmt.Printf("Failed to parse HTML: %v\n", err)
		return nil, nil
	}

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var link string
			switch n.Data {
			case "link", "img", "script":
				for _, attr := range n.Attr {
					if attr.Key == "href" || attr.Key == "src" {
						link = attr.Val
					}
				}
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link = attr.Val
					}
				}
			}

			// Normalize the link
			if link != "" {
				fullURL := normalizeURL(baseURL, link)
				if strings.Contains(fullURL, baseURL) {
					links = append(links, fullURL) // Add internal links for recursion
				} else {
					assets = append(assets, fullURL) // Treat external links as assets
				}
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return assets, links
}

// normalizeURL resolves relative URLs against a base URL.
func normalizeURL(baseURL, relative string) string {
	u, err := url.Parse(relative)
	if err != nil {
		return ""
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	return base.ResolveReference(u).String()
}

// downloadAsset downloads a single asset and saves it to the output directory.
func downloadAsset(assetURL, outputDir string) error {
	res, err := http.Get(assetURL)
	if err != nil {
		return fmt.Errorf("failed to fetch asset: %v", err)
	}
	defer res.Body.Close()

	// Create directories based on URL path
	parsedURL, err := url.Parse(assetURL)
	if err != nil {
		return fmt.Errorf("failed to parse asset URL: %v", err)
	}
	assetPath := filepath.Join(outputDir, parsedURL.Path)
	err = os.MkdirAll(filepath.Dir(assetPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	// Save the file
	file, err := os.Create(assetPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

// saveFile saves the given data to a specified directory and filename.
func saveFile(basePath, fileName string, data []byte) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the full file path
	filePath := filepath.Join(basePath, fileName)

	// Write data to the file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
