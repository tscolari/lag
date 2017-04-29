package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

type UI struct {
	entries      parser.Entries
	pointer      int
	gui          *gocui.Gui
	viewManagers []View
}

type View interface {
	SelectedEntry() *parser.Entry
	RemoveView() error
	SetCurrent() error
	Contents() parser.Entries
}

func New(entries parser.Entries) *UI {
	return &UI{
		entries:      entries,
		viewManagers: []View{},
	}
}

func (ui *UI) Start() error {
	var err error
	ui.gui, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		return err
	}
	defer ui.gui.Close()
	ui.gui.BgColor = gocui.ColorBlack
	ui.gui.SelBgColor = gocui.ColorGreen
	ui.gui.SelFgColor = gocui.ColorBlack

	ui.gui.SetManagerFunc(ui.render)
	ui.renderEntries(ui.gui, ui.entries, nil)

	if err := ui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.quit); err != nil {
		return err
	}

	if err := ui.gui.SetKeybinding("", 'Q', gocui.ModNone, ui.quit); err != nil {
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

	if err := setKeysbindings(ui.gui, viewName, navigateZoomInKeys, gocui.ModNone, ui.zoomIn); err != nil {
		return err
	}

	if err := setKeysbindings(ui.gui, viewName, navigateZoomOutKeys, gocui.ModNone, ui.zoomOut); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'e', gocui.ModNone, ui.filterErrored); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'D', gocui.ModNone, ui.deleteSimilar); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'i', gocui.ModNone, ui.sessionInfo); err != nil {
		return err
	}

	ui.viewManagers = append(ui.viewManagers, viewManager)
	return viewManager.SetCurrent()
}

func (ui *UI) zoomIn(g *gocui.Gui, v *gocui.View) error {
	topViewManager := ui.viewManagers[len(ui.viewManagers)-1]
	selEntry := topViewManager.SelectedEntry()

	if len(selEntry.Children) == 0 {
		return nil
	}

	return ui.renderEntries(g, selEntry.Children, selEntry)
}

func (ui *UI) zoomOut(g *gocui.Gui, v *gocui.View) error {
	if len(ui.viewManagers) <= 1 {
		return nil
	}

	totalViews := len(ui.viewManagers)
	topViewManager := ui.viewManagers[totalViews-1]

	if err := topViewManager.RemoveView(); err != nil {
		return err
	}

	ui.viewManagers = ui.viewManagers[:totalViews-1]
	if err := ui.viewManagers[totalViews-2].SetCurrent(); err != nil {
		return err
	}

	return nil
}

func (ui *UI) filterErrored(g *gocui.Gui, v *gocui.View) error {
	topViewManager := ui.viewManagers[len(ui.viewManagers)-1]
	currentContents := topViewManager.Contents()
	erroredEntries := currentContents.ErroredOnly()

	if len(erroredEntries) == 0 {
		return nil
	}

	return ui.renderEntries(g, erroredEntries, nil)
}

func (ui *UI) deleteSimilar(g *gocui.Gui, v *gocui.View) error {
	topViewManager := ui.viewManagers[len(ui.viewManagers)-1]
	currentContents := topViewManager.Contents()
	filteredContents := currentContents.RemoveSimilar(*topViewManager.SelectedEntry())

	if len(filteredContents) == 0 {
		return nil
	}

	return ui.renderEntries(g, filteredContents, nil)
}

func (ui *UI) sessionInfo(g *gocui.Gui, v *gocui.View) error {
	topViewManager := ui.viewManagers[len(ui.viewManagers)-1]
	entry := topViewManager.SelectedEntry()
	entries := topViewManager.Contents()
	sessionDuration := entries.SessionDuration(entry.Data.Session)
	message := fmt.Sprintf("Session duration: %s", sessionDuration.String())

	title := fmt.Sprintf("Session info: [%s]", entry.Data.Session)
	view := NewPopupView(g, title, message)

	maxX, maxY := g.Size()
	totalX := (maxX - 1) / 2
	totalY := (maxY - 1) / 3
	startX := (maxX - totalX) / 2
	startY := (maxY - totalY) / 2
	_, err := view.View(startX, startY, startX+totalX, startY+totalY)
	if err != nil {
		return err
	}

	if err := setKeysbindings(ui.gui, view.Name(), navigateZoomOutKeys, gocui.ModNone, ui.zoomOut); err != nil {
		return err
	}

	if err := g.SetKeybinding(view.Name(), 'i', gocui.ModNone, ui.zoomOut); err != nil {
		return err
	}

	ui.viewManagers = append(ui.viewManagers, view)
	return view.SetCurrent()
	return nil
}
