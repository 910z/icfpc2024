package evaluation

import (
	"icfpc/algorithms"
	"icfpc/database"
)

var Version = "1.0.0"

func EvaluateSolution(task database.Task, solution algorithms.Solution) database.EvalResult {
	return database.EvalResult{
		Score:       solution.Data,
		Explanation: "just because",
	}
}
