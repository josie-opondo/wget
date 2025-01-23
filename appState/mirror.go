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
func (app *AppState) downloadAndMirror(url, rejectTypes string, convertLink bool, pathRejects string) error {
	domain, err := utils.ExtractDomain(url)
	if err != nil {
		return fmt.Errorf("could not extract domain name for:\n%serror: %v", url, err)
	}

	app.muPages.Lock()
	if app.visitedPages[url] {
		app.muPages.Unlock()
		return nil
	}
	app.visitedPages[url] = true
	app.muPages.Unlock()

	// Check if we're at the root domain and force download of index.html
	if (strings.TrimRight(url, "/") == "http://"+domain || strings.TrimRight(url, "/") == "https://"+domain) && app.count == 0 {
		app.count++
		indexURL := strings.TrimRight(url, "/")
		app.downloadAsset(indexURL, domain, rejectTypes)
	}

	// Fetch and get the HTML of the page
	doc, err := fetchAndParsePage(url)
	if err != nil {
		return fmt.Errorf("error fetching or parsing page:\n%v", err)
	}

	// Function to handle links and assets found on the page
	handleLink := func(link, tagName string) {
		app.semaphore <- struct{}{}
		defer func() { <-app.semaphore }()

		baseURL := utils.ResolveURL(url, link)
		if utils.IsRejectedPath(baseURL, pathRejects) {
			fmt.Printf("Skipping Rejected file path: %s\n", baseURL)
			return
		}
		baseURLDomain, err := utils.ExtractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "\nError:", err)
			return
		}

		if baseURLDomain == domain {
			if tagName == "a" {
				if strings.HasSuffix(baseURL, "/") || strings.HasSuffix(baseURL, "/index.html") {
					// Ensure index.html is downloaded first
					indexURL := strings.TrimRight(baseURL, "/") + "/index.html"
					if !app.visitedPages[indexURL] {
						app.downloadAsset(indexURL, domain, rejectTypes)
						app.downloadAndMirror(indexURL, rejectTypes, convertLink, pathRejects)
					}
				} else {
					app.downloadAndMirror(baseURL, rejectTypes, convertLink, pathRejects)
				}
			}
			app.downloadAsset(baseURL, domain, rejectTypes)
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
		utils.ConvertLinks(url)
	}

	return nil
}

func (app *AppState) extractAndHandleStyleURLs(styleContent, baseURL, domain, rejectTypes string) {
	re := regexp.MustCompile(`url\(['"]?([^'"()]+)['"]?\)`)
	matches := re.FindAllStringSubmatch(styleContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			assetURL := utils.ResolveURL(baseURL, match[1])
			app.downloadAsset(assetURL, domain, rejectTypes)
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

func (app *AppState) downloadAsset(fileURL, domain, rejectTypes string) {
	app.muAssets.Lock()
	if app.visitedAssets[fileURL] {
		app.muAssets.Unlock()
		return
	}
	app.visitedAssets[fileURL] = true
	app.muAssets.Unlock()

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if utils.IsRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}
	fmt.Printf("Downloading: %s\n", fileURL)
	app.mirrorAsyncDownload("", fileURL, domain)
}
