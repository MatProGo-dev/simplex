package getConstraintMatrices

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

/*
NInequalityConstraints
Description:

	Computes the number of inequality constraints in a given problem.
*/
func NInequalityConstraints(problemIn *problem.OptimizationProblem) int {
	// Setup
	count := 0

	// Iterate through all constraints
	for _, constraint := range problemIn.Constraints {
		// Count the constraint if the constraint sense is NOT Equal
		if constraint.ConstrSense() != symbolic.SenseEqual {
			count++
		}
	}
	return count
}

/*
LinearConstraintMatrices
Description:

	This function collects all linear constraints from the associated problem
	(all constraints of the form `expr1` <= `expr2` OR `expr1` >= `expr2`)
	and attempts to create the large matrix equality
		A * x <= b
	from those inequalities.

Notes:

	This function assumes that the input problem is a Linear Program (may want to change this/add
	input checking for this soon.)
*/
// func LinearConstraintMatrices(problemIn *problem.OptimizationProblem) (symbolic.KMatrix, symbolic.KVector) {
// 	// Setup

// 	// Create a variable for tracking all scalar Left Hand Sides

// 	// Iterate through all constraints
// 	for _, constraint := range problemIn.Constraints {
// 		// Skip this constraint if it is an equality constraint
// 		switch constraint.ConstrSense() {
// 		case symbolic.SenseEqual:
// 			continue
// 		case symbolic.SenseLessThanEqual:

// 		}

// 	}

// 	// Algorithm
// 	return matOut, vecOut
// }
