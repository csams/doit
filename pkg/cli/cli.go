package cli

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/araddon/dateparse"
	"github.com/csams/doit/pkg/apis"
	generic "github.com/csams/doit/pkg/cli/client"
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
	Root *tview.Flex
	Me   *apis.User
}

func New(cfg CompletedConfig) (*CLI, error) {
	c := &CLI{
		CompletedConfig: cfg,
		Root:            tview.NewFlex(),
	}

	me, err := generic.Get[apis.User](c.Client, "me")
	if err != nil {
		return nil, err
	}

	c.Me = me

	table := NewTaskTable(c, me.AssignedTasks)
	table.SetTitle("Tasks for " + me.Name)

	c.Root.AddItem(table, 0, 1, true) // (item, fixedSize; 0 means not fixed, proportion, focus?)

	return c, nil
}

func (c *CLI) newQuitModal() {
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
			c.App.Stop()
		case "No":
			c.App.SetRoot(c.Root, true)
			c.App.SetFocus(c.Root)
		}
	})
	c.App.SetRoot(quitModal, false)
	c.App.SetFocus(quitModal)
}

func (c *CLI) newErrorModal(msg string) {
	modal := tview.NewModal()
	modal.SetTitle("Error")
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

func (c *CLI) newCreateTaskForm(table *TaskTable, orig *apis.Task) *tview.Form {
	form := styledForm()
	form.SetTitle("Create Task")

	ensureInt := func(t string, l rune) bool { _, err := strconv.Atoi(t); return err == nil }

	var due string

	fmtString := "2006-01-02 15:04:05 MST"
	if orig.Due != nil {
		due = orig.Due.Format(fmtString)
	}

	form.AddInputField("Description", "", 0, nil, func(text string) { orig.Description = text })
	form.AddInputField("Due", due, 30, nil, func(text string) { due = text })
	form.AddDropDown("State", []string{"open", "closed"}, 0, func(option string, index int) { orig.State = apis.State(option) })
	form.AddDropDown("Status", []string{"backlog", "todo", "doing", "done", "abandoned"}, 0, func(option string, index int) { orig.Status = apis.Status(option) })
	form.AddInputField("Priority", "0", 3, ensureInt, func(text string) {
		p, _ := strconv.Atoi(text)
		orig.Priority = apis.Priority(p)
	})
	form.AddCheckbox("Private", false, func(checked bool) { orig.Private = checked })

	cancel := func() { c.Root.RemoveItem(form); c.App.SetFocus(table.Table) }
	save := func() {
		if due != "" {
			dueDate, err := dateparse.ParseLocal(due)
			if err != nil {
				c.newErrorModal("Error parsing due date: " + err.Error())
				c.Root.RemoveItem(form)
				return
			}
			orig.Due = &dueDate
		}
		userId := fmt.Sprintf("%d", c.Me.ID)
		up, err := generic.Post(c.Client, "users/"+userId+"/tasks", orig)
		if err != nil {
			c.newErrorModal("Error creating task: " + err.Error())
			c.Root.RemoveItem(form)
			return
		}
		table.Tasks = append(table.Tasks, *up)
		c.Root.RemoveItem(form)
		c.App.SetFocus(table.Table)
		table.Update(false)
	}

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			save()
			return nil
		}
		return event
	})

	form.SetCancelFunc(cancel)
	form.AddButton("Cancel", cancel)
	form.AddButton("Save", save)

	c.Root.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}

func (c *CLI) newEditTaskForm(table *TaskTable, orig *apis.Task) *tview.Form {
	t := *orig // shallow copy for now...

	form := styledForm()
	form.SetTitle("Edit Task")

	var due string

	fmtString := "2006-01-02 15:04:05 MST"
	if orig.Due != nil {
		due = orig.Due.Format(fmtString)
	}

	form.AddInputField("Description", t.Description, 0, nil, func(text string) { t.Description = text })
	form.AddInputField("Due", due, 30, nil, func(text string) { due = text })
	form.AddDropDown("State", []string{"undefined", "open", "closed"}, getStateIndex(t.State), func(option string, index int) { t.State = apis.State(option) })
	form.AddDropDown("Status", []string{"undefined", "backlog", "todo", "doing", "done", "abandoned"}, getStatusIndex(t.Status), func(option string, index int) { t.Status = apis.Status(option) })
	form.AddInputField("Priority", strconv.Itoa(int(t.Priority)), 3, func(t string, l rune) bool { _, err := strconv.Atoi(t); return (err == nil) }, func(text string) {
		p, _ := strconv.Atoi(text)
		t.Priority = apis.Priority(p)
	})
	form.AddCheckbox("Private", t.Private, func(checked bool) { t.Private = checked })

	cancel := func() { c.Root.RemoveItem(form); c.App.SetFocus(table.Table) }
	save := func() {
		if due != "" {
			dueDate, err := dateparse.ParseLocal(due)
			if err != nil {
				c.newErrorModal("Error editing task: " + err.Error())
				c.Root.RemoveItem(form)
				return
			}
			t.Due = &dueDate
		}
		userId := fmt.Sprintf("%d", c.Me.ID)
		taskId := fmt.Sprintf("%d", t.ID)
		up, err := generic.Put(c.Client, "users/"+userId+"/tasks/"+taskId, &t)
		if err != nil {
			c.newErrorModal("Error editing task: " + err.Error())
			c.Root.RemoveItem(form)
			return
		}
		*orig = *up
		c.Root.RemoveItem(form)
		c.App.SetFocus(table.Table)
		table.Update(false)
	}

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			save()
			return nil
		}
		return event
	})

	form.SetCancelFunc(cancel)
	form.AddButton("Cancel", cancel)
	form.AddButton("Save", save)

	c.Root.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}

