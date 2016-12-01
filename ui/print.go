package ui

import (
	"fmt"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

func printEntries(v *gocui.View, entries parser.Entries) {
	for _, entry := range entries {
		printEntryHeader(v, entry)
	}
}

func printEntryHeader(v *gocui.View, entry *parser.Entry) {
	date := entry.Data.Timestamp.Format(time.RFC3339)
	message := entry.Data.Message
	if entry.Errored {
		message = redText(entry.Data.Message)
	}

	fmt.Fprintf(v, "%s [%s] (%d) %s \n",
		blueText(date),
		yellowText(entry.Data.Session),
		len(entry.Children),
		message)
}

func printEntryInfo(v *gocui.View, entry *parser.Entry) {
	fmt.Fprintln(v, " ")
	fmt.Fprintf(v, "  %s\n", logLevelToString(entry.Data.LogLevel))
	fmt.Fprintf(v, "  %s: %s\n", yellowText("Message"), entry.Data.Message)
	fmt.Fprintf(v, "  %s: %s\n", yellowText("Time"), blueText(entry.Data.Timestamp.Format(time.RFC3339)))
	if entry.Data.Error != nil {
		fmt.Fprintf(v, "  %s: %s\n", yellowText("Error"), blueText(entry.Data.Error.Error()))
	}
	fmt.Fprintf(v, "  %s\n", yellowText("--------------------------------------------------------"))

	for key, value := range entry.Data.Data {
		fmt.Fprintf(v, "  %s: %v\n", yellowText(key), value)
	}
}

func logLevelToString(logLevel lager.LogLevel) string {
	switch logLevel {
	case lager.DEBUG:
		return "DEBUG"
	case lager.INFO:
		return "INFO"
	case lager.FATAL:
		return redText("FATAL")
	case lager.ERROR:
		return redText("ERROR")
	default:
		return "Unknown"
	}
}
