package stanford_test

import (
	"strings"
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	stanford_algorithm1 "matprogo.dev/solvers/simplex/algorithms/stanford"
	"matprogo.dev/solvers/simplex/simplexSolver"
	"matprogo.dev/solvers/simplex/utils"
)

func TestStanfordAlgorithmState_NonBasicVariables1(t *testing.T) {
	// Setup
	exampleProblem1 := simplexSolver.GetTestProblem1()

	// Create the problem in standard form
	problemInStandardForm, slackVariables, err := exampleProblem1.ToLPStandardForm1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Collect the matrices that define the problem
	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	AAsDense := A.ToDense()
	bAsDense := b.ToVecDense()

	cAsPolynomialLike, tf := problemInStandardForm.Objective.Expression.(symbolic.PolynomialLikeScalar)
	if !tf {
		t.Errorf("Expected objective function to be a PolynomialLikeScalar, but got: %T", problemInStandardForm.Objective.Expression)
	}

	cAsVecDense := cAsPolynomialLike.LinearCoeff()

	// Create the initial state
	state0 := stanford_algorithm1.StanfordAlgorithmState{
		AllVariables:   problemInStandardForm.Variables,
		BasicVariables: slackVariables,
		NonBasicValues: nil, // We will set this later
		IterationCount: 0,
		A:              &AAsDense,
		B:              &bAsDense,
		C:              &cAsVecDense,
	}

	// Compute the non-basic variables
	varsOut := state0.GetNonBasicVariables()

	// Check that the non-basic variables do not include
	// the basic variables
	if len(varsOut) != len(utils.SetDifferenceOfVariables(varsOut, state0.BasicVariables)) {
		t.Errorf("Expected non-basic variables to not include basic variables. But there were some!")
	}

	// Check that the non basic variables are all the original variables (i.e., they do not contain the word "slack" in their names)
	for _, varOut := range varsOut {
		if strings.Contains(varOut.Name, "slack") {
			t.Errorf("Expected non-basic variable to not be a slack variable, but got: %v", varOut)
		}
	}

}

/*
TestStanfordAlgorithmState_BasicVariables1
Description:

	In this test, we verify that the GetBasicVariables() function
	correctly returns the basic variables in the current state of the algorithm.
	We test this on the initial set of basic variables which are all
	slack variables.
*/
func TestStanfordAlgorithmState_BasicVariables1(t *testing.T) {
	// Setup
	exampleProblem1 := simplexSolver.GetTestProblem1()

	// Create the problem in standard form
	problemInStandardForm, slackVariables, err := exampleProblem1.ToLPStandardForm1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Collect the matrices that define the problem
	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	AAsDense := A.ToDense()
	bAsDense := b.ToVecDense()

	cAsPolynomialLike, tf := problemInStandardForm.Objective.Expression.(symbolic.PolynomialLikeScalar)
	if !tf {
		t.Errorf("Expected objective function to be a PolynomialLikeScalar, but got: %T", problemInStandardForm.Objective.Expression)
	}

	cAsVecDense := cAsPolynomialLike.LinearCoeff()

	// Create the initial state
	state0 := stanford_algorithm1.StanfordAlgorithmState{
		AllVariables:   problemInStandardForm.Variables,
		BasicVariables: slackVariables,
		NonBasicValues: nil, // We will set this later
		IterationCount: 0,
		A:              &AAsDense,
		B:              &bAsDense,
		C:              &cAsVecDense,
	}

	// Compute the basic variables
	varsOut := state0.GetBasicVariables()

	// Check that the basic variables are all slack variables
	for _, varOut := range varsOut {
		if !strings.Contains(varOut.Name, "slack") {
			t.Errorf("Expected basic variable to be a slack variable, but got: %v", varOut)
		}
	}

	if len(varsOut) != len(state0.BasicVariables) {
		t.Errorf("Expected number of basic variables to match. Got %d expected %d",
			len(varsOut), len(state0.BasicVariables))
	}
}

