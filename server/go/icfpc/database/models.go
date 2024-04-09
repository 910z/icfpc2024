package database

import (
	"time"
)

type (
	Solution struct {
		Data int64
	}

	SolutionExplanation any
	EvalExplanation     any

	RunStatus string

	RunResult struct {
		ID               int64               `bun:"id,pk,autoincrement"`
		TaskID           int64               `bun:"unique:task_algorithm"`
		AlgorithmName    string              `bun:"unique:task_algorithm"`
		AlgorithmVersion string              `bun:"unique:task_algorithm"`
		Explanation      SolutionExplanation `bun:"type:jsonb"`
		Solution         Solution            `bun:"type:jsonb"`
		Status           RunStatus
		StartedAt        time.Time
		FinishedAt       time.Time
		Error            string
	}

	TaskData = int64

	Task struct {
		ID         int64    `bun:"id,pk,autoincrement"`  // суррогатный айдишник
		ExternalID string   `bun:"unique:external_data"` // то, как она называется на сервере ICFPC
		Data       TaskData `bun:"unique:external_data"`
	}

	EvalResult struct {
		Score       int64
		Explanation EvalExplanation `bun:"type:jsonb"`
	}

	RunEvalResult struct {
		EvalResult
		RunResultID int64
		RunResult   RunResult `bun:"rel:belongs-to,join:run_result_id=id"`
	}
)

const (
	RunStatusStarted  RunStatus = "started"
	RunStatusFinished RunStatus = "finished"
	RunStatusError    RunStatus = "error"
)

var allModels = []any{
	(*Task)(nil),
	(*RunResult)(nil),
	(*RunEvalResult)(nil),
}
