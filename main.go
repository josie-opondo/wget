package main

import (
	"fmt"
	"os"
	"wget/utils"
)

func main() {
	// Parse the commandline arguments
	cmd_args := os.Args[1:]

	//
	w := utils.WgetInstance()
	err := w.FlagsParser(cmd_args)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Downloader()

}