/*
TestStanfordAlgorithmState_ReducedCostVector1
Description:

	In this test, we verify that the GetReducedCostVector() function
	returns the expected reduced cost vector for the initial state
	of the example given in:

	https://web.stanford.edu/class/msande310/lecture09.pdf

	We expect for the reduced cost vector to be:

	[ -1, -2, 0, 0, 0 ]
*/
func TestStanfordAlgorithmState_ReducedCostVector1(t *testing.T) {
	// Setup
	exampleProblem1 := simplexSolver.GetTestProblem2()

	// Create the problem in standard form
	problemInStandardForm, slackVariables, err := exampleProblem1.ToLPStandardForm1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Collect the matrices that define the problem
	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	AAsDense := A.ToDense()
	bAsDense := b.ToVecDense()

	cAsPolynomialLike, tf := problemInStandardForm.Objective.Expression.(symbolic.PolynomialLikeScalar)
	if !tf {
		t.Errorf("Expected objective function to be a PolynomialLikeScalar, but got: %T", problemInStandardForm.Objective.Expression)
	}

	cAsVecDense := cAsPolynomialLike.LinearCoeff(problemInStandardForm.Variables)

	// Create the initial state
	state0 := stanford_algorithm1.StanfordAlgorithmState{
		AllVariables:   problemInStandardForm.Variables,
		BasicVariables: slackVariables,
		NonBasicValues: nil, // We will set this later
		IterationCount: 0,
		A:              &AAsDense,
		B:              &bAsDense,
		C:              &cAsVecDense,
	}

	// t.Errorf("C: %v", state0.C)
	// cBasic, err := state0.CBasic()
	// if err != nil {
	// 	t.Errorf("Expected no error, but got: %v", err)
	// }
	// t.Errorf("CBasic: %v", cBasic)

	// Compute the reduced cost vector
	reducedCostVector, err := state0.GetReducedCostVector()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the reduced cost vector is as expected
	expectedReducedCostVector := []float64{-1, -2, 0, 0, 0}
	if !mat.EqualApprox(
		reducedCostVector,
		mat.NewVecDense(
			len(expectedReducedCostVector),
			expectedReducedCostVector,
		),
		1e-10,
	) {
		t.Errorf("Expected reduced cost vector to be: %v, but got: %v",
			expectedReducedCostVector, reducedCostVector)
	}
}

// /*
// TestStanfordAlgorithmState_ReducedCostVector2
// Description:

// 	In this test, we verify that the GetReducedCostVector() function
// 	returns the expected reduced cost vector for the contrived state
// 	of the example given in:

// 	https://web.stanford.edu/class/msande310/lecture09.pdf

// 	See page 17. We expect for the reduced cost vector to be:

// 	[ 0, 0, 0, 1, 1 ]
// */
// func TestStanfordAlgorithmState_ReducedCostVector2(t *testing.T) {
// 	// Setup
// 	exampleProblem1 := simplexSolver.GetTestProblem2()

// 	// Create the problem in standard form
// 	problemInStandardForm, _, err := exampleProblem1.ToLPStandardForm1()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}

// 	// Collect the matrices that define the problem
// 	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}

// 	AAsDense := A.ToDense()
// 	bAsDense := b.ToVecDense()

// 	cAsPolynomialLike, tf := problemInStandardForm.Objective.Expression.(symbolic.PolynomialLikeScalar)
// 	if !tf {
// 		t.Errorf("Expected objective function to be a PolynomialLikeScalar, but got: %T", problemInStandardForm.Objective.Expression)
// 	}

// 	cAsVecDense := cAsPolynomialLike.LinearCoeff(problemInStandardForm.Variables)

// 	// Create the initial state
// 	state0 := stanford_algorithm1.StanfordAlgorithmState{
// 		AllVariables:   problemInStandardForm.Variables,
// 		BasicVariables: problemInStandardForm.Variables[:3], // First three variables are basic
// 		NonBasicValues: nil,                                 // We will set this later
// 		IterationCount: 0,
// 		A:              &AAsDense,
// 		B:              &bAsDense,
// 		C:              &cAsVecDense,
// 	}

// 	t.Errorf("C: %v", state0.C)
// 	cBasic, err := state0.CBasic()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}
// 	t.Errorf("CBasic: %v", cBasic)

// 	ABasic, err := state0.ABasic()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}
// 	t.Errorf("ABasic: %v", ABasic)

// 	ANonBasic, err := state0.ANonBasic()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}
// 	t.Errorf("ANonBasic: %v", ANonBasic)

// 	// Compute the reduced cost vector
// 	reducedCostVector, err := state0.GetReducedCostVector()
// 	if err != nil {
// 		t.Errorf("Expected no error, but got: %v", err)
// 	}

// 	// Check that the reduced cost vector is as expected
// 	expectedReducedCostVector := []float64{0, 0, 0, 1, 1}
// 	if !mat.EqualApprox(
// 		reducedCostVector,
// 		mat.NewVecDense(
// 			len(expectedReducedCostVector),
// 			expectedReducedCostVector,
// 		),
// 		1e-10,
// 	) {
// 		t.Errorf("Expected reduced cost vector to be: %v, but got: %v",
// 			expectedReducedCostVector, reducedCostVector)
// 	}
// }
