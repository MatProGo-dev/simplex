package utils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms/tableau/selection"
	"github.com/MatProGo-dev/simplex/utils"
	"github.com/MatProGo-dev/simplex/utils/examples"
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
	problemIn := examples.GetTestProblem3()

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
	problemIn := examples.GetTestProblem3()

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

/*
TestTableau_Pivot1
Description:

	In this test, we verify that the Pivot() function correctly performs a pivot operation on the tableau.

	We use an example problem from this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
	The initial tableau is:
		[	-15 	-25	  0 	  0 	  0 	  0 	  0 	]
		[	1 	 	1 	  1 	  0 	  0 	  0 	450		]
		[	0 	 	1 	  0 	  1 	  0 	  0 	300		]
		[	4 	 	5 	  0 	  0 	  1 	  0 	2000	]
		[	1 	 	0 	  0 	  0 	  0 	  1 	350		]
	And after pivoting using Bland's Rule, we expect the tableau to be:
		[	-15 	0 	  0 	  25 	  0 	  0 	  7500	]
		[	1 	 	0 	  1 	  -1 	  0 	  0 	150		]
		[	0 	 	1 	  0 	  1 	  0 	  0 	300		]
		[	4 	 	0 	  0 	  -5 	  1 	  0 	500		]
		[	1 	 	0 	  0 	  0 	  0 	  1 	350		]
*/
func TestTableau_Pivot1(t *testing.T) {
	// Setup

	// Create the test tableau
	testTableau, err := examples.GetTableauExample1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create the Bland's Rule selector
	selectionRule := selection.BlandsRule{}

	// Find the entering and exiting variables
	enteringVarIdx, exitingVarIdx, err := selectionRule.SelectEnteringAndExitingVariables(*testTableau)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	if enteringVarIdx != 1 {
		t.Errorf("Expected entering variable index to be 1, but got %d", enteringVarIdx)
	}
	if exitingVarIdx != 3 {
		t.Errorf("Expected exiting variable index to be 3, but got %d", exitingVarIdx)
	}

	// Perform the pivot operation
	newTab, err := testTableau.Pivot(enteringVarIdx, exitingVarIdx)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Compute the tableau as a dense matrix
	denseTableau := newTab.AsCompressedMatrix

	// Define the expected tableau
	expectedTableau := mat.NewDense(5, 7, []float64{
		-15, 0, 0, 25, 0, 0, 7500,
		1, 0, 1, -1, 0, 0, 150,
		0, 1, 0, 1, 0, 0, 300,
		4, 0, 0, -5, 1, 0, 500,
		1, 0, 0, 0, 0, 1, 350,
	})

	// Check that the tableau is correct
	if !mat.EqualApprox(denseTableau, expectedTableau, 1e-10) {
		t.Errorf("Expected tableau to be:\n%v\nbut got:\n%v", mat.Formatted(expectedTableau), mat.Formatted(denseTableau))
	}

}
