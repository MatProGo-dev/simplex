package dictionary_test

import (
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms/dictionary"
)

/*
TestDictionaryAlgorithmState_CheckBasicVariableIndicies1
Description:

	Tests that the CheckBasicVariableIndicies() method properly returns an
	error when one of the integers in the slice CheckBasicVariableIndicies
	is negative.
*/
func TestDictionaryAlgorithmState_CheckBasicVariableIndicies1(t *testing.T) {
	// Setup
	N := 3
	vv1 := symbolic.NewVariableVector(N)

	state0 := dictionary.DictionaryAlgorithmState{
		AllVariables:          vv1,
		BasicVariableIndicies: []int{0, -2, 1},
	}

	// Run algorithm
	err := state0.CheckBasicVariableIndicies()
	if err == nil {
		t.Errorf("Expected an error in test; received none!")
	}

}

/*
TestDictionaryAlgorithmState_CheckBasicVariableIndicies2
Description:

	Tests that the CheckBasicVariableIndicies() method properly returns an
	error when one of the integers in the slice CheckBasicVariableIndicies
	is too high.
*/
func TestDictionaryAlgorithmState_CheckBasicVariableIndicies2(t *testing.T) {
	// Setup
	N := 3
	vv1 := symbolic.NewVariableVector(N)

	state0 := dictionary.DictionaryAlgorithmState{
		AllVariables:          vv1,
		BasicVariableIndicies: []int{0, 1, 6},
	}

	// Run algorithm
	err := state0.CheckBasicVariableIndicies()
	if err == nil {
		t.Errorf("Expected an error in test; received none!")
	}

}
