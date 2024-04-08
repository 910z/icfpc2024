package database

import (
	"time"
)

type (
	RunStatus string

	RunResult struct {
		ID               int64  `bun:"id,pk,autoincrement"`
		TaskID           string `bun:"unique:task_algorithm"`
		AlgorithmName    string `bun:"unique:task_algorithm"`
		AlgorithmVersion string `bun:"unique:task_algorithm"`
		Explanation      any    `bun:"type:jsonb"`
		Error            string
		Solution         any `bun:"type:jsonb"`
		Status           RunStatus
		StartedAt        time.Time
		FinishedAt       time.Time
	}

	Task struct {
		ID          string
		Description string
		Data        int
	}
)

const (
	RunStatusStarted  RunStatus = "started"
	RunStatusFinished RunStatus = "finished"
	RunStatusError    RunStatus = "error"
)

var allModels = []any{
	(*RunResult)(nil),
	(*Task)(nil),
}
