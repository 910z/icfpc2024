package database

import (
	"time"
)

type (
	Solution struct {
		Data int
	}

	SolutionExplanation any
	EvalExplanation     any

	RunStatus string

	RunResult struct {
		ID               int64               `bun:"id,pk,autoincrement"`
		TaskID           string              `bun:"unique:task_algorithm"`
		AlgorithmName    string              `bun:"unique:task_algorithm"`
		AlgorithmVersion string              `bun:"unique:task_algorithm"`
		Explanation      SolutionExplanation `bun:"type:jsonb"`
		Solution         Solution            `bun:"type:jsonb"`
		Status           RunStatus
		StartedAt        time.Time
		FinishedAt       time.Time
		Error            string
	}

	Task struct {
		ID          string
		Description string
		Data        int
	}

	EvalResult struct {
		Score       int
		Explanation EvalExplanation `bun:"type:jsonb"`
	}

	RunEvalResult struct {
		EvalResult
		RunResult RunResult `bun:"rel:belongs-to,join:run_result_id=id"`
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
