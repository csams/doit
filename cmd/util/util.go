package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/csams/doit/pkg/apis/task"
	"github.com/spf13/pflag"
)

func GetDue(flags *pflag.FlagSet) (*time.Time, error) {
	dueStr, err := flags.GetString("due")
	if dueStr == "" || err != nil {
		return nil, err
	}

	due, err := dateparse.ParseStrict(dueStr)
	if err != nil {
		return nil, err
	}

	return &due, nil
}

func GetPriority(flags *pflag.FlagSet) (task.Priority, error) {
	prioStr, err := flags.GetString("priority")
	if prioStr == "" || err != nil {
		return task.Undefined, err
	}

	switch strings.ToLower(prioStr) {
	case "high", "h":
		return task.High, nil
	case "medium", "med", "m":
		return task.Medium, nil
	case "low", "l":
		return task.Low, nil
	default:
		return task.Undefined, fmt.Errorf("unrecognized priority: %s. use H, M, or L", prioStr)
	}
}

func GetStatus(flags *pflag.FlagSet) (*task.Status, error) {
	statusStr, err := flags.GetString("status")
	if statusStr == "" || err != nil {
		return nil, err
	}

	status := task.Status(strings.ToLower(statusStr))
	if !task.IsValidStatus(status) {
		statuses := task.StatusStrings()
		values := strings.Join(statuses, ", ")
		return nil, fmt.Errorf("unrecognized status: %s. valid values are %s", status, values)
	}

	return &status, nil
}

func GetTags(flags *pflag.FlagSet) ([]string, error) {
	tags, err := flags.GetStringSlice("tags")
	if tags == nil || err != nil {
		return nil, err
	}

	cleanTags := make([]string, 0, len(tags))
	for _, t := range tags {
		cleanTags = append(cleanTags, strings.TrimSpace(t))
	}

	return cleanTags, nil
}
