package cli

import (
	"strings"

	"github.com/csams/doit/pkg/auth"
	"github.com/gdamore/tcell/v2"
	"github.com/go-logr/logr"
	"github.com/rivo/tview"
)

func GetApplication(log logr.Logger, app *tview.Application, oidcFlow *auth.TokenProvider) (tview.Primitive, error) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 10, 40
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter).SetSelectable(r > 0 && c > 0))

			word = (word + 1) % len(lorem)
		}
	}
	table.Select(1, 1)
	table.SetSelectedFunc(func(row, col int) {
		table.GetCell(row, 1).SetTextColor(tcell.ColorRed)
	})
	return table, nil
}
