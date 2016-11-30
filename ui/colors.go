package ui

import "fmt"

func redText(text string) string {
	return fmt.Sprintf("\x1b[0;31m%s\x1b", text)
}

func blueText(text string) string {
	return fmt.Sprintf("\x1b[0;34m%s\x1b", text)
}

func yellowText(text string) string {
	return fmt.Sprintf("\x1b[0;33m%s\x1b", text)
}
