package getConstraintMatrices_test

import (
	"testing"

	"matprogo.dev/solvers/simplex/simplexSolver"
	"matprogo.dev/solvers/simplex/simplexSolver/getConstraintMatrices"
)

/*
Test_NInequalityConstraints1
Description:

	In this test, we test that the NInequalityConstraints() function correctly counts the number of constraints
	in the example optimization problem #1. The answer should be 1.
*/
func Test_NInequalityConstraints1(t *testing.T) {
	// Get Example Problem 1
	problem1 := simplexSolver.GetTestProblem1()

	// There should only be 1 inequality constraint in this problem
	if getConstraintMatrices.NInequalityConstraints(problem1) != 1 {
		t.Errorf(
			"Expected Problem 1 to have 1 Inequality constraint; function claims that there are %v!",
			getConstraintMatrices.NInequalityConstraints(problem1),
		)
	}
}

/*
Test_NInequalityConstraints2
Description:

	In this test, we test that the NInequalityConstraints() function correctly counts the number of constraints
	in the example optimization problem #1. The answer should be 5.
*/
func Test_NInequalityConstraints2(t *testing.T) {
	// Get Example Problem 1
	problem1 := simplexSolver.GetTestProblem2()

	// There should only be 1 inequality constraint in this problem
	if getConstraintMatrices.NInequalityConstraints(problem1) != 5 {
		t.Errorf(
			"Expected Problem 1 to have 5 Inequality constraint; function claims that there are %v!",
			getConstraintMatrices.NInequalityConstraints(problem1),
		)
	}
}

/*
Test_
*/
