package mirror

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// DownloadAndMirror downloads a page and its assets, recursively visiting links
func DownloadAndMirror(url, rejectTypes string, convertLink bool, pathRejects string) {
	domain, err := extractDomain(url)
	if err != nil {
		fmt.Println("Could not extract domain name for:", url, "Error:", err)
		return
	}
	// fmt.Println(url)

	state.MuPages.Lock()
	if state.VisitedPages[url] {
		state.MuPages.Unlock()
		return
	}
	state.VisitedPages[url] = true
	state.MuPages.Unlock()

	// Check if we're at the root domain and force download of index.html
	if (strings.TrimRight(url, "/") == "http://"+domain || strings.TrimRight(url, "/") == "https://"+domain) && state.Count == 0 {
		state.Count++
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
		state.Semaphore <- struct{}{}
		defer func() { <-state.Semaphore }()

		baseURL := resolveURL(url, link)
		if isRejectedPath(baseURL, pathRejects) {
			fmt.Printf("Skipping Rejected file path: %s\n", baseURL)
			return
		}
		baseURLDomain, err := extractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "Error:", err)
			return
		}

		if baseURLDomain == domain {
			if tagName == "a" {
				if strings.HasSuffix(baseURL, "/") || strings.HasSuffix(baseURL, "/index.html") {
					// Ensure index.html is downloaded first
					indexURL := strings.TrimRight(baseURL, "/") + "/index.html"
					if !state.VisitedPages[indexURL] {
						downloadAsset(indexURL, domain, rejectTypes)
						DownloadAndMirror(indexURL, rejectTypes, convertLink, pathRejects)
					}
				} else {
					DownloadAndMirror(baseURL, rejectTypes, convertLink, pathRejects)
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
				if isValidAttribute(n.Data, attr.Key) {
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
					extractAndHandleStyleURLs(attr.Val, url, domain, rejectTypes)
				}
			}
			// Check for <style> tags
			if n.Data == "style" && n.FirstChild != nil {
				extractAndHandleStyleURLs(n.FirstChild.Data, url, domain, rejectTypes)
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

// fetchAndParsePage fetches the content of the URL and parses it as HTML
func fetchAndParsePage(url string) (*html.Node, error) {
	resp, err := HttpRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

func resolveURL(base, rel string) string {
	// Remove fragment identifiers (anything starting with #)
	if fragmentIndex := strings.Index(rel, "#"); fragmentIndex != -1 {
		rel = rel[:fragmentIndex]
	}

	if strings.HasPrefix(rel, "http") {
		return rel
	}

	if strings.HasPrefix(rel, "//") {
		protocol := "http:"
		if strings.HasPrefix(base, "https") {
			protocol = "https:"
		}
		return protocol + rel
	}

	if strings.HasPrefix(rel, "/") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	if strings.HasPrefix(rel, "./") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel[1:]
	}
	if strings.HasPrefix(rel, "//") && strings.Contains(rel[2:], "/") {
		baseParts := strings.Split(base, "/")
		return baseParts[0] + "//" + baseParts[2] + rel[1:]
	}

	baseParts := strings.Split(base, "/")
	return baseParts[0] + "//" + baseParts[2] + "/" + rel
}

func downloadAsset(fileURL, domain, rejectTypes string) {
	state.MuAssets.Lock()
	if state.VisitedAssets[fileURL] {
		state.MuAssets.Unlock()
		return
	}
	state.VisitedAssets[fileURL] = true
	state.MuAssets.Unlock()

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if isRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}
	fmt.Printf("Downloading: %s\n", fileURL)
	MirrorAsyncDownload("", fileURL, "", domain)
}
