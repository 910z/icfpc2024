package algorithms

import (
	"context"
	"icfpc/database"
)

type (
	Solution struct {
		Data int
	}

	Explanation any
)

type IAlgorithm interface {
	Version() string
	Solve(context.Context, database.Task) (Solution, Explanation, error)
}
