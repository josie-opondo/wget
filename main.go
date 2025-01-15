package main

import (
	"fmt"
	"os"
	"wget/appState"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <URL> [options]")
		return
	}

	_, err := appState.GetAppState()
	if err != nil {
		fmt.Println(err)
		return
	}
}
