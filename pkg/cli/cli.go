package cli

import (
	"fmt"
	"strconv"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/cli/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
Example Task:
{
  "id": 3,
  "created_at": "2022-12-02T16:58:15.698228072-06:00",
  "updated_at": "2022-12-02T16:58:15.698228072-06:00",
  "owner_id": 1,
  "assignee_id": 1,
  "desc": "",
  "due": null,
  "priority": 10,
  "private": false,
  "state": "open",
  "status": "",
  "comments": null,
  "annotations":
}
*/

type CLI struct {
	CompletedConfig
	Flex *tview.Flex
	Me   *apis.User
}

func New(cfg CompletedConfig) (*CLI, error) {
	flex := tview.NewFlex().SetFullScreen(true)

	c := &CLI{
		CompletedConfig: cfg,
		Flex:            flex,
	}

	table := tview.NewTable().
		SetFixed(1, 1).            // the first column and first row always visible even when scrolling
		SetSelectable(true, false) // (rows?, cols?) - select entire rows not columns
	table.SetBorder(true)
	table.SetSeparator(tview.Borders.Vertical)

	flex.AddItem(table, 0, 1, true) // (item, fixedSize; 0 means not fixed, proportion, focus?)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			newQuitModal(c.App, flex)
			return nil
		case tcell.KeyEnter:
			r, _ := table.GetSelection()
			t := table.GetCell(r, 0).GetReference().(apis.Task)

			form := c.editTaskForm(&t)
			c.App.SetFocus(form)
			return nil
		}

		switch event.Rune() {
		case 'Q', 'q':
			newQuitModal(c.App, flex)
			return nil
		case 'n':
			form := c.createTaskForm()
			c.App.SetFocus(form)
			return nil
		}

		return event // returning the event means other handlers also see it
	})

	user, err := client.Get[apis.User](c.Client, "http://localhost:9090/me", c.Tokens)
	if err != nil {
		return nil, err
	}

	c.Me = user

	table.SetTitle("Tasks for " + user.Name)

	table_headers := []string{
		"Id",
		"Created",
		"Description",
		"Due",
		"Priority",
		"State",
		"Status",
		"Private",
	}
	for c, h := range table_headers {
		table.SetCell(0, c,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorViolet).
				SetSelectable(false).
				SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	tasks := user.AssignedTasks

	fmtString := "2006-01-02 15:04:05 MST"
	for r, task := range tasks {
		r = r + 1
		id := fmt.Sprintf("%d", task.ID)
		priority := fmt.Sprintf("%d", task.Priority)
		createdAt := task.CreatedAt.Format(fmtString)

		var due string
		if task.Due != nil {
			due = task.Due.Format(fmtString)
		}

		table.SetCell(r, 0, tview.NewTableCell(id).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetReference(task))
		table.SetCell(r, 1, tview.NewTableCell(string(createdAt)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 2, tview.NewTableCell(task.Description).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(4))
		table.SetCell(r, 3, tview.NewTableCell(due).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 4, tview.NewTableCell(priority).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 5, tview.NewTableCell(string(task.State)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 6, tview.NewTableCell(string(task.Status)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 7, tview.NewTableCell(privateMap[task.Private]).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
	}
	table.Select(1, 1)
	return c, nil
}

var (
	privateMap = map[bool]string{true: "✓", false: "✗"}
	statusMap  = map[apis.Status]int{
		apis.Status(""): 0,
		apis.Backlog:    1,
		apis.Todo:       2,
		apis.Doing:      3,
		apis.Done:       4,
		apis.Abandoned:  5,
	}
	stateMap = map[apis.State]int{
		apis.State(""): 0,
		apis.Open:      1,
		apis.Closed:    2,
	}
)

func getStatusIndex(s apis.Status) int {
	i, found := statusMap[s]
	if !found {
		return 0
	}
	return i
}

func getStateIndex(s apis.State) int {
	i, found := stateMap[s]
	if !found {
		return 0
	}
	return i
}

func newQuitModal(app *tview.Application, oldRoot tview.Primitive) {
	prevFocus := app.GetFocus()

	quitModal := tview.NewModal()
	quitModal.SetTitle("Quit?")
	quitModal.SetText("Do you want to quit?")
	quitModal.SetBackgroundColor(tcell.ColorDarkBlue)
	quitModal.SetTextColor(tcell.ColorWheat)
	quitModal.SetButtonBackgroundColor(tcell.ColorDarkViolet)
	quitModal.SetButtonTextColor(tcell.ColorWheat)

	quitModal.AddButtons([]string{"Yes", "No"})

	quitModal.SetDoneFunc(func(i int, l string) {
		switch l {
		case "Yes":
			app.Stop()
		case "No":
			app.SetRoot(oldRoot, true)
			app.SetFocus(prevFocus)
		}
	})
	app.SetRoot(quitModal, false)
	app.SetFocus(quitModal)
}

func styledForm() *tview.Form {
	form := tview.NewForm()
	form.SetTitleAlign(tview.AlignLeft)
	form.SetBorder(true)
	form.SetFieldTextColor(tcell.ColorWheat)
	form.SetFieldBackgroundColor(tcell.ColorDarkBlue)

	form.SetButtonBackgroundColor(tcell.ColorDarkViolet)
	form.SetButtonTextColor(tcell.ColorWheat)

	return form
}

func (c *CLI) createTaskForm() *tview.Form {
	oldFocus := c.App.GetFocus()
	form := styledForm()
	form.SetTitle("Create Task")

	form.AddTextArea("Description", "", 0, 0, 0, nil)
	form.AddDropDown("State", []string{"undefined", "open", "closed"}, 1, nil)
	form.AddDropDown("Status", []string{"undefined", "backlog", "todo", "doing", "done", "abandoned"}, 1, nil)
	form.AddInputField("Priority", "0", 3, func(t string, l rune) bool { _, err := strconv.Atoi(t); return (err == nil) }, nil)
	form.AddCheckbox("Private", true, nil)

	form.AddButton("Cancel", func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })
	form.AddButton("Save", func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })
	form.SetCancelFunc(func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })

	c.Flex.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}

func (c *CLI) editTaskForm(t *apis.Task) *tview.Form {
	oldFocus := c.App.GetFocus()

	form := styledForm()
	form.SetTitle("Edit Task")

	form.AddTextArea("Description", t.Description, 0, 0, 0, nil)
	form.AddDropDown("State", []string{"undefined", "open", "closed"}, getStateIndex(t.State), nil)
	form.AddDropDown("Status", []string{"undefined", "backlog", "todo", "doing", "done", "abandoned"}, getStatusIndex(t.Status), nil)
	form.AddInputField("Priority", strconv.Itoa(int(t.Priority)), 3, func(t string, l rune) bool { _, err := strconv.Atoi(t); return (err == nil) }, nil)
	form.AddCheckbox("Private", t.Private, nil)

	form.AddButton("Cancel", func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })
	form.AddButton("Save", func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })
	form.SetCancelFunc(func() { c.Flex.RemoveItem(form); c.App.SetFocus(oldFocus) })

	c.Flex.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}
