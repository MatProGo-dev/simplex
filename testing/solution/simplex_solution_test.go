package solution_test

import (
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
)

/*
TestSimplexSolution_GetValueMap1
Description:

	Tests that the GetValueMap() method returns the correct variable values map.
*/
func TestSimplexSolution_GetValueMap1(t *testing.T) {
	// Setup
	variableValues := map[uint64]float64{
		0: 1.0,
		1: 2.0,
		2: 3.0,
	}

	sol := simplex_solution.SimplexSolution{
		VariableValues: variableValues,
		Objective:      10.0,
		Status:         solution_status.OPTIMAL,
		Iterations:     5,
	}

	// Test
	result := sol.GetValueMap()

	// Verify
	if len(result) != len(variableValues) {
		t.Errorf("Expected %d values, but got %d", len(variableValues), len(result))
	}

	for key, expectedValue := range variableValues {
		if result[key] != expectedValue {
			t.Errorf("Expected value %v for variable %d, but got %v", expectedValue, key, result[key])
		}
	}
}

/*
TestSimplexSolution_GetOptimalValue1
Description:

	Tests that the GetOptimalValue() method returns the correct objective value.
*/
func TestSimplexSolution_GetOptimalValue1(t *testing.T) {
	// Setup
	expectedObjective := 42.5

	sol := simplex_solution.SimplexSolution{
		VariableValues: map[uint64]float64{0: 1.0},
		Objective:      expectedObjective,
		Status:         solution_status.OPTIMAL,
		Iterations:     3,
	}

	// Test
	result := sol.GetOptimalValue()

	// Verify
	if result != expectedObjective {
		t.Errorf("Expected objective value %v, but got %v", expectedObjective, result)
	}
}

/*
TestSimplexSolution_GetStatus1
Description:

	Tests that the GetStatus() method returns the correct solution status
	for an optimal solution.
*/
func TestSimplexSolution_GetStatus1(t *testing.T) {
	// Setup
	sol := simplex_solution.SimplexSolution{
		VariableValues: map[uint64]float64{0: 1.0},
		Objective:      10.0,
		Status:         solution_status.OPTIMAL,
		Iterations:     2,
	}

	// Test
	result := sol.GetStatus()

	// Verify
	if result != solution_status.OPTIMAL {
		t.Errorf("Expected status %v, but got %v", solution_status.OPTIMAL, result)
	}
}

/*
TestSimplexSolution_GetStatus2
Description:

	Tests that the GetStatus() method returns the correct solution status
	for an infeasible solution.
*/
func TestSimplexSolution_GetStatus2(t *testing.T) {
	// Setup
	sol := simplex_solution.SimplexSolution{
		VariableValues: map[uint64]float64{},
		Objective:      0.0,
		Status:         solution_status.INFEASIBLE,
		Iterations:     1,
	}

	// Test
	result := sol.GetStatus()

	// Verify
	if result != solution_status.INFEASIBLE {
		t.Errorf("Expected status %v, but got %v", solution_status.INFEASIBLE, result)
	}
}

/*
TestSimplexSolution_GetStatus3
Description:

	Tests that the GetStatus() method returns the correct solution status
	for an unbounded solution.
*/
func TestSimplexSolution_GetStatus3(t *testing.T) {
	// Setup
	sol := simplex_solution.SimplexSolution{
		VariableValues: map[uint64]float64{},
		Objective:      0.0,
		Status:         solution_status.UNBOUNDED,
		Iterations:     0,
	}

	// Test
	result := sol.GetStatus()

	// Verify
	if result != solution_status.UNBOUNDED {
		t.Errorf("Expected status %v, but got %v", solution_status.UNBOUNDED, result)
	}
}

/*
TestSimplexSolution_GetProblem1
Description:

	Tests that the GetProblem() method returns the correct original problem.
*/
func TestSimplexSolution_GetProblem1(t *testing.T) {
	// Setup
	originalProblem := problem.NewProblem("TestProblem")
	originalProblem.AddVariableVector(2)

	sol := simplex_solution.SimplexSolution{
		VariableValues:  map[uint64]float64{0: 1.0, 1: 2.0},
		Objective:       10.0,
		Status:          solution_status.OPTIMAL,
		Iterations:      4,
		OriginalProblem: originalProblem,
	}

	// Test
	result := sol.GetProblem()

	// Verify
	if result == nil {
		t.Errorf("Expected a non-nil problem, but got nil")
	}

	if result != originalProblem {
		t.Errorf("Expected the same problem reference, but got a different one")
	}

	if result.Name != "TestProblem" {
		t.Errorf("Expected problem name 'TestProblem', but got '%s'", result.Name)
	}
}

/*
TestSimplexSolution_GetProblem2
Description:

	Tests that the GetProblem() method returns nil when no problem is attached.
*/
func TestSimplexSolution_GetProblem2(t *testing.T) {
	// Setup
	sol := simplex_solution.SimplexSolution{
		VariableValues:  map[uint64]float64{0: 1.0},
		Objective:       5.0,
		Status:          solution_status.OPTIMAL,
		Iterations:      2,
		OriginalProblem: nil,
	}

	// Test
	result := sol.GetProblem()

	// Verify
	if result != nil {
		t.Errorf("Expected nil problem, but got %v", result)
	}
}

/*
TestSimplexSolution_EmptyVariableValues
Description:

	Tests that the SimplexSolution can handle an empty variable values map.
*/
func TestSimplexSolution_EmptyVariableValues(t *testing.T) {
	// Setup
	sol := simplex_solution.SimplexSolution{
		VariableValues: map[uint64]float64{},
		Objective:      0.0,
		Status:         solution_status.INFEASIBLE,
		Iterations:     0,
	}

	// Test
	result := sol.GetValueMap()

	// Verify
	if result == nil {
		t.Errorf("Expected non-nil map, but got nil")
	}

	if len(result) != 0 {
		t.Errorf("Expected empty map, but got %d values", len(result))
	}
}

/*
TestSimplexSolution_CompleteStructure
Description:

	Tests a complete SimplexSolution with all fields populated.
*/
func TestSimplexSolution_CompleteStructure(t *testing.T) {
	// Setup
	variableValues := map[uint64]float64{
		0: 125.0,
		1: 300.0,
		2: 25.0,
	}
	objective := 9875.0
	iterations := 10

	originalProblem := problem.NewProblem("CompleteProblem")
	originalProblem.AddVariableVector(3)

	sol := simplex_solution.SimplexSolution{
		VariableValues:  variableValues,
		Objective:       objective,
		Status:          solution_status.OPTIMAL,
		Iterations:      iterations,
		OriginalProblem: originalProblem,
	}

	// Test all methods
	if sol.GetOptimalValue() != objective {
		t.Errorf("Expected objective %v, but got %v", objective, sol.GetOptimalValue())
	}

	if sol.GetStatus() != solution_status.OPTIMAL {
		t.Errorf("Expected status OPTIMAL, but got %v", sol.GetStatus())
	}

	valueMap := sol.GetValueMap()
	if len(valueMap) != len(variableValues) {
		t.Errorf("Expected %d variables, but got %d", len(variableValues), len(valueMap))
	}

	prob := sol.GetProblem()
	if prob == nil || prob.Name != "CompleteProblem" {
		t.Errorf("Expected problem 'CompleteProblem', but got %v", prob)
	}

	if sol.Iterations != iterations {
		t.Errorf("Expected %d iterations, but got %d", iterations, sol.Iterations)
	}
}
