package types

import "time"

type EvalResult struct {
	Score       int
	Explanation string
}
type RunEvalResult struct {
	EvalResult
	RunResult *RunResult `bun:"rel:belongs-to,join:run_result_id=id"`
}

type Explanation interface{}

type RunResult struct {
	Id               int64       `bun:"id,pk,autoincrement"`
	TaskId           string      `bun:"unique:task_algorithm"`
	AlgorithmName    string      `bun:"unique:task_algorithm"`
	AlgorithmVersion string      `bun:"unique:task_algorithm"`
	Explanation      Explanation `bun:"type:jsonb"`
	Solution         *Solution   `bun:"type:jsonb"`
	StartedAt        time.Time
	FinishedAt       time.Time
	RunStatus        RunStatus
}

type RunStatus int

const (
	RunStatusStarted RunStatus = iota
	RunStatusFinished
	RunStatusError
)
