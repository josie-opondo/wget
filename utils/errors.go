package utils

import (
	"io"
	"log"
)

// CheckError checks for errors
func CheckError(err error) {
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}