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

func SendSolution(ctx context.Context, taskId string, solution database.Solution) error {
	slog.InfoContext(ctx, "sent best to API", slog.String("taskId", taskId), slog.Any("solution", solution))
	return nil
}
