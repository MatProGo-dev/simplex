package simplex_solution

import solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"

// SimplexSolution represents the result of solving a linear program using the simplex method.
// It contains the values of the decision variables, the objective value, and the solution status.
type SimplexSolution struct {
	// VariableValues maps variable IDs (as uint64) to their solution values.
	// The uint64 key typically represents the unique identifier or index of a variable in the model.
	VariableValues map[uint64]float64
	// Objective is the value of the objective function at the solution.
	Objective      float64
	// Status indicates the status of the solution (e.g., optimal, infeasible).
	Status         solution_status.SolutionStatus
}