func (c *CLI) newDeleteModal(table *TaskTable, orig *apis.Task) *tview.Modal {
	modal := tview.NewModal()
	modal.SetTitle("Delete?")
	modal.SetText("Do you want to delete task [")
	modal.SetBackgroundColor(tcell.ColorDarkBlue)
	modal.SetTextColor(tcell.ColorWheat)
	modal.SetButtonBackgroundColor(tcell.ColorDarkViolet)
	modal.SetButtonTextColor(tcell.ColorWheat)

	modal.AddButtons([]string{"Yes", "No"})

	cancel := func() { c.App.SetRoot(c.Root, true); c.App.SetFocus(table.Table) }
	modal.SetDoneFunc(func(i int, l string) {
		switch l {
		case "Yes":
			userId := fmt.Sprintf("%d", c.Me.ID)
			taskId := fmt.Sprintf("%d", orig.ID)
			_, err := generic.Delete[apis.Task](c.Client, "users/"+userId+"/tasks/"+taskId)
			if err != nil {
				c.newErrorModal(err.Error())
				c.Root.RemoveItem(modal)
				return
			}
			if err = table.Remove(orig); err != nil {
				c.App.SetRoot(c.Root, true)
				c.App.SetFocus(table.Table)
				return
			}
			table.Update(true)
			c.App.SetRoot(c.Root, true)
			c.App.SetFocus(table.Table)
		case "No":
			cancel()
		}
	})
	c.App.SetRoot(modal, false)
	c.App.SetFocus(modal)
	return modal
}

type TaskTable struct {
	*tview.Table
	CLI   *CLI
	Tasks []apis.Task
}

func NewTaskTable(c *CLI, tasks []apis.Task) *TaskTable {
	table := tview.NewTable().
		SetFixed(1, 1).            // the first column and first row are visible even when scrolling
		SetSelectable(true, false) // (rows?, cols?) - select entire rows not columns
	table.SetBorder(true)
	table.SetSeparator(tview.Borders.Vertical)

	tt := &TaskTable{
		CLI:   c,
		Tasks: tasks,
		Table: table,
	}

	table.SetSelectedFunc(func(row, col int) {
		t := table.GetCell(row, 0).GetReference().(*apis.Task)

		form := c.newEditTaskForm(tt, t)
		c.App.SetFocus(form)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			c.newQuitModal()
			return nil
		}

		switch event.Rune() {
		case 'd':
			row, _ := table.GetSelection()
			t := table.GetCell(row, 0).GetReference().(*apis.Task)
			modal := c.newDeleteModal(tt, t)
			c.App.SetFocus(modal)
			return nil
		case 'n':
			orig := &apis.Task{State: apis.Open, Status: apis.Backlog}
			form := c.newCreateTaskForm(tt, orig)
			c.App.SetFocus(form)
			return nil
		case 'q':
			c.newQuitModal()
			return nil
		case 'Q':
			c.App.Stop()
			return nil
		}

		return event // returning the event means other handlers also see it
	})

	tt.Update(false)
	return tt
}

func (t *TaskTable) Update(clear bool) {
	table := t.Table
	tasks := t.Tasks
	curR, curC := table.GetSelection()
	if clear {
		t.Clear()
	}
	for c, h := range table_headers {
		table.SetCell(0, c,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorViolet).
				SetSelectable(false).
				SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	fmtString := "2006-01-02 15:04:05 MST"
	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	statusOrder := map[apis.Status]int{
		apis.Doing:     0,
		apis.Todo:      1,
		apis.Backlog:   2,
		apis.Done:      3,
		apis.Abandoned: 4,
	}

	sort.SliceStable(tasks, func(i, j int) bool {
		if tasks[i].Due == nil {
			return false
		}
		if tasks[j].Due == nil {
			return true
		}
		return tasks[i].Due.After(*tasks[j].Due)
	})

	sort.SliceStable(tasks, func(i, j int) bool {
		return statusOrder[tasks[i].Status] < statusOrder[tasks[j].Status]
	})

	for r := range tasks {
		task := &tasks[r]
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
	if curR == 0 || curC == 0 {
		curR = 1
		curC = 1
	}
	table.Select(curR, curC)
}

func (t *TaskTable) Remove(task *apis.Task) error {
	toRemove := -1
	for i, e := range t.Tasks {
		if e.ID == task.ID {
			toRemove = i
			break
		}
	}

	taskId := fmt.Sprintf("%d", task.ID)
	if toRemove == -1 {
		return errors.New("could not find task to delete: " + taskId)
	}

	newTasks := make([]apis.Task, 0)
	newTasks = append(newTasks, t.Tasks[:toRemove]...)
	newTasks = append(newTasks, t.Tasks[toRemove+1:]...)

	t.Tasks = newTasks
	return nil
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
	table_headers = []string{
		"Id",
		"Created",
		"Description",
		"Due",
		"Priority",
		"State",
		"Status",
		"Private",
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
