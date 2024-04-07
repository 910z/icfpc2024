package algorithms

import (
	"icfpc/database"
)

type Doubler struct{}

var _ IAlgorithm = Doubler{}

func (d Doubler) Version() string {
	return "1.0.0"
}

func (d Doubler) Solve(task database.Task) (Solution, Explanation, error) {
	return Solution{Data: task.Data * 2}, nil, nil
}
