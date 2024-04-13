package algorithms

import (
	"context"
	"icfpc/database"
)

type Doubler struct{}

var _ IAlgorithm = Doubler{}

func (d Doubler) Version() string {
	return "1.0.0"
}

func (d Doubler) Solve(_ context.Context, task database.Task) (database.Solution, database.SolutionExplanation, error) {
	return database.Solution{Data: task.Data * 2}, nil, nil
}
