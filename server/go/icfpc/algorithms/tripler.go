package algorithms

import (
	"icfpc/types"
)

type Tripler struct{}

func (r *Tripler) Version() string {
	return "1.0.0"
}

func (r *Tripler) Solve(task *types.Task) (solution *types.Solution, explanation *types.Explanation) {
	solution = &types.Solution{Data: task.Data * 3}
	return
}

var _ Algorithm = &Tripler{}
