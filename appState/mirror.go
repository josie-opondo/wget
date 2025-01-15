package appState

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"wget/utils"

	"golang.org/x/net/html"
)

// DownloadAndMirror downloads a page and its assets, recursively visiting links
func (app *AppState) DownloadAndMirror(url, rejectTypes string, convertLink bool, pathRejects string) {
	domain, err := utils.ExtractDomain(url)
	if err != nil {
		fmt.Println("Could not extract domain name for:", url, "Error:", err)
		return
	}

	app.MuPages.Lock()
	if app.VisitedPages[url] {
		app.MuPages.Unlock()
		return
	}
	app.VisitedPages[url] = true
	app.MuPages.Unlock()

	// Check if we're at the root domain and force download of index.html
	if (strings.TrimRight(url, "/") == "http://"+domain || strings.TrimRight(url, "/") == "https://"+domain) && app.Count == 0 {
		app.Count++
		indexURL := strings.TrimRight(url, "/")
		downloadAsset(indexURL, domain, rejectTypes)
	}

	// Fetch and get the HTML of the page
	doc, err := fetchAndParsePage(url)
	if err != nil {
		fmt.Println("Error fetching or parsing page:", err)
		return
	}

	// Function to handle links and assets found on the page
	handleLink := func(link, tagName string) {
		app.Semaphore <- struct{}{}
		defer func() { <-app.Semaphore }()

		baseURL := utils.ResolveURL(url, link)
		if utils.IsRejectedPath(baseURL, pathRejects) {
			fmt.Printf("Skipping Rejected file path: %s\n", baseURL)
			return
		}
		baseURLDomain, err := utils.ExtractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "Error:", err)
			return
		}

		if baseURLDomain == domain {
			if tagName == "a" {
				if strings.HasSuffix(baseURL, "/") || strings.HasSuffix(baseURL, "/index.html") {
					// Ensure index.html is downloaded first
					indexURL := strings.TrimRight(baseURL, "/") + "/index.html"
					if !app.VisitedPages[indexURL] {
						downloadAsset(indexURL, domain, rejectTypes)
						app.DownloadAndMirror(indexURL, rejectTypes, convertLink, pathRejects)
					}
				} else {
					app.DownloadAndMirror(baseURL, rejectTypes, convertLink, pathRejects)
				}
			}
			downloadAsset(baseURL, domain, rejectTypes)
		}
	}

	var wg sync.WaitGroup
	var processNode func(n *html.Node)

	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if utils.IsValidAttribute(n.Data, attr.Key) {
					link := attr.Val
					if link != "" {
						wg.Add(1)
						go func(link, tagName string) {
							defer wg.Done()
							handleLink(link, tagName)
						}(link, n.Data)
					}
				}
				// Check for inline styles
				if attr.Key == "style" {
					app.extractAndHandleStyleURLs(attr.Val, url, domain, rejectTypes)
				}
			}
			// Check for <style> tags
			if n.Data == "style" && n.FirstChild != nil {
				app.extractAndHandleStyleURLs(n.FirstChild.Data, url, domain, rejectTypes)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	// Start processing the document
	processNode(doc)

	// Wait for all goroutines to complete
	wg.Wait()

	// Convert links if the flag is set
	if convertLink {
		convertLinks(url)
	}
}

func (app *AppState) extractAndHandleStyleURLs(styleContent, baseURL, domain, rejectTypes string) {
	re := regexp.MustCompile(`url\(['"]?([^'"()]+)['"]?\)`)
	matches := re.FindAllStringSubmatch(styleContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			assetURL := utils.ResolveURL(baseURL, match[1])
			downloadAsset(assetURL, domain, rejectTypes)
		}
	}
}

// fetchAndParsePage fetches the content of the URL and parses it as HTML
func fetchAndParsePage(url string) (*html.Node, error) {
	resp, err := utils.HttpRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

func downloadAsset(fileURL, domain, rejectTypes string) {
	app.MuAssets.Lock()
	if app.VisitedAssets[fileURL] {
		app.MuAssets.Unlock()
		return
	}
	app.VisitedAssets[fileURL] = true
	app.MuAssets.Unlock()

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if isRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}
	fmt.Printf("Downloading: %s\n", fileURL)
	mirrorAsyncDownload("", fileURL, domain)
}
