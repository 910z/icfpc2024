package integration

import "icfpc/database"

var allTasks = []database.Task{
	{ExternalID: "1", Data: 5},
	{ExternalID: "2", Data: 10},
}

type Error struct {
	error
}

func GetTasks() ([]database.Task, error) {
	return allTasks, nil
}
