package solver_test

import (
	"strings"
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"matprogo.dev/solvers/simplex/solver"
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

/*
TestInitialTableau1
Description:

	In this test, we verify that the InitialTableau() function correctly creates a tableau from a linear
	problem in standard form.
*/
func TestInitialTableau1(t *testing.T) {
	// Setup
	problemIn := solver.GetTestProblem3()

	transformed_problem, _ := solver.ToStandardFormWithSlackVariables(problemIn)

	// Create the tableau
	tableau := solver.InitialTableau(transformed_problem)
	if tableau.Problem != transformed_problem {
		t.Errorf("Expected problem to be %v, but got %v", transformed_problem, tableau.Problem)
	}

	t.Errorf("Number of variables in tableau: %d", len(tableau.Problem.Variables))
	for i, variable := range tableau.Problem.Variables {
		t.Errorf("Variable %d: %v", i, variable)
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

}
