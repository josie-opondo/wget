package appState

import (
	"sync"
)

// AppState holds global variables and synchronization mechanisms
// Use a Singleton pattern to ensure only one instance exists
var (
	instance *AppState
	once     sync.Once
)

// GetAppState provides access to the Singleton instance of AppState
func GetAppState() *AppState {
	once.Do(func() {
		instance = &AppState{
			VisitedPages:   make(map[string]bool),
			VisitedAssets:  make(map[string]bool),
			Semaphore:      make(chan struct{}, 50),
			Count:          0,
			TempConfigFile: "progress_config.txt",
		}
		instance.ProcessedURLs.URLs = make(map[string]bool)
		instance.ParseArgs()
	})
	return instance
}
