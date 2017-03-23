package ui

import (
	"fmt"
	"math/rand"

	"github.com/jroimartin/gocui"
	"github.com/tscolari/lag/parser"
)

type PopupView struct {
	gui     *gocui.Gui
	id      string
	message string
	title   string
	view    *gocui.View
	size    size
}

func NewPopupView(g *gocui.Gui, title string, message string) *PopupView {
	id := fmt.Sprintf("%s-%d", title, rand.Int())
	return &PopupView{
		gui:     g,
		title:   title,
		id:      id,
		message: message,
	}
}

func (vm *PopupView) View(startX, startY, endX, endY int) (*gocui.View, error) {
	if vm.view != nil {
		return vm.view, nil
	}

	if v, err := vm.gui.SetView(vm.id, startX, startY, endX, 2*endY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		vm.size = size{startX, startY, endX, endY}
		v.Highlight = false

		v.Title = vm.title
		v.BgColor = gocui.ColorBlack
		vm.view = v

		fmt.Fprintf(v, "%s\n", vm.message)
	}

	return vm.view, nil
}

func (vm *PopupView) SetCurrent() error {
	_, err := vm.gui.SetCurrentView(vm.id)
	return err
}

func (vm *PopupView) SelectedEntry() *parser.Entry {
	return nil
}

func (vm *PopupView) Contents() parser.Entries {
	return parser.Entries{}
}

func (vm *PopupView) Name() string {
	return vm.id
}

func (vm *PopupView) RemoveView() error {
	return vm.gui.DeleteView(vm.id)
}
