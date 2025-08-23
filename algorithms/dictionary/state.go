package dictionary

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

type DictionaryAlgorithmState struct {
	AllVariables          []symbolic.Variable
	BasicVariableIndicies []int
	IterationCount        int

	// Structure of Constraints
	ObjectiveExpression   symbolic.Expression
	DictionaryConstraints []symbolic.ScalarConstraint
}

/*
Check
Description:

	Returns an error if:
	- The indices in BasicVariableIndicies ARE OUTSIDE of the range acceptable for AllVariables
	- There exists a constraint that:
		+ is not an equality constraint, or
		+ contains variables that are not in AllVariables
*/
func (state *DictionaryAlgorithmState) Check() error {
	// Check values of BasicVariableIndicies
	err := state.CheckBasicVariableIndicies()
	if err != nil {
		return err
	}

	// Check the constraints
	for idx, constraint := range state.DictionaryConstraints {
		if constraint.ConstrSense() != symbolic.SenseEqual {
			return fmt.Errorf(
				"Constraint #%v is not an Equality constraint; found sense %v",
				idx,
				constraint.ConstrSense(),
			)
		}
	}

	// All checks passed!
	return nil
}

/*
CheckBasicVariableIndicies
Description:

	Verifies that the BasicVariableIndicies variable contains indicies that make sense
	given the value of AllVariables.
*/
func (state *DictionaryAlgorithmState) CheckBasicVariableIndicies() error {
	// Setup
	N := len(state.AllVariables)
	allowableIndiciesSuffix := fmt.Sprintf("only values between 0 and %v supported", N)

	// Check all indicies in the BasicVariableIndicies slice.
	for _, bvIndex := range state.BasicVariableIndicies {
		if bvIndex < 0 {
			return fmt.Errorf(
				"BasicVariableIndicies contains a negative value (%v); %v",
				bvIndex,
				allowableIndiciesSuffix,
			)
		}

		if bvIndex >= N {
			return fmt.Errorf(
				"BasicVariableIndicies contains a value outside of the allowable range (%v); %v",
				bvIndex,
				allowableIndiciesSuffix,
			)
		}
	}

	// All Checks Passed!
	return nil
}
