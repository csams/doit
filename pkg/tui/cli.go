package tui

import (
	"fmt"
	"strconv"

	"github.com/csams/doit/pkg/apis"
	generic "github.com/csams/doit/pkg/tui/client"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
	table.SetTitle("Tasks assigned to " + me.Username)
	c.Root.AddItem(table, 0, 1, true) // (item, fixedSize; 0 means not fixed, proportion, focus?)

	return c, nil
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

func (c *CLI) newTaskForm(table *TaskTable, task *apis.Task, title string, save func(*taskFormData) error) *tview.Form {
	formData := formDataFromTask(task)

	form := styledForm()
	form.SetTitle(title)

	form.AddInputField("Description", formData.Description, 0, nil, func(text string) { formData.Description = text })
	form.AddInputField("Due", formData.Due, 30, nil, func(text string) { formData.Due = text })
	form.AddDropDown("State", []string{"open", "closed"}, getStateIndex(formData.State), func(option string, index int) { formData.State = apis.State(option) })
	form.AddDropDown("Status", []string{"backlog", "todo", "doing", "done", "abandoned"}, getStatusIndex(formData.Status), func(option string, index int) { formData.Status = apis.Status(option) })
	form.AddInputField("Priority", strconv.Itoa(int(formData.Priority)), 3, EnsureInt, func(text string) {
		p, _ := strconv.Atoi(text)
		formData.Priority = apis.Priority(p)
	})
	form.AddCheckbox("Private", formData.Private, func(checked bool) { formData.Private = checked })

	doSave := func() {
		if err := save(formData); err != nil {
			c.Root.RemoveItem(form)
			c.newErrorModal("Error saving task: " + err.Error())
		} else {
			c.Root.RemoveItem(form)
			c.App.SetFocus(table.Table)
			table.Update(false)
		}
	}

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			doSave()
			return nil
		}
		return event
	})

	cancel := func() {
		c.Root.RemoveItem(form)
		c.App.SetFocus(table.Table)
	}
	form.SetCancelFunc(cancel)
	form.AddButton("Cancel", cancel)
	form.AddButton("Save", doSave)

	c.Root.SetDirection(tview.FlexRow).AddItem(form, 0, 1, true)
	return form
}

func (c *CLI) newDeleteModal(table *TaskTable, orig *TaskModel) *tview.Modal {
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
			if err = table.Remove(orig.Task); err != nil {
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

var (
	statusMap = map[apis.Status]int{
		apis.Backlog:   0,
		apis.Todo:      1,
		apis.Doing:     2,
		apis.Done:      3,
		apis.Abandoned: 4,
	}
	stateMap = map[apis.State]int{
		apis.Open:   0,
		apis.Closed: 1,
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
