package mirror

import (
	"io"
	"sync"
	"time"
)

// processedURLs is a global map to keep track of processed URLs
var processedURLs = struct {
	sync.Mutex
	urls map[string]bool
}{
	urls: make(map[string]bool),
}

// tempConfigFile is the name of the temporary file used to store the showProgress state
const tempConfigFile = "progress_config.txt"

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

// Global variables to keep track of visited URLs and synchronization
var (
	visitedPages  = make(map[string]bool)
	visitedAssets = make(map[string]bool)
	muPages       sync.Mutex
	muAssets      sync.Mutex
	semaphore         = make(chan struct{}, 50)
	count         int = 0
)

type RateLimitedReader struct {
	reader     io.Reader
	rateLimit  int64
	bucket     int64
	lastFilled time.Time
}
