package appState

import (
	"io"
	"sync"
	"time"
)

// UrlArgs struct with exported fields (Uppercase names)
type UrlArgs struct {
	URL              string
	File             string
	RateLimit        string
	Path             string
	Sourcefile       string
	WorkInBackground bool
	Mirroring        bool
	RejectFlag       string
	ExcludeFlag      string
	ConvertLinksFlag bool
}

type RateLimitedReader struct {
	Reader     io.Reader
	RateLimit  int64
	Bucket     int64
	LastFilled time.Time
}

type ProcessedURLs struct {
	sync.Mutex
	URLs map[string]bool
}

// AppState encapsulates global variables and synchronization primitives
type AppState struct {
	UrlArgs       UrlArgs
	RateLimitedReader RateLimitedReader
	ProcessedURLs ProcessedURLs
	VisitedPages  map[string]bool
	VisitedAssets map[string]bool
	MuPages       sync.Mutex
	MuAssets      sync.Mutex
	Semaphore     chan struct{}
	Count         int
	TempConfigFile string
}
