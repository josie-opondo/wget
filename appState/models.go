package appState

import (
	"sync"
)

// UrlArgs struct with exported fields (Uppercase names)
type UrlArgs struct {
	url              string
	file             string
	rateLimit        string
	path             string
	sourceFile       string
	workInBackground bool
	mirroring        bool
	rejectFlag       string
	excludeFlag      string
	convertLinksFlag bool
}

type ProcessedURLs struct {
	sync.Mutex
	urls map[string]bool
}

// AppState encapsulates global variables and synchronization primitives
type AppState struct {
	urlArgs           UrlArgs
	// rateLimitedReader RateLimitedReader
	processedURLs     ProcessedURLs
	visitedPages      map[string]bool
	visitedAssets     map[string]bool
	muPages           sync.Mutex
	muAssets          sync.Mutex
	semaphore         chan struct{}
	count             int
	tempConfigFile    string
}

func newAppstate() *AppState {
	return &AppState{
		visitedPages: make(map[string]bool),
		visitedAssets: make(map[string]bool),
		processedURLs: ProcessedURLs{
			urls: make(map[string]bool),
		},
		tempConfigFile: "progress_config.txt",
	}
}