package integration

import (
	"errors"
	"icfpc/database"
)

var allTasks = []database.Task{
	{ExternalID: "icfpc_id_1", Data: 5},
	{ExternalID: "icfpc_id_1", Data: 10},
}

var Error = errors.New("integration error")

func GetTasks() ([]database.Task, error) {
	return allTasks, nil
}
