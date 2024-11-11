package utils

import (
	"io"
	"log"
)


func CheckError(err error) {
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}