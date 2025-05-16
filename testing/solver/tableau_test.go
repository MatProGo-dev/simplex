package solver_test

import (
	"strings"
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	simplexSolver "matprogo.dev/solvers/simplex/simplexSolver"
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

/*
TestInitialTableau1
Description:

	In this test, we verify that the InitialTableau() function correctly creates a tableau from a linear
	problem in standard form.
*/
func TestGetInitialTableau1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem3()

	// Create the tableau
	tableau, err := simplexSolver.GetInitialTableau(problemIn)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if tableau.Problem == problemIn {
		t.Errorf("Expected problem to have changed from %v, but got %v", problemIn, tableau.Problem)
	}

	// Check that the number of basic variables is 3
	if len(tableau.BasicVariables) != 3 {
		t.Errorf("Expected 3 basic variables, but got %d", len(tableau.BasicVariables))
	}

	// Check that all basic variables are slack variables
	for _, basicVar := range tableau.BasicVariables {
		if !strings.Contains(basicVar.Name, "slack") {
			t.Errorf("Expected basic variable to be a slack variable, but got \"%v\"", basicVar)
		}
	}

	// Check that the number of non-basic variables is 3
	if len(tableau.NonBasicVariables) != len(tableau.Problem.Variables)-len(tableau.BasicVariables) {
		t.Errorf("Expected 3 non-basic variables, but got %d", len(tableau.NonBasicVariables))
	}

}

/*
TestComputeFeasibleSolution1
Description:

	In this test, we verify that the ComputeFeasibleSolution() function correctly computes a feasible
	solution for the basic variables.

	We use an example problem from this youtube video:
		https://www.youtube.com/watch?v=QAR8zthQypc&t=483s
	which is written in standard form as:
		maximize 4 x1 + 3 x2 + 5 x3
		subject to
			x1 + 2 x2 + 2 x3 + slack1 = 4
			3 x1 + 4 x3 + slack2 = 12
			2 x1 + x2 + 4 x3 + slack3 = 8

		x1, x2, x3 >= 0
		slack1, slack2, slack3 >= 0

	We expect the feasible solution to be:
		[4, 12, 8]
	for the initial tableau.
*/
func TestComputeFeasibleSolution1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem3()

	// Create the tableau
	tableau, err := simplexSolver.GetInitialTableau(problemIn)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Compute the feasible solution
	solution, err := tableau.ComputeFeasibleSolution()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the solution is correct
	expectedSolution := mat.NewVecDense(3, []float64{4, 12, 8})
	if !mat.EqualApprox(solution, expectedSolution, 1e-10) {
		t.Errorf("Expected feasible solution to be %v, but got %v", expectedSolution, solution)
	}
}
