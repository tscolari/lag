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

	if v, err := vm.gui.SetView(vm.id, startX, startY, endX, endY); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		v.Highlight = true

		if err := vm.gui.SetKeybinding(vm.id, gocui.KeyArrowDown, gocui.ModNone, vm.moveDown); err != nil {
			return nil, err
		}

		if err := vm.gui.SetKeybinding(vm.id, gocui.KeyArrowUp, gocui.ModNone, vm.moveUp); err != nil {
			return nil, err
		}

		v.Title = vm.title
		printEntries(v, vm.contents)
		vm.view = v
	}

	return vm.view, nil
}

func (vm *ViewManager) SelectedEntry() *parser.Entry {
	return vm.contents[vm.pointer]
}

func (vm *ViewManager) Name() string {
	return vm.id
}

func (vm *ViewManager) moveUp(g *gocui.Gui, v *gocui.View) error {
	if vm.pointer > 0 {
		vm.pointer--
	}

	if vm.view != nil {
		ox, oy := vm.view.Origin()
		cx, cy := vm.view.Cursor()
		if err := vm.view.SetCursor(cx, cy-1); err != nil && 0 < oy {
			if err := vm.view.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *ViewManager) moveDown(g *gocui.Gui, v *gocui.View) error {
	if len(vm.contents) > vm.pointer {
		vm.pointer++
	}

	if vm.view != nil {
		cx, cy := vm.view.Cursor()
		ox, oy := vm.view.Origin()
		//  ok,  := movm.viewable(vm.view, oy+cy+1)
		//  _, maxY := vm.view.Size()
		if err := vm.view.SetCursor(cx, cy+1); err != nil {
			if err := vm.view.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return nil
}
