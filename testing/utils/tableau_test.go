package utils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	simplexSolver "github.com/MatProGo-dev/simplex/simplexSolver"
	"github.com/MatProGo-dev/simplex/utils"
	"gonum.org/v1/gonum/mat"
)

/*
TestInitialTableau1
Description:

	In this test, we verify that the InitialTableau() function correctly creates a tableau from a linear
	problem in standard form.
*/
func TestGetInitialTableau1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem3()

	// Create the tableau using the initial state + problem in standard form
	tableau, err := utils.GetInitialTableauFrom(problemIn)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	fmt.Println("Basic Variables: ", tableau.BasicVariableIndicies)
	fmt.Println("Nonbasic variables: ", tableau.NonBasicVariableIndicies())

	// Check that the number of basic variables is 3
	if len(tableau.BasicVariables()) != 3 {
		t.Errorf("Expected 3 basic variables, but got %d", len(tableau.BasicVariables()))
	}

	// Check that all basic variables are slack variables
	for _, basicVar := range tableau.BasicVariables() {
		if !strings.Contains(basicVar.Name, "slack") {
			t.Errorf("Expected basic variable to be a slack variable, but got \"%v\"", basicVar)
		}
	}

	// Check that the number of non-basic variables is 3
	if len(tableau.NonBasicVariables()) != len(tableau.Variables)-len(tableau.BasicVariables()) {
		t.Errorf(
			"Expected %v non-basic variables, but got %d",
			len(tableau.Variables)-len(tableau.BasicVariables()),
			len(tableau.NonBasicVariables()),
		)
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
	tableau, err := utils.GetInitialTableauFrom(problemIn)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Compute the feasible solution
	zerosVector := symbolic.ZerosVector(3)
	solution, err := tableau.ComputeFeasibleSolution(&zerosVector)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the solution is correct
	expectedSolution := mat.NewVecDense(3, []float64{4, 6, 8})
	if !mat.EqualApprox(solution, expectedSolution, 1e-10) {
		t.Errorf("Expected feasible solution to be %v, but got %v", expectedSolution, solution)
	}
}

// /*
// TestAsDense1
// Description:

// 	In this test, we verify that the AsDense() function correctly
// 	computes a feasible solution for the basic variables.

// 	We use an example problem from this youtube video:
// 		https://www.youtube.com/watch?v=QAR8zthQypc&t=483s
// 	which is written in standard form as:
// 		maximize 4 x1 + 3 x2 + 5 x3
// 		subject to
// 			x1 + 2 x2 + 2 x3 + slack1 = 4
// 			3 x1 + 4 x3 + slack2 = 12
// 			2 x1 + x2 + 4 x3 + slack3 = 8

// 		x1, x2, x3 >= 0
// 		slack1, slack2, slack3 >= 0

// 	We expect the feasible solution to be:
// 		[4, 12, 8]
// 	for the initial tableau.
// */
// func TestAsDense1(t *testing.T) {
// 	// Setup
// 	problemIn := simplexSolver.GetTestProblem3()

// 	// Create the tableau
// 	tableau, err := utils.GetInitialTableauFrom(problemIn)
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}

// 	// Compute the tableau as a dense matrix
// 	denseTableau, err := tableau.AsDense()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}
// 	nRows, _ := denseTableau.Dims()

// 	// Check that the dense tableau has zeros in the correct places
// 	for ii := 1; ii < nRows; ii++ {
// 		if denseTableau.At(ii, 0) != 0.0 {
// 			t.Errorf("Expected zero in row %d, column 0, but got %f", ii, denseTableau.At(ii, 0))
// 		}
// 	}

// }
