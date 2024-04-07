package runner

import (
	"icfpc/algorithms"
	"icfpc/database"
)

var (
	AllTasks = []database.Task{
		{ID: "1", Description: "Task1", Data: 5},
		{ID: "2", Description: "Task2", Data: 10},
	}

	AllAlgorithms = []algorithms.IAlgorithm{
		&algorithms.Doubler{},
		&algorithms.Tripler{},
	}
)
