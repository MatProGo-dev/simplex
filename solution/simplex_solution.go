package simplex_solution

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/MatProInterface.go/solution"
	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
)

// SimplexSolution represents the result of solving a linear program using the simplex method.
// It contains the values of the decision variables, the objective value, and the solution status.
type SimplexSolution struct {
	// VariableValues maps variable IDs (as uint64) to their solution values.
	// The uint64 key typically represents the unique identifier or index of a variable in the model.
	VariableValues map[uint64]float64
	// Status indicates the status of the solution (e.g., optimal, infeasible).
	Status     solution_status.SolutionStatus
	Iterations int
	// originalProblem is the original optimization problem that was solved to obtain this solution.
	// It is included for reference and may be nil if not applicable.
	OriginalProblem *problem.OptimizationProblem
}

func (sol *SimplexSolution) GetValueMap() map[uint64]float64 {
	return sol.VariableValues
}

func (sol *SimplexSolution) GetOptimalValue() float64 {
	// Use the symbolic.Solution interface to compute the optimal value
	optVal, err := solution.GetOptimalObjectiveValue(sol)
	if err != nil {
		return 0.0
	}
	return optVal
}

func (sol *SimplexSolution) GetStatus() solution_status.SolutionStatus {
	return sol.Status
}

func (sol *SimplexSolution) GetProblem() *problem.OptimizationProblem {
	return sol.OriginalProblem
}
