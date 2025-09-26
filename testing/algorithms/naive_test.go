package algorithms_test

import (
	"testing"

	"github.com/MatProGo-dev/simplex/algorithms"
	"github.com/MatProGo-dev/simplex/simplexSolver"
	"gonum.org/v1/gonum/mat"
)

/*
Test_NaiveAlgorithm_ComputeFeasibleSolution1
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
func Test_NaiveAlgorithm_ComputeFeasibleSolution1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem3()

	// Create the solver
	solver, err := simplexSolver.For(problemIn, simplexSolver.Configuration{IterationLimit: 100})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Extract the simplex algorithm and internal state from the solver
	algo, err := solver.CreateAlgorithm(algorithms.TypeNaive)
	if err != nil {
		t.Errorf("Expected no error when creating algorithm; received %v", err)
	}

	naiveAlgo, ok := algo.(*algorithms.NaiveAlgorithm)
	if !ok {
		t.Errorf("Failed to convert the algorithm interface into a Naive Algorithm object (which should be possible given our inputs!)")
	}

	initialState := solver.InternalState

	// Compute the feasible solution
	solution, err := naiveAlgo.ComputeFeasibleBasicSolution(initialState)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the solution is correct
	expectedSolution := mat.NewVecDense(3, []float64{4, 6, 8})
	if !mat.EqualApprox(solution, expectedSolution, 1e-10) {
		t.Errorf("Expected feasible solution to be %v, but got %v", expectedSolution, solution)
	}
}

/*
Test_NaiveAlgorithm_SolveLoop1
Description:

	In this test, we verify that the operations
	of the solver loop are correct.
	We will use a problem from this youtube video:
		https://youtu.be/XMLysZSPsug?si=KMoouByHAV3TTK7h&t=377
	Our method should find that the first Basic Feasible Solution
	is {3 2 4 2} and the objective value should be zero.
*/
func Test_NaiveAlgorithm_SolveLoop1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem4()

	// Create the solver
	solver, err := simplexSolver.For(problemIn, simplexSolver.Configuration{IterationLimit: 100})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Extract the simplex algorithm and internal state from the solver
	algo, err := solver.CreateAlgorithm(algorithms.TypeNaive)
	if err != nil {
		t.Errorf("Expected no error when creating algorithm; received %v", err)
	}

	naiveAlgo, ok := algo.(*algorithms.NaiveAlgorithm)
	if !ok {
		t.Errorf("Failed to convert the algorithm interface into a Naive Algorithm object (which should be possible given our inputs!)")
	}

	initialState := solver.InternalState

	// Find the first basic feasible solution
	basic1, err := naiveAlgo.ComputeFeasibleBasicSolution(initialState)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the solution is correct
	expectedFirstBasicSolution := mat.NewVecDense(4, []float64{3, 2, 4, 2})
	if !mat.EqualApprox(basic1, expectedFirstBasicSolution, 1e-10) {
		t.Errorf("Expected feasible solution to be %v, but got %v", expectedFirstBasicSolution, basic1)
	}

	// Compute the value of the objective function that corresponds to the basic solution
	obj1, err := naiveAlgo.ComputeObjectiveFunctionValueWithFeasibleBasicSolution(initialState, basic1)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check that the objective value is correct (should be zero)
	expectedObj1 := 0.0
	if obj1 != expectedObj1 {
		t.Errorf("Expected objective value to be %v, but got %v", expectedObj1, obj1)
	}

}

/*
Test_NaiveAlgorithm_Solve1
Description:

	In this test, we verify that the operations
	of the NaiveAlgorithm.Solve() are correct.
	We will use a problem from this youtube video:
		https://youtu.be/XMLysZSPsug?si=KMoouByHAV3TTK7h&t=377
	Our method should find that the first Basic Feasible Solution
	is {3 2 4 2} and the objective value should be zero.
*/
func Test_NaiveAlgorithm_Solve1(t *testing.T) {
	// Setup
	problemIn := simplexSolver.GetTestProblem4()

	// Create the solver
	solver, err := simplexSolver.For(problemIn, simplexSolver.Configuration{IterationLimit: 10})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Extract the simplex algorithm and internal state from the solver
	algo, err := solver.CreateAlgorithm(algorithms.TypeNaive)
	if err != nil {
		t.Errorf("Expected no error when creating algorithm; received %v", err)
	}

	naiveAlgo, ok := algo.(*algorithms.NaiveAlgorithm)
	if !ok {
		t.Errorf("Failed to convert the algorithm interface into a Naive Algorithm object (which should be possible given our inputs!)")
	}

	initialState := solver.InternalState

	// Find the first basic feasible solution
	_, err = naiveAlgo.Solve(initialState)
	// sol1, err := naiveAlgo.Solve(initialState)
	if err != nil {
		t.Errorf("Unexpected error during solve loops: %v", err)
	}

	// t.Errorf(
	// 	"Objective value after %v loops: %v",
	// 	naiveAlgo.IterationLimit,
	// 	sol1.Objective,
	// )

}
