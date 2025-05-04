package getConstraintMatrices

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
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
func LinearConstraintMatrices(problemIn *problem.OptimizationProblem) (symbolic.KMatrix, symbolic.KVector) {
	// Setup

	// Create a variable for tracking all scalar Left Hand Sides
	scalarLHSs := []symbolic.ScalarExpression{}
	scalarRHSs := []symbolic.ScalarExpression{}
	vectorLHSs := []symbolic.VectorExpression{}
	vectorRHSs := []symbolic.VectorExpression{}

	// Iterate through all constraints
	for _, constraint := range problemIn.Constraints {
		// Skip this constraint if it is an equality constraint
		if constraint.ConstrSense() == symbolic.SenseEqual {
			continue
		}

		// Otherwise, let's add this constraint to the list of constraints
		// depending on if it is a scalar or vector
		switch constraint := constraint.(type) {
		case symbolic.ScalarConstraint:
			newLHS := constraint.Left()
			newLHS = newLHS.Minus(constraint.Right()).(symbolic.ScalarExpression)

			newRHS := symbolic.K(
				constraint.Right().(symbolic.ScalarExpression).Constant(),
			)

			if constraint.ConstrSense() == symbolic.SenseGreaterThanEqual {
				newLHS = newLHS.Multiply(-1.0)
				newRHS = newRHS.Multiply(-1.0)
			}

			// Add the newLHS to the list of scalar LHSs
			scalarLHSs = append(scalarLHSs, newLHS)
			scalarRHSs = append(scalarRHSs, newRHS)

		case symbolic.VectorConstraint:
			newLHS := constraint.Left()
			newLHS = newLHS.Minus(constraint.Right()).(symbolic.VectorExpression)

			newRHS := symbolic.KVector(
				constraint.Right().(symbolic.VectorExpression).Constant(),
			)
			if constraint.ConstrSense() == symbolic.SenseGreaterThanEqual {
				newLHS = newLHS.Multiply(-1.0)
				newRHS = newRHS.Multiply(-1.0)
			}
			// Add the newLHS to the list of vector LHSs
			vectorLHSs = append(vectorLHSs, newLHS)
			vectorRHSs = append(vectorRHSs, newRHS)
		default:
			panic(
				fmt.Sprintf("Unknown constraint type: %T", constraint),
			)
		}
	}

	// Now, let's create the output matrix and vector using the collection of scalar and vector constraints
	scalarConstraintsExist := len(scalarLHSs) > 0
	vectorConstraintsExist := len(vectorLHSs) > 0

	var AOut symbolic.KMatrix
	var bOut symbolic.KVector

	if scalarConstraintsExist {
		// Construct A Matrix
		lhs0AsPolynomial := scalarLHSs[0].(symbolic.PolynomialLikeScalar)
		lhs0A := symbolic.KVector(
			lhs0AsPolynomial.LinearCoeff(),
		).Transpose()
		AOut = lhs0A

		for ii := 1; ii < len(scalarLHSs); ii++ {
			lhsA := scalarLHSs[ii].(symbolic.PolynomialLikeScalar)
			lhsB := getKVector.From(lhsA.LinearCoeff()).Transpose()
			AOut = symbolic.VStack(
				AOut,
				lhsB,
			).(symbolic.KMatrix)
		}

		// Construct b Matrix
		rhs0AsK := scalarRHSs[0].(symbolic.KScalar)
		bOut = symbolic.KVector(*mat.NewVecDense(1, []float64{rhs0AsK.Constant()}))
		for ii := 1; ii < len(scalarRHSs); ii++ {
			rhsAsK := scalarRHSs[ii].(symbolic.KScalar)
			bOut = symbolic.VStack(
				bOut,
				getKVector.From(*mat.NewVecDense(1, []float64{rhsAsK.Constant()})),
			).(symbolic.KVector)
		}
	}

	if vectorConstraintsExist {
		// Construct A Matrix, if it doesn't exist already
		if !scalarConstraintsExist {
			lhs0AsPolynomial := vectorLHSs[0].(symbolic.PolynomialLikeVector)
			lhs0A := lhs0AsPolynomial.LinearCoeff()
			AOut = lhs0A
		}
		// Add to the matrix
		startIndex = 0
		if !scalarConstraintsExist {
			startIndex = 1
		}
		for ii := startIndex; ii < len(vectorLHSs); ii++ {
			lhsA := vectorLHSs[ii].(symbolic.PolynomialLikeVector)
			lhsB := lhsA.LinearCoeff() // TODO(kwesi): Fix this function
			AOut = symbolic.VStack(
				AOut,
				lhsB,
			).(symbolic.KMatrix)
		}

		// Construct b Matrix, if it doesn't exist already
		if !scalarConstraintsExist {
			rhs0AsK := vectorRHSs[0].(symbolic.KVector)
			bOut = rhs0AsK
		}
		// Add to the vector
		for ii := startIndex; ii < len(vectorRHSs); ii++ {
			rhsAsK := vectorRHSs[ii].(symbolic.KVector)
			bOut = symbolic.VStack(
				bOut,
				rhsAsK,
			).(symbolic.KVector)
		}
	}

	// Algorithm
	return AOut, vecOut
}
