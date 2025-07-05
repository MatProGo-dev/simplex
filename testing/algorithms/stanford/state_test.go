package stanford_test

import (
	"strings"
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
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
	varsOut := state0.NonBasicVariables()

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
