package mirror

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
	reader     io.Reader
	rateLimit  int64
	bucket     int64
	lastFilled time.Time
}

type ProcessedURLs struct {
	sync.Mutex
	URLs map[string]bool
}

// MirrorState encapsulates global variables and synchronization primitives
type MirrorState struct {
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

// MirrorState holds global variables and synchronization mechanisms
// Use a Singleton pattern to ensure only one instance exists
var (
	instance *MirrorState
	once     sync.Once
)

// GetMirrorState provides access to the Singleton instance of MirrorState
func GetMirrorState() *MirrorState {
	once.Do(func() {
		instance = &MirrorState{
			VisitedPages:  make(map[string]bool),
			VisitedAssets: make(map[string]bool),
			Semaphore:     make(chan struct{}, 50),
			Count:         0,
			TempConfigFile: "progress_config.txt",
		}
		instance.ProcessedURLs.URLs = make(map[string]bool)
	})
	return instance
}

// state is the global MirrorState instance
var state = GetMirrorState()
