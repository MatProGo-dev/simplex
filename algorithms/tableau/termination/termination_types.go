package tableau_termination

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
)

type TerminationType string

const DidNotTerminate TerminationType = "Did Not Terminate"
const MaximumIterationsReached TerminationType = "Maximum Iterations Reached"
const OptimalSolutionFound TerminationType = "Optimal Solution Found"

func (tt TerminationType) ToOptimizationStatus() problem.OptimizationStatus {
	switch tt {
	case DidNotTerminate:
		return problem.OptimizationStatus_INPROGRESS
	case MaximumIterationsReached:
		return problem.OptimizationStatus_ITERATION_LIMIT
	case OptimalSolutionFound:
		return problem.OptimizationStatus_OPTIMAL
	default:
		return problem.OptimizationStatus_INPROGRESS
	}
}
