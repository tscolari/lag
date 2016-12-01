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

func whiteInBlueText(text string) string {
	return color.New(color.FgWhite, color.BgBlue).SprintFunc()(text)
}

func whiteInRedText(text string) string {
	return color.New(color.FgWhite, color.BgRed).SprintFunc()(text)
}

func whiteInMagentaText(text string) string {
	return color.New(color.FgWhite, color.BgHiMagenta).SprintFunc()(text)
}
