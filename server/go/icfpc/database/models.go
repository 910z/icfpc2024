package database

import (
	"time"
)

type (
	Score = int64

	Solution struct {
		Data int64
	}

	SolutionExplanation any
	EvalExplanation     any

	ProgressStatus   string
	SubmissionStatus string

	Submission struct {
		SubmissionStatus    SubmissionStatus
		SubmittedAt         time.Time
		SubmissionCheckedAt time.Time
		SubmissionError     string
		SubmissionToken     string
		ExternalScore       Score
	}

	RunResult struct {
		Submission
		ID                  int64               `bun:"id,pk,autoincrement"`
		TaskID              int64               `bun:"unique:task_algorithm"`
		Task                Task                `bun:"rel:has-one,join:task_id=id"`
		AlgorithmName       string              `bun:"unique:task_algorithm"`
		AlgorithmVersion    string              `bun:"unique:task_algorithm"`
		Explanation         SolutionExplanation `bun:"type:jsonb"`
		Solution            Solution            `bun:"type:jsonb"`
		AlgorithmStatus     ProgressStatus
		AlgorithmStartedAt  time.Time
		AlgorithmFinishedAt time.Time
		Error               string
	}

	TaskData = int64

	Task struct {
		ID         int64    `bun:"id,pk,autoincrement"`  // суррогатный айдишник
		ExternalID string   `bun:"unique:external_data"` // то, как она называется на сервере ICFPC
		Data       TaskData `bun:"unique:external_data"`
	}

	EvalResult struct {
		Score       Score
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

	SubmissionStatusPending      SubmissionStatus = "pending"
	SubmissionStatusChecked      SubmissionStatus = "checked" // в том числе получена другая оценка
	SubmissionStatusError        SubmissionStatus = "error"
	SubmissionStatusNotSubmitted SubmissionStatus = ""
)

var allModels = []any{
	(*Task)(nil),
	(*RunResult)(nil),
	(*RunEvalResult)(nil),
}
