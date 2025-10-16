package tableau_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	tableau_termination "github.com/MatProGo-dev/simplex/algorithms/tableau/termination"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
	"github.com/MatProGo-dev/simplex/utils"
	"gonum.org/v1/gonum/mat"
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

func (algo *TableauAlgorithm) Solve(prob problem.OptimizationProblem) (simplex_solution.SimplexSolution, error) {
	// Setup

	// Create initial Tableau state from the problem
	initialTableau, mapFromOriginalVariablesToStandardFormVariables, err := utils.GetInitialTableauFrom(&prob)
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("there was an issue creating the initial tableau: %v", err)
	}
	stateII := TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Loop
	sol := simplex_solution.SimplexSolution{}
	for iter := 0; iter < algo.IterationLimit; iter++ {
		// Test for Termination
		condition, err := algo.CheckTerminationConditions(stateII)
		if err != nil {
			return simplex_solution.SimplexSolution{},
				fmt.Errorf(
					"There was an issue checking the termination condition at iteration %v: %v",
					iter,
					err,
				)
		}

		if condition != tableau_termination.DidNotTerminate {
			sol, err = stateII.ToSolution(condition, mapFromOriginalVariablesToStandardFormVariables, &prob)
			if err != nil {
				return simplex_solution.SimplexSolution{},
					fmt.Errorf(
						"There was an issue converting the final state to a solution at iteration %v: %v",
						iter,
						err,
					)
			}
			// Exit the loop
			break
		}

		fmt.Println("Iteration: ", iter)
		fmt.Println("Matrix: ", mat.Formatted(stateII.Tableau.AsCompressedMatrix))

		// Update the state
		stateII, err = stateII.CalculateNextState()
		if err != nil {
			return simplex_solution.SimplexSolution{},
				fmt.Errorf(
					"There was an issue updating the state at iteration %v: %v",
					iter,
					err,
				)
		}

	}

	return sol, nil
}
