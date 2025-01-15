package appState

import (
	"testing"
)

func TestMirrorAsyncDownload(t *testing.T) {
	app, err := GetAppState()
	if err != nil {
		return
	}
	// Mock input
	outputFileName := "testfile.txt"
	urlStr := "https://example.com/testfile.txt"
	directory := "./testdir"

	// Run the function
	err = app.mirrorAsyncDownload(outputFileName, urlStr, directory)

	// Check if the error is nil (indicating success)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestDownloadInBackground(t *testing.T) {
	// Initialize AppState
	app, err := GetAppState()
	if err != nil {
		return
	}

	// Mock input
	file := "testfile.txt"
	urlStr := "https://example.com/testfile.txt"
	rateLimit := "500k" // Example rate limit

	// Run the function
	err = app.DownloadInBackground(file, urlStr, rateLimit)

	// Check if the error is nil (indicating success)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestDownloadAndMirror(t *testing.T) {
	// Setup app state
	app, err := GetAppState()
	if err != nil {
		return
	}

	// Mock input
	url := "https://example.com"
	rejectTypes := "image/png,image/jpg"
	convertLink := true
	pathRejects := "/ignore"

	// Run the function
	err = app.DownloadAndMirror(url, rejectTypes, convertLink, pathRejects)

	// Check if the error is nil (indicating success)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}
