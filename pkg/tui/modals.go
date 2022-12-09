package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *CLI) newHelp(bindings []KeyBinding) {
	grid := tview.NewGrid()
	grid.SetColumns(0, -3, 0)
	grid.SetRows(0)
	table := tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(false, false)
	table.SetBorder(true)
	table.SetSeparator(tview.Borders.Vertical)
	grid.AddItem(table, 0, 1, 1, 1, 0, 0, true)

	for c, h := range []string{"Key", "Description"} {
		table.SetCell(0, c,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorViolet).
				SetSelectable(false).
				SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	r := 1
	for _, binding := range bindings {
		table.SetCell(r, 0, tview.NewTableCell(binding.Key).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(1))
		table.SetCell(r, 1, tview.NewTableCell(binding.Description).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(4))
		r += 1
	}

	hide := func() {
		c.App.SetRoot(c.Root, true)
		c.App.SetFocus(c.Root)
	}

	table.SetDoneFunc(func(key tcell.Key) {
		hide()
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			hide()
			return nil
		}
		return event
	})

	c.App.SetRoot(grid, true)
}

func (c *CLI) newQuitModal() {
	quitModal := tview.NewModal()
	quitModal.SetText("Do you want to quit?")
	quitModal.SetBackgroundColor(tcell.ColorDarkBlue)
	quitModal.SetTextColor(tcell.ColorWheat)
	quitModal.SetButtonBackgroundColor(tcell.ColorDarkViolet)
	quitModal.SetButtonTextColor(tcell.ColorWheat)

	quitModal.AddButtons([]string{"Yes", "No"})

	quitModal.SetDoneFunc(func(i int, l string) {
		switch l {
		case "Yes":
			c.App.Stop()
		case "No":
			c.App.SetRoot(c.Root, true)
			c.App.SetFocus(c.Root)
		}
	})
	c.App.SetRoot(quitModal, false)
	c.App.SetFocus(quitModal)
}

func (c *CLI) newMessageModal(title, msg string) {
	modal := tview.NewModal()
	modal.SetTitle(title)
	modal.SetText(msg)
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	modal.SetTextColor(tcell.ColorWheat)
	modal.SetButtonBackgroundColor(tcell.ColorDarkViolet)
	modal.SetButtonTextColor(tcell.ColorWheat)

	modal.AddButtons([]string{"OK"})

	modal.SetDoneFunc(func(i int, l string) {
		c.App.SetRoot(c.Root, true)
		c.App.SetFocus(c.Root)
	})

	c.App.SetRoot(modal, false)
	c.App.SetFocus(modal)
}

func (c *CLI) newErrorModal(msg string) {
	c.newMessageModal("Error", msg)
}
