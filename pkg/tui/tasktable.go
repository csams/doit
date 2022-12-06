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
		task := table.GetCell(row, 0).GetReference().(*apis.Task)

		form := c.newTaskForm(tt, task, "Edit task", func(formData *taskFormData) error {
			proposedTask := *task
			err := formData.ApplyTo(&proposedTask)
			if err != nil {
				return err
			}
			userId := fmt.Sprintf("%d", c.Me.ID)
			taskId := fmt.Sprintf("%d", task.ID)
			up, err := generic.Put(c.Client, "users/"+userId+"/tasks/"+taskId, &proposedTask)
			if err != nil {
				return err
			}
			*task = *up
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

	r, _ := t.Table.GetSelection()
	ref := t.Table.GetCell(r, 0).GetReference()
	var focusedTask *apis.Task = nil
	if ref != nil {
		focusedTask = ref.(*apis.Task)
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
		createdAt := task.CreatedAt.Format(dateSpec)

		var due string
		if task.Due != nil {
			due = task.Due.Format(dateSpec)
		}

		table.SetCell(r, 0, tview.NewTableCell(id).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetReference(task))
		table.SetCell(r, 1, tview.NewTableCell(string(createdAt)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 2, tview.NewTableCell(task.Description).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(4))
		table.SetCell(r, 3, tview.NewTableCell(due).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 4, tview.NewTableCell(priority).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 5, tview.NewTableCell(string(task.State)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 6, tview.NewTableCell(string(task.Status)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 7, tview.NewTableCell(privateMap[task.Private]).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))

		if focusedTask != nil {
			if task.ID == focusedTask.ID {
				table.Select(r, 1)
			}
		}
	}
	if focusedTask == nil {
		table.Select(1, 1)
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

	newTasks := make([]apis.Task, 0)
	newTasks = append(newTasks, t.Tasks[:toRemove]...)
	newTasks = append(newTasks, t.Tasks[toRemove+1:]...)

	t.Tasks = newTasks
	return nil
}

var (
	dateSpec    = "2006-01-02 15:04:05 MST"
	statusOrder = map[apis.Status]int{
		apis.Doing:     0,
		apis.Todo:      1,
		apis.Backlog:   2,
		apis.Done:      3,
		apis.Abandoned: 4,
	}
)
