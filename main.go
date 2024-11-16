package main

import (
	"os"
	"wget/utils"
)

func main() {
	// Parse the commandline arguments
	cmd_args := os.Args[1:]

	//
	w := utils.WgetInstance()
	w.FlagsParser(cmd_args)
	w.Downloader()

}
