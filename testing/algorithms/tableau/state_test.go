package tableau

import (
	"testing"

	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
	tableau_algorithm1 "github.com/MatProGo-dev/simplex/algorithms/tableau"
	tableau_termination "github.com/MatProGo-dev/simplex/algorithms/tableau/termination"
	"github.com/MatProGo-dev/simplex/utils"
	"github.com/MatProGo-dev/simplex/utils/examples"
)

/*
TestTableau_CalculateOptimalValue1
Description:

	In this test, we verify that the CalculateOptimalValue() function correctly computes the optimal value
	of the objective function for a known tableau.

	We use an example problem from this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
	The optimal solution of the objective function is:
		x1 = 125
		x2 = 300
		s1 = 25
		s2 = 0
		s3 = 0
		s4 = 225
*/
func TestTableau_CalculateOptimalSolution1(t *testing.T) {
	// Setup

	// Create the test problem
	testProblem := examples.GetTestProblem5()

	// Use the optimization solver to solve the problem

	// Create initial Tableau state from the problem
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}
	state0 := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Update state two times
	state1, err := state0.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}
	state2, err := state1.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}

	// Calculate the optimal solution
	solVec, err := state2.CalculateOptimalSolution()
	if err != nil {
		t.Errorf("there was an issue calculating the optimal value: %v", err)
	}

	// Check that the solution is correct
	expectedSol := []float64{125.0, 300.0, 25.0, 0.0, 0.0, 225.0}
	for ii, val := range expectedSol {
		if solVec.AtVec(ii) != val {
			t.Errorf("Expected solution value %v at index %d, but got %v", val, ii, solVec.AtVec(ii))
		}
	}
}

/*
TestTableauAlgorithmState_Check1
Description:

	Tests that the Check() method returns no error for a valid state.
*/
func TestTableauAlgorithmState_Check1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	err = state.Check()

	// Verify
	if err != nil {
		t.Errorf("Expected no error for valid state, but got: %v", err)
	}
}

/*
TestTableauAlgorithmState_Check2
Description:

	Tests that the Check() method returns an error when IterationCount is negative.
*/
func TestTableauAlgorithmState_Check2(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: -5,
	}

	// Test
	err = state.Check()

	// Verify
	if err == nil {
		t.Errorf("Expected an error for negative iteration count, but got none")
	}
}

/*
TestTableauAlgorithmState_CheckTerminationCondition1
Description:

	Tests that CheckTerminationCondition() returns false for initial state
	(not yet optimal).
*/
func TestTableauAlgorithmState_CheckTerminationCondition1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	terminated, err := state.CheckTerminationCondition()

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if terminated {
		t.Errorf("Expected state to not be terminated initially, but it was")
	}
}

/*
TestTableauAlgorithmState_CheckTerminationCondition2
Description:

	Tests that CheckTerminationCondition() returns true for an optimal state.
*/
func TestTableauAlgorithmState_CheckTerminationCondition2(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state0 := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Move to optimal state
	state1, err := state0.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}
	state2, err := state1.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}

	// Test
	terminated, err := state2.CheckTerminationCondition()

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if !terminated {
		t.Errorf("Expected state to be terminated at optimal, but it wasn't")
	}
}

/*
TestTableauAlgorithmState_GetBasicVariables1
Description:

	Tests that GetBasicVariables() returns the correct basic variables.
*/
func TestTableauAlgorithmState_GetBasicVariables1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	basicVars := state.GetBasicVariables()

	// Verify
	if len(basicVars) == 0 {
		t.Errorf("Expected non-empty basic variables, but got empty slice")
	}

	// For the initial tableau, we should have slack variables as basic variables
	// (4 constraints = 4 slack variables)
	if len(basicVars) != 4 {
		t.Errorf("Expected 4 basic variables, but got %d", len(basicVars))
	}
}

/*
TestTableauAlgorithmState_GetNonBasicVariables1
Description:

	Tests that GetNonBasicVariables() returns the correct non-basic variables.
*/
func TestTableauAlgorithmState_GetNonBasicVariables1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	nonBasicVars := state.GetNonBasicVariables()

	// Verify
	if len(nonBasicVars) == 0 {
		t.Errorf("Expected non-empty non-basic variables, but got empty slice")
	}

	// For the initial tableau, we should have original variables as non-basic
	// (2 original variables)
	if len(nonBasicVars) != 2 {
		t.Errorf("Expected 2 non-basic variables, but got %d", len(nonBasicVars))
	}
}

/*
TestTableauAlgorithmState_NumberOfIterations1
Description:

	Tests that NumberOfIterations() returns the correct iteration count.
*/
func TestTableauAlgorithmState_NumberOfIterations1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 7,
	}

	// Test
	iterations := state.NumberOfIterations()

	// Verify
	if iterations != 7 {
		t.Errorf("Expected 7 iterations, but got %d", iterations)
	}
}

