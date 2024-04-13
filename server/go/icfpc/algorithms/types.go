package algorithms

import (
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
	Solve(database.Task) (Solution, Explanation, error)
}
