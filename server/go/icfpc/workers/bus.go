package workers

import (
	"icfpc/database"
)

type bus struct {
	tasksAdded        chan struct{}
	algorithmFinish   chan struct{}
	solutionEvaluated chan struct{}
}

func NewBus() bus {
	return bus{
		tasksAdded:        make(chan struct{}, 1),
		algorithmFinish:   make(chan struct{}, 1),
		solutionEvaluated: make(chan struct{}, 1),
	}
}

func (b bus) onTasksAdded([]database.Task) {
	select {
	case b.tasksAdded <- struct{}{}:
	default:
	}
}

func (b bus) onAlgorithmFinish(database.RunResult) {
	select {
	case b.algorithmFinish <- struct{}{}:
	default:
	}
}

func (b bus) onSolutionEvaluated(database.RunEvalResult) {
	select {
	case b.solutionEvaluated <- struct{}{}:
	default:
	}
}
