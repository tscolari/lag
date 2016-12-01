package ui

import "github.com/fatih/color"

func redText(text string) string {
	return color.New(color.FgRed).SprintFunc()(text)
}

func blueText(text string) string {
	return color.New(color.FgBlue).SprintFunc()(text)
}

func yellowText(text string) string {
	return color.New(color.FgYellow).SprintFunc()(text)
}
