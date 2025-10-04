package tableau_termination

import (
	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
)

type TerminationType string

const DidNotTerminate TerminationType = "Did Not Terminate"
const MaximumIterationsReached TerminationType = "Maximum Iterations Reached"
const OptimalSolutionFound TerminationType = "Optimal Solution Found"

func (tt TerminationType) ToOptimizationStatus() solution_status.SolutionStatus {
	switch tt {
	case DidNotTerminate:
		return solution_status.INPROGRESS
	case MaximumIterationsReached:
		return solution_status.ITERATION_LIMIT
	case OptimalSolutionFound:
		return solution_status.OPTIMAL
	default:
		return solution_status.INPROGRESS
	}
}
