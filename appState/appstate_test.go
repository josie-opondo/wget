package appState

import "testing"

func (test *AppState) TestMirrorAsyncDownload(t *testing.T) {
	_, err := GetAppState()
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
	// Mock input
	outputFileName := "testfile.txt"
	urlStr := "https://example.com/testfile.txt"
	directory := "./testdir"

	// Run the function
	err = test.mirrorAsyncDownload(outputFileName, urlStr, directory)

	// Check if the error is nil (indicating success)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}
