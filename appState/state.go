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
func GetAppState() (*AppState, error) {
	var err error

	once.Do(func() {
		instance = newAppstate()
		err = instance.parseArgs()
		err = instance.taskManager(err)
	})

	if err != nil {
		return nil, err
	}
	return instance, nil
}
