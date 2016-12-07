package ui

import (
	"fmt"
	"math/rand"

	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

type ViewManager struct {
	contents parser.Entries
	gui      *gocui.Gui
	id       string
	pointer  int
	title    string
	view     *gocui.View
	size     size
}

type size struct {
	startX, startY, endX, endY int
}

func NewViewManager(g *gocui.Gui, title string, contents parser.Entries) *ViewManager {
	id := fmt.Sprintf("%s-%d", title, rand.Int())
	return &ViewManager{
		gui:      g,
		title:    title,
		contents: contents,
		id:       id,
	}
}

func (vm *ViewManager) View(startX, startY, endX, endY int) (*gocui.View, error) {
	if vm.view != nil {
		return vm.view, nil
	}

	if v, err := vm.gui.SetView(vm.id, startX, startY, endX, 2*endY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		vm.size = size{startX, startY, endX, endY}
		v.Highlight = true

		if err := setKeysbindings(vm.gui, vm.id, navigateDownKeys, gocui.ModNone, vm.moveDown); err != nil {
			return nil, err
		}

		if err := vm.gui.SetKeybinding(vm.id, gocui.KeyPgdn, gocui.ModNone, vm.movePageDown); err != nil {
			return nil, err
		}

		if err := setKeysbindings(vm.gui, vm.id, navigateUpKeys, gocui.ModNone, vm.moveUp); err != nil {
			return nil, err
		}

		if err := vm.gui.SetKeybinding(vm.id, gocui.KeyPgup, gocui.ModNone, vm.movePageUp); err != nil {
			return nil, err
		}

		v.Title = vm.title
		v.BgColor = gocui.ColorBlack
		printEntries(v, vm.contents)
		vm.view = v

		vm.updateInfoView(vm.contents[0])
	}

	return vm.view, nil
}

func (vm *ViewManager) SelectedEntry() *parser.Entry {
	return vm.contents[vm.pointer]
}

func (vm *ViewManager) SetCurrent() error {
	_, err := vm.gui.SetCurrentView(vm.id)
	return err
}

func (vm *ViewManager) Contents() parser.Entries {
	return vm.contents
}

func (vm *ViewManager) Name() string {
	return vm.id
}

func (vm *ViewManager) RemoveView() error {
	return vm.gui.DeleteView(vm.id)
}

func (vm *ViewManager) updateInfoView(entry *parser.Entry) error {
	viewName := fmt.Sprintf("%s-info", vm.id)
	vm.gui.DeleteView(viewName)

	if entry == nil {
		return nil
	}

	if v, err := vm.gui.SetView(viewName, vm.size.startX, 2*vm.size.endY/3+1, vm.size.endX, vm.size.endY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "details"
		printEntryInfo(v, entry)
	}
	return nil
}

func (vm *ViewManager) moveUp(g *gocui.Gui, v *gocui.View) error {
	if vm.pointer <= 0 {
		return nil
	}

	vm.pointer--
	if vm.view != nil {
		ox, oy := vm.view.Origin()
		cx, cy := vm.view.Cursor()
		if err := vm.view.SetCursor(cx, cy-1); err != nil && 0 < oy {
			if err := vm.view.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return vm.updateInfoView(vm.contents[vm.pointer])
}

func (vm *ViewManager) moveDown(g *gocui.Gui, v *gocui.View) error {
	if vm.pointer >= len(vm.contents)-1 {
		return nil
	}

	vm.pointer++
	if vm.view != nil {
		cx, cy := vm.view.Cursor()
		ox, oy := vm.view.Origin()
		if err := vm.view.SetCursor(cx, cy+1); err != nil {
			if err := vm.view.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return vm.updateInfoView(vm.contents[vm.pointer])
}

func (vm *ViewManager) movePageUp(g *gocui.Gui, v *gocui.View) error {
	_, y := g.CurrentView().Size()

	for i := 0; i < y; i++ {
		if err := vm.moveUp(g, v); err != nil {
			return err
		}
	}

	return nil
}

func (vm *ViewManager) movePageDown(g *gocui.Gui, v *gocui.View) error {
	_, y := g.CurrentView().Size()

	for i := 0; i < y; i++ {
		if err := vm.moveDown(g, v); err != nil {
			return err
		}
	}

	return nil
}
