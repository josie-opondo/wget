package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func isRejected(url, rejectTypes string) bool {
	if rejectTypes == "" {
		return false
	}

	rejectedTypes := strings.Split(rejectTypes, ",")
	for _, ext := range rejectedTypes {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}

func isRejectedPath(url, pathRejects string) bool {
	if pathRejects == "" {
		return false
	}

	rejects := strings.Split(pathRejects, ",")
	for _, path := range rejects {
		if path[0] != '/' {
			continue
		}
		if contains(url, path) {
			return true
		}
	}

	return false
}

func contains(str, substr string) bool {
	for i := 0; i < len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// func extractAndHandleStyleURLs(styleContent, baseURL, domain, rejectTypes string) {
// 	re := regexp.MustCompile(`url\(['"]?([^'"()]+)['"]?\)`)
// 	matches := re.FindAllStringSubmatch(styleContent, -1)
// 	for _, match := range matches {
// 		if len(match) > 1 {
// 			assetURL := resolveURL(baseURL, match[1])
// 			downloadAsset(assetURL, domain, rejectTypes)
// 		}
// 	}
// }

func extractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// isValidAttribute checks if an HTML tag attribute is valid for processing
func isValidAttribute(tagName, attrKey string) bool {
	return (tagName == "a" && attrKey == "href") ||
		(tagName == "img" && attrKey == "src") ||
		(tagName == "script" && attrKey == "src") ||
		(tagName == "link" && attrKey == "href")
}

func HttpRequest(url string) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request with a User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers to mimic a Chrome browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	return resp, err
}

// expandPath expands shorthand notations to full paths
func expandPath(path string) string {
	// 1. Expand `~` to the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			return ""
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// 2. Expand environment variables like $HOME, $USER, etc.
	path = os.ExpandEnv(path)

	// 3. Convert relative paths (./ or ../) to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	return absPath
}
