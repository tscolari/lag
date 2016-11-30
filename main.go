package main

import (
	"fmt"
	"os"

	"github.com/tscolari/lag/parser"
	"github.com/tscolari/lag/ui"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: lag <LOG FILE>")
		os.Exit(1)
	}

	logFilePath := os.Args[1]
	logFile, err := os.Open(logFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %s\n", err.Error())
		os.Exit(1)
	}

	entries, _ := parser.Parse(logFile)

	ui := ui.New(entries)
	err = ui.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start UI: %s\n", err.Error())
		os.Exit(1)
	}
}
