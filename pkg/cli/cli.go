package cli

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

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

type taskFormData struct {
	Description string
	Due         string
	State       apis.State
	Status      apis.Status
	Priority    apis.Priority
	Private     bool
}

func formDataFromTask(task *apis.Task) *taskFormData {
	fmtString := "2006-01-02 15:04:05 MST"
	var due string
	if task.Due != nil {
		due = task.Due.Format(fmtString)
	}
	return &taskFormData{
		Description: task.Description,
		Due:         due,
		State:       task.State,
		Status:      task.Status,
		Priority:    task.Priority,
		Private:     task.Private,
	}
}

func taskFromFormData(data *taskFormData) (*apis.Task, error) {
	var due *time.Time
	if data.Due != "" {
		if dueDate, err := dateparse.ParseLocal(data.Due); err != nil {
			return nil, err
		} else {
			due = &dueDate
		}
	}
	return &apis.Task{
		Description: data.Description,
		Due:         due,
		State:       data.State,
		Status:      data.Status,
		Priority:    data.Priority,
		Private:     data.Private,
	}, nil
}

func (c *CLI) newTaskForm(table *TaskTable, t *taskFormData, title string, save func() error) *tview.Form {
	form := styledForm()
	form.SetTitle(title)

	ensureInt := func(t string, l rune) bool { _, err := strconv.Atoi(t); return err == nil }

	form.AddInputField("Description", t.Description, 0, nil, func(text string) { t.Description = text })
	form.AddInputField("Due", t.Due, 30, nil, func(text string) { t.Due = text })
	form.AddDropDown("State", []string{"undefined", "open", "closed"}, getStateIndex(t.State), func(option string, index int) { t.State = apis.State(option) })
	form.AddDropDown("Status", []string{"undefined", "backlog", "todo", "doing", "done", "abandoned"}, getStatusIndex(t.Status), func(option string, index int) { t.Status = apis.Status(option) })
	form.AddInputField("Priority", strconv.Itoa(int(t.Priority)), 3, ensureInt, func(text string) {
		p, _ := strconv.Atoi(text)
		t.Priority = apis.Priority(p)
	})
	form.AddCheckbox("Private", t.Private, func(checked bool) { t.Private = checked })

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			if err := save(); err != nil {
				c.newErrorModal("Error saving task: " + err.Error())
				c.Root.RemoveItem(form)
			} else {
				c.Root.RemoveItem(form)
				c.App.SetFocus(table)
				table.Update(false)
			}
			return nil
		}
		return event
	})

	cancel := func() { c.Root.RemoveItem(form); c.App.SetFocus(table.Table) }
	form.SetCancelFunc(cancel)
	form.AddButton("Cancel", cancel)
	form.AddButton("Save", func() {
		if err := save(); err != nil {
			c.newErrorModal("Error saving task: " + err.Error())
			c.Root.RemoveItem(form)
			return
		} else {
			c.Root.RemoveItem(form)
			c.App.SetFocus(table)
			table.Update(false)
		}
	})

	c.Root.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}

func (c *CLI) newDeleteModal(table *TaskTable, orig *apis.Task) *tview.Modal {
	modal := tview.NewModal()
	modal.SetTitle("Delete?")
	modal.SetText("Do you want to delete task [" + orig.Description + "]")
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

		formData := formDataFromTask(t)
		form := c.newTaskForm(tt, formData, "Edit task", func() error {
			proposedTask, err := taskFromFormData(formData)
			if err != nil {
				return err
			}
			proposedTask.ID = t.ID
			userId := fmt.Sprintf("%d", c.Me.ID)
			taskId := fmt.Sprintf("%d", t.ID)
			up, err := generic.Put(c.Client, "users/"+userId+"/tasks/"+taskId, proposedTask)
			if err != nil {
				return err
			}
			*t = *up
			return nil
		})
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
			day := 24 * time.Hour
			due := time.Now().Add(day).Round(day)
			orig := &apis.Task{State: apis.Open, Status: apis.Backlog, Due: &due}
			formData := formDataFromTask(orig)
			form := c.newTaskForm(tt, formData, "Create task", func() error {
				t, err := taskFromFormData(formData)
				if err != nil {
					return err
				}
				userId := fmt.Sprintf("%d", c.Me.ID)
				up, err := generic.Post(c.Client, "users/"+userId+"/tasks", t)
				if err != nil {
					return err
				}
				tt.Tasks = append(tt.Tasks, *up)
				return nil
			})
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
	statusOrder := map[apis.Status]int{
		apis.Doing:     0,
		apis.Todo:      1,
		apis.Backlog:   2,
		apis.Done:      3,
		apis.Abandoned: 4,
	}

	// primary sort by status and then secondary sorts by due date and priority
	sort.SliceStable(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	sort.SliceStable(tasks, func(i, j int) bool {
		if tasks[i].Due == nil {
			return false
		}
		if tasks[j].Due == nil {
			return true
		}
		return tasks[i].Due.Before(*tasks[j].Due)
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
