package algorithms

import (
	"context"
	"icfpc/database"
)

type ()

type IAlgorithm interface {
	Version() string
	Solve(context.Context, database.Task) (database.Solution, database.SolutionExplanation, error)
}
