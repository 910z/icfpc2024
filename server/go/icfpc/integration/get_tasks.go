package integration

import (
	"context"
	"errors"
	"icfpc/database"
	"log/slog"
)

var allTasks = []database.Task{
	{ExternalID: "icfpc_id_1", Data: 5},
	{ExternalID: "icfpc_id_2", Data: 10},
}

var Error = errors.New("integration error")

func GetTasks() ([]database.Task, error) {
	return allTasks, nil
}

func SendSolution(ctx context.Context, taskId string, solution database.Solution) error {
	slog.InfoContext(ctx, "sent best", slog.String("taskId", taskId))
	return nil
}
