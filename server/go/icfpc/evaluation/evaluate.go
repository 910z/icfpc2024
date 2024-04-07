package evaluation

import (
	"icfpc/types"
)

var Version = "1.0.0"

func EvaluateSolution(task *types.Task, solution *types.Solution) types.EvalResult {
	return types.EvalResult{Score: solution.Data, Explanation: "just because"}
}
