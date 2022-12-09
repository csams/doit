package tui

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/csams/doit/pkg/apis"
	generic "github.com/csams/doit/pkg/tui/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TaskTable struct {
	*tview.Table
	CLI   *CLI
	Tasks []TaskModel
}

func (t *TaskTable) editTask(task *TaskModel) {
	form := t.CLI.newTaskForm(t, task.Task, "Edit task", func(formData *taskFormData) error {
		proposedTask := *task.Task
		err := formData.ApplyTo(&proposedTask)
		if err != nil {
			return err
		}
		userId := fmt.Sprintf("%d", t.CLI.Me.ID)
		taskId := fmt.Sprintf("%d", task.ID)
		up, err := generic.Put(t.CLI.Client, "users/"+userId+"/tasks/"+taskId, &proposedTask)
		if err != nil {
			return err
		}
		*task.Task = *up
		task.LastTouched = true
		return nil
	})
	t.CLI.App.SetFocus(form)
}

type TaskModel struct {
	*apis.Task
	LastTouched bool
}

func NewTaskTable(c *CLI, tasks []apis.Task) *TaskTable {
	table := tview.NewTable().
		SetFixed(1, 1). // the first column and first row are visible even when scrolling
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical) // (rows?, cols?) - select entire rows not columns
	table.SetBorder(true)

	tt := &TaskTable{
		CLI:   c,
		Table: table,
	}

	tt.SetTasks(tasks)

	table.SetSelectedFunc(func(row, col int) {
		ref := table.GetCell(row, 0).GetReference()
		if ref == nil {
			return
		}
		task := ref.(*TaskModel)
		tt.editTask(task)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			c.newQuitModal()
			return nil
		}

		switch event.Rune() {
		case 'a':
			userId := fmt.Sprintf("%d", c.Me.ID)
			taskList, err := generic.Get[apis.TaskList](c.Client, "users/"+userId+"/tasks?assignee=1")
			if err != nil {
				c.newErrorModal(err.Error())
				return nil
			}
			tt.SetTasks(taskList.Tasks)
			tt.SetTitle("Tasks assigned to " + c.Me.Username)
			return nil
		case 'o':
			userId := fmt.Sprintf("%d", c.Me.ID)
			taskList, err := generic.Get[apis.TaskList](c.Client, "users/"+userId+"/tasks")
			if err != nil {
				c.newErrorModal(err.Error())
				return nil
			}
			tt.SetTasks(taskList.Tasks)
			tt.SetTitle("Tasks owned by " + c.Me.Username)
			return nil

		case '?':
			c.newHelp(taskTableKeyBindings)
			return nil
		case 'd':
			row, _ := table.GetSelection()
			ref := table.GetCell(row, 0).GetReference()
			if ref != nil {
				t := ref.(*TaskModel)
				modal := c.newDeleteModal(tt, t)
				c.App.SetFocus(modal)
			}
			return nil
		case 'n':
			day := 24 * time.Hour
			due := time.Now().Add(day).Round(day)
			row, _ := table.GetSelection()
			ref := table.GetCell(row, 0).GetReference()
			orig := &apis.Task{State: apis.Open, Status: apis.Backlog, Due: &due}
			form := c.newTaskForm(tt, orig, "Create task", func(formData *taskFormData) error {
				t := &apis.Task{}
				err := formData.ApplyTo(t)
				if err != nil {
					return err
				}
				userId := fmt.Sprintf("%d", c.Me.ID)
				up, err := generic.Post(c.Client, "users/"+userId+"/tasks", t)
				if err != nil {
					return err
				}
				if ref != nil {
					prev := ref.(*TaskModel)
					prev.LastTouched = false
				}
				tt.Tasks = append(tt.Tasks, TaskModel{Task: up, LastTouched: true})
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

	return tt
}

func (t *TaskTable) SetTasks(tasks []apis.Task) {
	model := make([]TaskModel, len(tasks))
	for i := range tasks {
		model[i] = TaskModel{
			Task:        &tasks[i],
			LastTouched: false,
		}
	}
	t.Tasks = model
	t.Update(true)
}

func (t *TaskTable) Update(clear bool) {
	table := t.Table
	tasks := t.Tasks
	if clear {
		t.Clear()
	}

	// add table header
	for c, h := range taskTableHeaders {
		table.SetCell(0, c,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorViolet).
				SetSelectable(false).
				SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	// primary sort tasks by status and then secondary sort by due date and priority
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

	// add tasks to the table
	for r := range tasks {
		task := &tasks[r]
		r = r + 1

		priority := fmt.Sprintf("%d", task.Priority)
		createdAt := task.CreatedAt.Format(dateSpec)

		var due string
		if task.Due != nil {
			due = task.Due.Format(dateSpec)
		}

		// id := fmt.Sprintf("%d", task.ID)
		table.SetCell(r, 0, tview.NewTableCell(string(createdAt)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetReference(task))
		table.SetCell(r, 1, tview.NewTableCell(task.Description).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(4))
		table.SetCell(r, 2, tview.NewTableCell(due).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 3, tview.NewTableCell(priority).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 4, tview.NewTableCell(string(task.State)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 5, tview.NewTableCell(string(task.Status)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 6, tview.NewTableCell(privateMap[task.Private]).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))

		if task.LastTouched {
			table.Select(r, 0)
		}
	}
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

	newTasks := make([]TaskModel, 0)
	newTasks = append(newTasks, t.Tasks[:toRemove]...)
	newTasks = append(newTasks, t.Tasks[toRemove+1:]...)

	t.Tasks = newTasks
	return nil
}

var (
	dateSpec = "2006-01-02 15:04:05 MST"

	privateMap = map[bool]string{true: "✓", false: "✗"}

	statusOrder = map[apis.Status]int{
		apis.Doing:     0,
		apis.Todo:      1,
		apis.Backlog:   2,
		apis.Done:      3,
		apis.Abandoned: 4,
	}

	taskTableHeaders = []string{
		"Created",
		"Description",
		"Due",
		"Priority",
		"State",
		"Status",
		"Private",
	}

	taskTableKeyBindings = []KeyBinding{
		{"n", "Create a new Task"},
		{"<Enter>", "Edit the selected task"},
		{"o", "See tasks I own"},
		{"a", "See tasks assigned to me"},
		{"q", "Quit with prompt"},
		{"Q", "Quit Immediately"},
		{"?", "Show this help"},
	}
)
