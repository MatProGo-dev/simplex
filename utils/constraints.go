package utils

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

/*
ExtractScalarConstraints
Description:

	This method analyzes all constraints in an OptimizationProblem and converts them all
	into scalar constraints.

TODO:

	Add this to SymbolicMath.go.
*/
func ExtractScalarConstraints(constraints []symbolic.Constraint) []symbolic.ScalarConstraint {
	// Setup
	var out []symbolic.ScalarConstraint

	// Iterate through all constraints
	for _, constraint := range constraints {
		// Switch statement based on the type of the constraint
		switch concreteConstraint := constraint.(type) {
		case symbolic.ScalarConstraint:
			out = append(out, concreteConstraint)
		case symbolic.VectorConstraint:
			for ii := 0; ii < concreteConstraint.Len(); ii++ {
				out = append(out, concreteConstraint.AtVec(ii))
			}
		case symbolic.MatrixConstraint:
			dims := concreteConstraint.Dims()
			for rowIdx := 0; rowIdx < dims[0]; rowIdx++ {
				for colIdx := 0; colIdx < dims[1]; colIdx++ {
					out = append(out, concreteConstraint.At(rowIdx, colIdx))
				}
			}
		default:
			panic(
				fmt.Errorf(
					"The received constraint type (%T) is not supported by ExtractScalarConstraints!",
					constraint,
				),
			)
		}
	}

	return out
}
