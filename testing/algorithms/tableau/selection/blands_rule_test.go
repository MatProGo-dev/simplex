package tableau_test

import (
	"testing"

	"github.com/MatProGo-dev/simplex/algorithms/tableau/selection"
	"github.com/MatProGo-dev/simplex/utils/examples"
)

/*
TestBlandsRule_SelectEnteringVariable1
Description:

	This test will verify that the Bland's Rule selection algorithm correctly selects
	the entering variable for the problem shown in minute 21:42 of this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
*/
func TestBlandsRule_SelectEnteringVariable1(t *testing.T) {
	// Setup

	// Create the test tableau
	testTableau, err := examples.GetTableauExample1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create the Bland's Rule selector
	selectionRule := selection.BlandsRule{}

	// Find the entering variable
	enteringVarIdx := selectionRule.SelectEnteringVariable(*testTableau)

	if enteringVarIdx != 1 {
		t.Errorf("Expected entering variable index to be 1, but got %d", enteringVarIdx)
	}
}

/*
TestBlandsRule_SelectExitingVariable1
Description:

	This test will verify that the Bland's Rule selection algorithm correctly selects
	the exiting variable for the problem shown in minute 21:42 of this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
	The exiting variable selected should be the variable with index 3 (the slack variable
	for the second constraint).
*/
func TestBlandsRule_SelectExitingVariable1(t *testing.T) {
	// Setup

	// Create the test tableau
	testTableau, err := examples.GetTableauExample1()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Create the Bland's Rule selector
	selectionRule := selection.BlandsRule{}

	// Find the entering variable
	enteringVarIdx := selectionRule.SelectEnteringVariable(*testTableau)
	if enteringVarIdx != 1 {
		t.Errorf("Expected entering variable index to be 1, but got %d", enteringVarIdx)
	}

	// Find the exiting variable
	exitingVarIdx := selectionRule.SelectExitingVariable(*testTableau, enteringVarIdx)
	if exitingVarIdx != 3 {
		t.Errorf("Expected exiting variable index to be 3, but got %d", exitingVarIdx)
	}
}

/*
TestBlandsRule_SelectEnteringAndExitingVariables1
Description:

	This test will verify that the Bland's Rule selection algorithm correctly selects
	the entering and exiting variables for the problem shown in minute 21:42 of this youtube video:
		https://www.youtube.com/watch?v=-7mCHWpQ9Fw&t=883s
	The entering variable selected should be the variable with index 1 (x2),
	and the exiting variable selected should be the variable with index 3 (the slack variable
	for the second constraint).
*/
func TestBlandsRule_SelectEnteringAndExitingVariables1(t *testing.T) {
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
}
