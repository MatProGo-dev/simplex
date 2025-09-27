package tableau_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	tableau_termination "github.com/MatProGo-dev/simplex/algorithms/tableau/termination"
	"github.com/MatProGo-dev/simplex/utils"
)

type TableauAlgorithm struct {
	IterationLimit int
}

func (algo *TableauAlgorithm) CheckTerminationConditions(state TableauAlgorithmState) (tableau_termination.TerminationType, error) {
	// Input Checking
	err := state.Check()
	if err != nil {
		return tableau_termination.DidNotTerminate, err
	}

	// Check If the iteration limit has been reached
	if state.IterationCount >= algo.IterationLimit {
		return tableau_termination.MaximumIterationsReached, nil
	}

	// Check that the reduced costs are all non-negative
	if state.Tableau.CanNotBeImproved() {
		return tableau_termination.OptimalSolutionFound, nil
	}

	return tableau_termination.DidNotTerminate, nil
}

func (algo *TableauAlgorithm) Solve(prob problem.OptimizationProblem) (problem.Solution, error) {
	// Setup

	// Create initial Tableau state from the problem
	initialTableau, err := utils.GetInitialTableauFrom(&prob)
	if err != nil {
		return problem.Solution{}, fmt.Errorf("there was an issue creating the initial tableau: %v", err)
	}
	stateII := TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Loop
	sol := problem.Solution{}
	for iter := 0; iter < algo.IterationLimit; iter++ {
		// Test for Termination
		condition, err := algo.CheckTerminationConditions(stateII)
		if err != nil {
			return problem.Solution{},
				fmt.Errorf(
					"There was an issue checking the termination condition at iteration %v: %v",
					iter,
					err,
				)
		}

		if condition != tableau_termination.DidNotTerminate {
			sol.Status = condition.ToOptimizationStatus()
			break
		}

		// Update the state
		stateII, err = stateII.CalculateNextState()
		if err != nil {
			return problem.Solution{},
				fmt.Errorf(
					"There was an issue updating the state at iteration %v: %v",
					iter,
					err,
				)
		}

	}

	return sol, nil
}
