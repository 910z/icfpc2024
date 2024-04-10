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

	ProgressStatus string

	RunResult struct {
		ID               int64               `bun:"id,pk,autoincrement"`
		TaskID           int64               `bun:"unique:task_algorithm"`
		Task             Task                `bun:"rel:has-one,join:task_id=id"`
		AlgorithmName    string              `bun:"unique:task_algorithm"`
		AlgorithmVersion string              `bun:"unique:task_algorithm"`
		Explanation      SolutionExplanation `bun:"type:jsonb"`
		Solution         Solution            `bun:"type:jsonb"`
		Status           ProgressStatus
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
		ID          int64 `bun:"id,pk,autoincrement"` // суррогатный айдишник
		Error       string
		StartedAt   time.Time
		FinishedAt  time.Time
		Status      ProgressStatus
		RunResultID int64     `bun:"unique:version_result_id"`
		RunResult   RunResult `bun:"rel:belongs-to,join:run_result_id=id"`
		Version     string    `bun:"unique:version_result_id"`
	}
)

const (
	ProgressStatusStarted  ProgressStatus = "started"
	ProgressStatusFinished ProgressStatus = "finished"
	ProgressStatusError    ProgressStatus = "error"
)

var allModels = []any{
	(*Task)(nil),
	(*RunResult)(nil),
	(*RunEvalResult)(nil),
}
