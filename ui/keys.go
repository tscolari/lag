package ui

import "github.com/jroimartin/gocui"

var navigateDownKeys = []interface{}{
	gocui.KeyArrowDown,
	'j',
}

var navigateUpKeys = []interface{}{
	gocui.KeyArrowUp,
	'k',
}

var navigateZoomOutKeys = []interface{}{
	gocui.KeyBackspace,
	gocui.KeyBackspace2,
	gocui.KeyArrowLeft,
	'h',
	'q',
}

var navigateZoomInKeys = []interface{}{
	gocui.KeyEnter,
	gocui.KeyArrowRight,
	'l',
}

func setKeysbindings(g *gocui.Gui, viewname string, keys []interface{}, mod gocui.Modifier, handler func(*gocui.Gui, *gocui.View) error) error {
	for _, key := range keys {
		if err := g.SetKeybinding(viewname, key, mod, handler); err != nil {
			return nil
		}
	}

	return nil
}
