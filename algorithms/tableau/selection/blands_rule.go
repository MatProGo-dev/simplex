package selection

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/utils"
)

type BlandsRule struct{}

/*
Description:

	SelectEnteringVariable selects an entering variable from the given tableau using the Bland's Rule.
	If the result of this function is called:
		`idx := SelectEnteringVariable(tableau)`
	then `idx` is the index of the entering variable in `tableau.Variables`.
	If no entering variable can be found (i.e., the current solution is optimal),
	then -1 is returned as the variable index.
*/
func (br BlandsRule) SelectEnteringVariable(tableau utils.Tableau) int {
	// Setup
	minIndex := -1
	minValue := 0.0

	// Get the cost coefficients
	costCoefficients := tableau.C()

	// Iterate through the non-basic variables and find
	// the coefficient with the smallest index that is negative
	for _, nonBasicVarIdx := range tableau.NonBasicVariableIndicies() {
		if costCoefficients.AtVec(nonBasicVarIdx) < minValue {
			minIndex = nonBasicVarIdx
			minValue = costCoefficients.AtVec(nonBasicVarIdx)
		}
	}

	// If no entering variable is found, return -1
	return minIndex
}

/*
SelectExitingVariable

Description:

	SelectExitingVariable selects an exiting variable from the given tableau using the Blands Rule.

	If the result of this function is called:
		`idx := SelectExitingVariable(tableau, enteringVarIdx)`
	then `idx` is the index of the exiting variable in `tableau.Variables`.
	If no exiting variable can be found (i.e., the problem is unbounded),
	then -1 is returned as the variable index.
*/
func (br BlandsRule) SelectExitingVariable(tableau utils.Tableau, enteringVarIdx int) int {
	// Setup
	minIndex := -1
	minRatio := float64(symbolic.Infinity)

	// Get the relevant matrices
	A := tableau.A()
	b := tableau.B()

	// Create the vector of ratios
	ratios := make([]float64, tableau.NumberOfConstraints())
	for i := 0; i < tableau.NumberOfConstraints(); i++ {
		if A.At(i, enteringVarIdx) > 0 {
			ratios[i] = b.AtVec(i) / A.At(i, enteringVarIdx)
		} else {
			ratios[i] = -1.0 // Indicate that this variable cannot be used
		}
	}

	// Iterate through the ratios and find the smallest one
	for i, ratio := range ratios {
		if ratio >= 0 { // Only consider valid ratios
			if ratio < minRatio || (ratio == minRatio && tableau.BasicVariableIndicies[i] < tableau.BasicVariableIndicies[minIndex]) {
				minRatio = ratio
				minIndex = i
			}
		}
	}

	// If no exiting variable is found, return -1 (this is guaranteed by the structure of the algorithm)
	if minIndex == -1 {
		return -1
	}

	// Return the index of the exiting variable in tableau.Variables
	minIndex = tableau.BasicVariableIndicies[minIndex]
	return minIndex
}

func (br BlandsRule) SelectEnteringAndExitingVariables(tableau utils.Tableau) (int, int, error) {
	// Select the entering variable
	enteringVarIdx := br.SelectEnteringVariable(tableau)
	if enteringVarIdx == -1 {
		return -1, -1, nil // Optimal solution found, no entering variable
	}

	fmt.Println("Entering variable index:", enteringVarIdx)

	// Select the exiting variable
	exitingVarIdx := br.SelectExitingVariable(tableau, enteringVarIdx)
	if exitingVarIdx == -1 {
		return enteringVarIdx, -1, fmt.Errorf("BlandsRule: No exiting variable found, problem is unbounded (?)")
	}

	return enteringVarIdx, exitingVarIdx, nil
}
