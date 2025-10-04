package simplex_solution

import solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"

type SimplexSolution struct {
	VariableValues map[uint64]float64
	Objective      float64
	Status         solution_status.SolutionStatus
}
