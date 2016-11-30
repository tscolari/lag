package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

type UI struct {
	entries parser.Entries
	pointer int
	zoom    *parser.Entry
	gui     *gocui.Gui
	views   []*gocui.View
}

func New(entries parser.Entries) *UI {
	return &UI{
		entries: entries,
		zoom:    nil,
		views:   []*gocui.View{},
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

func (ui *UI) moveUp(g *gocui.Gui, v *gocui.View) error {
	if ui.pointer > 0 {
		ui.pointer--
	}

	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && 0 < oy {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ui *UI) addView(g *gocui.Gui) error {
	return ui.renderEntries(g, ui.zoom.Children, ui.zoom)
}

func (ui *UI) renderEntries(g *gocui.Gui, entries parser.Entries, parent *parser.Entry) error {
	maxX, maxY := g.Size()
	viewName := fmt.Sprintf("entries-%d", len(ui.views))

	if v, err := g.SetView(viewName, 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true

		if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone, ui.zoomIn); err != nil {
			return err
		}

		if err := g.SetKeybinding(viewName, gocui.KeyArrowLeft, gocui.ModNone, ui.zoomOut); err != nil {
			return err
		}

		if err := g.SetKeybinding(viewName, gocui.KeyArrowDown, gocui.ModNone, ui.moveDown); err != nil {
			return err
		}

		if err := g.SetKeybinding(viewName, gocui.KeyArrowUp, gocui.ModNone, ui.moveUp); err != nil {
			return err
		}

		ui.views = append(ui.views, v)
		if parent != nil {
			v.Title = parent.Data.Message
		}
		printEntries(v, entries)
	}

	_, err := g.SetCurrentView(viewName)
	return err
}

func (ui *UI) moveDown(g *gocui.Gui, v *gocui.View) error {
	if len(ui.entries) > ui.pointer {
		ui.pointer++
	}

	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		//  ok,  := movable(v, oy+cy+1)
		//  _, maxY := v.Size()
		if err := v.SetCursor(cx, cy+1); err != nil {
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ui *UI) zoomIn(g *gocui.Gui, v *gocui.View) error {
	if ui.zoom == nil {
		if len(ui.entries[ui.pointer].Children) == 0 {
			return nil
		}

		ui.zoom = ui.entries[ui.pointer]
	} else {
		if len(ui.zoom.Children[ui.pointer].Children) == 0 {
			return nil
		}

		ui.zoom = ui.zoom.Children[ui.pointer]
	}

	ui.pointer = 0
	ui.addView(g)
	return nil
}

func (ui *UI) zoomOut(g *gocui.Gui, v *gocui.View) error {
	if ui.zoom == nil {
		return nil
	}

	ui.pointer = 0
	ui.zoom = ui.zoom.Parent
	curViewIdx := len(ui.views) - 1

	if err := g.DeleteView(ui.views[curViewIdx].Name()); err != nil {
		return err
	}

	ui.views = ui.views[:curViewIdx]
	if _, err := g.SetCurrentView(ui.views[curViewIdx-1].Name()); err != nil {
		return err
	}

	return nil
}
