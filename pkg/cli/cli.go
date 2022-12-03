package cli

import (
	"fmt"

	"encoding/json"
	"io"
	"net/http"

	"github.com/csams/doit/pkg/apis"
	"github.com/csams/doit/pkg/auth"
	"github.com/gdamore/tcell/v2"
	"github.com/go-logr/logr"
	"github.com/rivo/tview"
)

/*
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

var (
	privateMap = map[bool]string{true: "true", false: "false"}
)

func GetApplication(log logr.Logger, app *tview.Application, tokenProvider *auth.TokenProvider) (tview.Primitive, error) {
	flex := tview.NewFlex().SetFullScreen(true)
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	table.SetBorder(true)
	table.SetSeparator(tview.Borders.Vertical)

	flex.AddItem(table, 0, 1, true)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	client := auth.CreateClient(true)
	req, err := http.NewRequest("GET", "http://localhost:8080/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "todo-app-client")

	token, err := tokenProvider.GetIdToken()
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("BEARER %s", token)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user apis.User

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}

	table.SetTitle("Tasks for " + user.Name)

	headers := []string{
		"id",
		"created_at",
		"description",
		"due",
		"priority",
		"state",
		"status",
		"private",
	}

	for c, h := range headers {
		table.SetCell(0, c,
			tview.NewTableCell(h).
				SetTextColor(tcell.ColorViolet).
				SetSelectable(false).
				SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	tasks := user.AssignedTasks

	for r, task := range tasks {
		r = r + 1
		id := fmt.Sprintf("%d", task.ID)
		priority := fmt.Sprintf("%d", task.Priority)
		createdAt, _ := task.CreatedAt.MarshalText()

		var due []byte
		if task.Due != nil {
			due, _ = task.Due.MarshalText()
		}

		table.SetCell(r, 0, tview.NewTableCell(id).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignLeft).SetReference(task))
		table.SetCell(r, 1, tview.NewTableCell(string(createdAt)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 2, tview.NewTableCell(task.Description).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft).SetExpansion(4))
		table.SetCell(r, 3, tview.NewTableCell(string(due)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 4, tview.NewTableCell(priority).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 5, tview.NewTableCell(string(task.State)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 6, tview.NewTableCell(string(task.Status)).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
		table.SetCell(r, 7, tview.NewTableCell(privateMap[task.Private]).SetTextColor(tcell.ColorWheat).SetAlign(tview.AlignLeft))
	}
	table.Select(1, 1)
	return flex, nil
}
