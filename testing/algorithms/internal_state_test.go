package algorithms_test

import (
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms"
	"gonum.org/v1/gonum/mat"
)

/*
TestToTableau1
Description:

	In this test, we verify that the ToTableau() function correctly panics when the input problem is not
	a linear program.
*/
func TestToTableau1(t *testing.T) {
	// Create a non-linear optimization problem
	nonLinearProblem := problem.NewProblem("TestToTableau1 Problem")

	// Create non-linear objective
	x := symbolic.NewVariableVector(2)
	nonLinearProblem.SetObjective(
		x.Transpose().Multiply(x),
		problem.SenseMinimize,
	)

	// Create an internal state object
	internalState := &algorithms.AlgorithmInternalState{
		BasicVariables:    []symbolic.Variable{symbolic.NewVariable(), symbolic.NewVariable()},
		NonBasicVariables: []symbolic.Variable{symbolic.NewVariable()},
		NonBasicValues:    mat.NewVecDense(2, nil),
		IterationCount:    0,
	}

	// Expect a panic when calling ToTableau with a non-linear problem
	_, err := internalState.ToTableau(nonLinearProblem)
	if err == nil {
		// If no error occurs, fail the test
		t.Errorf("Expected an error, but got none")
	}

}
