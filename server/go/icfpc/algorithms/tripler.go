package algorithms

import (
	"context"
	"icfpc/database"
)

type Tripler struct{}

var _ IAlgorithm = Tripler{}

func (t Tripler) Version() string {
	return "1.0.0"
}

func (t Tripler) Solve(_ context.Context, task database.Task) (database.Solution, database.SolutionExplanation, error) {
	return database.Solution{Data: task.Data * 3}, nil, nil
}