/*
TestTableauAlgorithmState_NumberOfVariables1
Description:

	Tests that NumberOfVariables() returns the correct total variable count.
*/
func TestTableauAlgorithmState_NumberOfVariables1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	numVars := state.NumberOfVariables()

	// Verify
	// 2 original variables + 4 slack variables = 6 total
	if numVars != 6 {
		t.Errorf("Expected 6 variables, but got %d", numVars)
	}
}

/*
TestTableauAlgorithmState_NumberOfConstraints1
Description:

	Tests that NumberOfConstraints() returns the correct constraint count.
*/
func TestTableauAlgorithmState_NumberOfConstraints1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	numConstraints := state.NumberOfConstraints()

	// Verify
	// Test problem 5 has 4 constraints
	if numConstraints != 4 {
		t.Errorf("Expected 4 constraints, but got %d", numConstraints)
	}
}

/*
TestTableauAlgorithmState_GetReducedCostVector1
Description:

	Tests that GetReducedCostVector() returns a vector without errors.
*/
func TestTableauAlgorithmState_GetReducedCostVector1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	reducedCost, err := state.GetReducedCostVector()

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if reducedCost == nil {
		t.Errorf("Expected non-nil reduced cost vector")
	}

	// The reduced cost vector should have length equal to number of variables
	if reducedCost.Len() != state.NumberOfVariables() {
		t.Errorf("Expected reduced cost vector length %d, but got %d", 
			state.NumberOfVariables(), reducedCost.Len())
	}
}

/*
TestTableauAlgorithmState_GetShadowPrice1
Description:

	Tests that GetShadowPrice() returns a vector without errors.
*/
func TestTableauAlgorithmState_GetShadowPrice1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, _, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Test
	shadowPrice, err := state.GetShadowPrice()

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if shadowPrice == nil {
		t.Errorf("Expected non-nil shadow price vector")
	}

	// The shadow price vector should have length equal to number of variables
	if shadowPrice.Len() != state.NumberOfVariables() {
		t.Errorf("Expected shadow price vector length %d, but got %d", 
			state.NumberOfVariables(), shadowPrice.Len())
	}
}

/*
TestTableauAlgorithmState_CreateOptimalValuesMap1
Description:

	Tests that CreateOptimalValuesMap() creates a proper values map.
*/
func TestTableauAlgorithmState_CreateOptimalValuesMap1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, varMap, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state0 := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Move to optimal state
	state1, err := state0.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}
	state2, err := state1.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}

	// Test
	valuesMap, err := state2.CreateOptimalValuesMap(varMap)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if valuesMap == nil {
		t.Errorf("Expected non-nil values map")
	}

	// Should have values for the original variables
	if len(valuesMap) == 0 {
		t.Errorf("Expected non-empty values map")
	}

	// Check expected values for TestProblem5
	// x1 = 125, x2 = 300
	for varID, value := range valuesMap {
		if varID == 0 && value != 125.0 {
			t.Errorf("Expected x1 = 125.0, but got %v", value)
		}
		if varID == 1 && value != 300.0 {
			t.Errorf("Expected x2 = 300.0, but got %v", value)
		}
	}
}

/*
TestTableauAlgorithmState_ToSolution1
Description:

	Tests that ToSolution() creates a valid SimplexSolution.
*/
func TestTableauAlgorithmState_ToSolution1(t *testing.T) {
	// Setup
	testProblem := examples.GetTestProblem5()
	initialTableau, varMap, err := utils.GetInitialTableauFrom(testProblem)
	if err != nil {
		t.Errorf("there was an issue creating the initial tableau: %v", err)
	}

	state0 := tableau_algorithm1.TableauAlgorithmState{
		Tableau:        &initialTableau,
		IterationCount: 0,
	}

	// Move to optimal state
	state1, err := state0.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}
	state2, err := state1.CalculateNextState()
	if err != nil {
		t.Errorf("there was an issue calculating the next state: %v", err)
	}

	// Test
	solution, err := state2.ToSolution(
		tableau_termination.OptimalSolutionFound,
		varMap,
		testProblem,
	)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check solution status
	if solution.GetStatus() != solution_status.OPTIMAL {
		t.Errorf("Expected OPTIMAL status, but got %v", solution.GetStatus())
	}

	// Check iteration count
	if solution.Iterations != 2 {
		t.Errorf("Expected 2 iterations, but got %d", solution.Iterations)
	}

	// Check that we have variable values
	if len(solution.VariableValues) == 0 {
		t.Errorf("Expected non-empty variable values")
	}

	// Check that problem is attached
	if solution.GetProblem() == nil {
		t.Errorf("Expected non-nil problem")
	}
}

