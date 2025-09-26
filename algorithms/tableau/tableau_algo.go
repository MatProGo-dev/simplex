package tableau_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"gonum.org/v1/gonum/mat"
)

type TableauAlgorithm struct {
	ProblemInStandardForm *problem.OptimizationProblem
	CurrentSolution       *mat.VecDense
	IterationLimit        int
}

func (algo *TableauAlgorithm) Solve(initialState TableauAlgorithmState) (problem.Solution, error) {
	// Setup
	var stateII TableauAlgorithmState = initialState
	var status problem.OptimizationStatus = problem.OptimizationStatus_INPROGRESS

	// Loop
	for iter := 0; iter < algo.IterationLimit; iter++ {
		// Test for Termination
		terminated, err := stateII.CheckTerminationCondition()
		if err != nil {
			return problem.Solution{},
				fmt.Errorf(
					"There was an issue checking the termination condition at iteration %v: %v",
					iter,
					err,
				)
		}

		if terminated {
			break
		}

		// Compute XB, y and r
		_, err = stateII.XBasic()
		if err != nil {
			return problem.Solution{},
				fmt.Errorf(
					"There was an issue computing the value XBasic() at iteration #%v: %v",
					iter,
					err,
				)
		}

		// If we reach the limit, then return limit reached as status
		if iter == algo.IterationLimit-1 {
			status = problem.OptimizationStatus_ITERATION_LIMIT
		}

	}

	return stateII.ToSolution(status), nil
}
