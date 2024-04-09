package runner

import (
	"icfpc/algorithms"
)

var (
	AllAlgorithms = []algorithms.IAlgorithm{
		&algorithms.Doubler{},
		&algorithms.Tripler{},
	}
)
