package solver_test

import (
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	simplexSolver "matprogo.dev/solvers/simplex/solver"
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

	// Create function to catch a panic

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Expected panic, but got none")
		}
	}()

	// Expect a panic when calling ToTableau with a non-linear problem
	simplexSolver.ToTableau(nonLinearProblem)

	t.Errorf("Expected panic, but got none")
}
