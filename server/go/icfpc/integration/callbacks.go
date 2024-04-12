package integration

import (
	"context"
	"errors"
	"icfpc/database"
	"log/slog"
)

var Error = errors.New("integration error")

var allTasks = []database.Task{
	{ExternalID: "icfpc_id_1", Data: 5},
	{ExternalID: "icfpc_id_2", Data: 10},
}

func GetTasks() ([]database.Task, error) {
	return allTasks, nil
}

func SendSolution(ctx context.Context, taskId string, solution database.Solution) (string, error) {
	slog.InfoContext(ctx, "sent best to API", slog.String("taskId", taskId), slog.Any("solution", solution))
	return taskId, nil
}

type CheckedSubmission struct {
	Status database.SubmissionStatus
	Score  database.Score
}

type CheckedSubmissions = map[string]CheckedSubmission

func GetSubmissionsStatus(_ context.Context, runResults []database.RunResult) (CheckedSubmissions, error) {
	m := map[string]CheckedSubmission{}
	for _, task := range allTasks {
		for _, res := range runResults {
			if res.Task.ExternalID == task.ExternalID {
				m[task.ExternalID] = CheckedSubmission{
					Status: database.SubmissionStatusChecked,
					Score:  res.Solution.Data - 13,
				}
			}
		}
	}
	return m, nil
}
