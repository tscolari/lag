package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

type UI struct {
	entries parser.Entries
	pointer int
	gui     *gocui.Gui
	views   []*ViewManager
}

func New(entries parser.Entries) *UI {
	return &UI{
		entries: entries,
		views:   []*ViewManager{},
	}
}

func (ui *UI) Start() error {
	var err error
	ui.gui, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	defer ui.gui.Close()
	ui.gui.SelBgColor = gocui.ColorGreen
	ui.gui.SelFgColor = gocui.ColorBlack

	ui.gui.SetManagerFunc(ui.render)

	ui.renderEntries(ui.gui, ui.entries, nil)
	ui.gui.SetCurrentView("entries-0")

	if err := ui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}

	if err := ui.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (ui *UI) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) render(g *gocui.Gui) error {
	return nil
}

// func (ui *UI) addView(g *gocui.Gui) error {
// 	return ui.renderEntries(g, ui.zoom.Children, ui.zoom)
// }

func (ui *UI) renderEntries(g *gocui.Gui, entries parser.Entries, parent *parser.Entry) error {
	maxX, maxY := g.Size()
	title := ""
	if parent != nil {
		title = parent.Data.Message
	}

	viewManager := NewViewManager(g, title, entries)
	_, err := viewManager.View(0, 0, maxX-1, maxY-1)
	if err != nil {
		return err
	}

	viewName := viewManager.Name()
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, ui.zoomIn); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, gocui.KeyArrowLeft, gocui.ModNone, ui.zoomOut); err != nil {
		return err
	}

	ui.views = append(ui.views, viewManager)
	_, err = g.SetCurrentView(viewName)
	return err
}

func (ui *UI) zoomIn(g *gocui.Gui, v *gocui.View) error {
	topViewManager := ui.views[len(ui.views)-1]
	selEntry := topViewManager.SelectedEntry()

	if len(selEntry.Children) == 0 {
		return nil
	}

	return ui.renderEntries(g, selEntry.Children, selEntry)
}

func (ui *UI) zoomOut(g *gocui.Gui, v *gocui.View) error {
	if len(ui.views) <= 1 {
		return nil
	}

	totalViews := len(ui.views)
	topViewManager := ui.views[totalViews-1]

	if err := g.DeleteView(topViewManager.Name()); err != nil {
		return err
	}

	ui.views = ui.views[:totalViews-1]
	if _, err := g.SetCurrentView(ui.views[totalViews-2].Name()); err != nil {
		return err
	}

	return nil
}
