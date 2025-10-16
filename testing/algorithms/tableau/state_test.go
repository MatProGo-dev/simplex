package tableau

import (
	"testing"

	tableau_algorithm1 "github.com/MatProGo-dev/simplex/algorithms/tableau"
	"github.com/MatProGo-dev/simplex/utils"
	"github.com/MatProGo-dev/simplex/utils/examples"
)

/*
TestTableau_CalculateOptimalValue1
Description:

	In this test, we verify that the CalculateOptimalValue() function correctly computes the optimal value
	of the objective function for a known tableau.

	We use an example problem from this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
	The optimal solution of the objective function is:
		x1 = 125
		x2 = 300
		s1 = 25
		s2 = 0
		s3 = 0
		s4 = 225
*/
func TestTableau_CalculateOptimalSolution1(t *testing.T) {
	// Setup

	// Create the test problem
	testProblem := examples.GetTestProblem5()

	// Use the optimization solver to solve the problem

	// Create initial Tableau state from the problem
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}
	state0 := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Update state two times
	state1, err := state0.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}
	state2, err := state1.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}

	// Calculate the optimal solution
	solVec, err := state2.CalculateOptimalSolution()
	if err != nil {
		t.Errorf("there was an issue calculating the optimal value: %v", err)
	}

	// Check that the solution is correct
	expectedSol := []float64{125.0, 300.0, 25.0, 0.0, 0.0, 225.0}
	for ii, val := range expectedSol {
		if solVec.AtVec(ii) != val {
			t.Errorf("Expected solution value %v at index %d, but got %v", val, ii, solVec.AtVec(ii))
		}
	}
}
