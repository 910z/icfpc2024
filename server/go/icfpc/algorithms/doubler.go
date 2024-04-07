package algorithms

import (
	"icfpc/types"
)

type Doubler struct{}

func (r *Doubler) Version() string {
	return "1.0.0"
}

func (r *Doubler) Solve(task *types.Task) (solution *types.Solution, explanation *types.Explanation) {
	solution = &types.Solution{Data: task.Data * 2}
	return
}

var _ Algorithm = &Doubler{}
