package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/csams/doit/pkg/apis"
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

func GetPriority(flags *pflag.FlagSet) (apis.Priority, error) {
	prioStr, err := flags.GetString("priority")
	if prioStr == "" || err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(prioStr)

	if err != nil {
		return 0, fmt.Errorf("unrecognized priority: %s. use H, M, or L", prioStr)
	}
	return apis.Priority(v), nil
}

func GetStatus(flags *pflag.FlagSet) (*apis.Status, error) {
	statusStr, err := flags.GetString("status")
	if statusStr == "" || err != nil {
		return nil, err
	}

	status := apis.Status(strings.ToLower(statusStr))
	if !apis.IsValidStatus(status) {
		statuses := apis.StatusStrings()
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
