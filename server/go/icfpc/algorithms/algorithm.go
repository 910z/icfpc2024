package algorithms

import (
	"icfpc/types"
)

type Algorithm interface {
	Version() string
	Solve(*types.Task) (*types.Solution, *types.Explanation)
}
