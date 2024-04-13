package evaluation

import (
	"context"
	"icfpc/database"
)

var Version = "1.0.0"

func EvaluateSolution(
	_ context.Context,
	task database.Task,
	solution database.Solution,
) (database.EvalResult, error) {
	return database.EvalResult{
		Score:       solution.Data,
		Explanation: "just because",
	}, nil
}
