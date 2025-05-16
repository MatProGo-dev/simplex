package simplexSolver

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

/*
ToStandardFormWithSlackVariables
Description:

	Transforms the given linear program (represented in an OptimizationProblem object)
	into a standard form (i.e., only linear equality constraints and a linear objective function).

		sense c^T * x
		subject to
		A * x = b
		x >= 0

	Where A is a matrix of coefficients, b is a vector of constants, and c is the vector of coefficients
	for the objective function. This method also returns the slack variables (i.e., the variables that
	are added to the problem to convert the inequalities into equalities).
*/
func ToStandardFormWithSlackVariables(problemIn *problem.OptimizationProblem) (*problem.OptimizationProblem, []symbolic.Variable) {

	// Setup

	// Create a new problem
	problemInStandardForm := problem.NewProblem(
		problemIn.Name + " (In Standard Form)",
	)

	// Copy over each of the

	// Add all variables to the new problem
	mapFromInToNewVariables := make(map[symbolic.Variable]symbolic.Expression)
	for _, varII := range problemIn.Variables {
		problemInStandardForm.AddVariable()
		nVariables := len(problemInStandardForm.Variables)
		mapFromInToNewVariables[varII] = problemInStandardForm.Variables[nVariables-1]
	}

	// Add all constraints to the new problem
	slackVariables := []symbolic.Variable{}
	for _, constraint := range problemIn.Constraints {
		// Create a new expression by substituting the variables according
		// to the map we created above
		oldLHS := constraint.Left()
		newLHS := oldLHS.SubstituteAccordingTo(mapFromInToNewVariables)

		oldRHS := constraint.Right()
		newRHS := oldRHS.SubstituteAccordingTo(mapFromInToNewVariables)

		switch constraint.ConstrSense() {
		case symbolic.SenseEqual:
			// No need to do anything
		case symbolic.SenseGreaterThanEqual:
			switch concreteConstraint := constraint.(type) {
			case symbolic.ScalarConstraint:
				// Add a new SCALAR slack variable to the right hand side
				problemInStandardForm.AddVariableClassic(0.0, symbolic.Infinity.Constant(), symbolic.Continuous)
				nVariables := len(problemInStandardForm.Variables)
				problemInStandardForm.Variables[nVariables-1].Name = problemInStandardForm.Variables[nVariables-1].Name + " (slack)"
				slackVariables = append(
					slackVariables,
					problemInStandardForm.Variables[nVariables-1],
				)

				newRHS = newRHS.Plus(problemInStandardForm.Variables[nVariables-1])
			case symbolic.VectorConstraint:
				// Add a new VECTOR slack variable to the right hand side
				// TODO(Kwesi): Revisit this when we have a proper Len() method for constraints.
				dims := concreteConstraint.Dims()
				nRows := dims[0]
				problemInStandardForm.AddVariableVectorClassic(
					nRows,
					0.0,
					symbolic.Infinity.Constant(),
					symbolic.Continuous,
				)
				nVariables := len(problemInStandardForm.Variables)
				for jj := nRows - 1; jj >= nRows; jj-- {
					problemInStandardForm.Variables[nVariables-1-jj].Name = problemInStandardForm.Variables[nVariables-1-jj].Name + " (slack)"
					slackVariables = append(
						slackVariables,
						problemInStandardForm.Variables[nVariables-1-jj],
					)
				}

				// Add the slack variable to the right hand side
				newRHS = newRHS.Plus(
					symbolic.VariableVector(problemInStandardForm.Variables[nVariables-1-nRows : nVariables-1]),
				)
			default:
				panic(
					fmt.Sprintf(
						"Unexpected constraint type: %T for \"ToStandardFormWithSlackVariables\" with %v sense",
						constraint,
						constraint.ConstrSense(),
					),
				)

			}
		case symbolic.SenseLessThanEqual:
			// Use a switch statement to handle different dimensions of the constraint
			switch concreteConstraint := constraint.(type) {
			case symbolic.ScalarConstraint:
				// Add a new SCALAR slack variable to the left hand side
				problemInStandardForm.AddVariableClassic(0.0, symbolic.Infinity.Constant(), symbolic.Continuous)
				nVariables := len(problemInStandardForm.Variables)
				problemInStandardForm.Variables[nVariables-1].Name = problemInStandardForm.Variables[nVariables-1].Name + " (slack)"
				slackVariables = append(
					slackVariables,
					problemInStandardForm.Variables[nVariables-1],
				)
				newLHS = newLHS.Plus(problemInStandardForm.Variables[nVariables-1])
			case symbolic.VectorConstraint:
				// Add a new VECTOR slack variable to the left hand side
				// TODO(Kwesi): Revisit this when we have a proper Len() method for constraints.
				dims := concreteConstraint.Dims()
				nRows := dims[0]
				problemInStandardForm.AddVariableVectorClassic(
					nRows,
					0.0,
					symbolic.Infinity.Constant(),
					symbolic.Continuous,
				)
				nVariables := len(problemInStandardForm.Variables)
				for jj := nRows - 1; jj >= 0; jj-- {
					problemInStandardForm.Variables[nVariables-1-jj].Name = problemInStandardForm.Variables[nVariables-1-jj].Name + " (slack)"
					slackVariables = append(
						slackVariables,
						problemInStandardForm.Variables[nVariables-1-jj],
					)
					// fmt.Printf("Slack variable %d: %v\n", jj, problemInStandardForm.Variables[nVariables-1-jj])
				}
				// Add the slack variable to the left hand side
				newLHS = newLHS.Plus(
					symbolic.VariableVector(problemInStandardForm.Variables[nVariables-1-nRows : nVariables-1]),
				)
			default:
				panic(
					fmt.Sprintf(
						"Unexpected constraint type %T for \"ToStandardFormWithSlackVariables\" with %v sense",
						constraint,
						constraint.ConstrSense(),
					),
				)
			}
		default:
			panic("Unknown constraint sense: " + constraint.ConstrSense().String())
		}

		newConstraint := newLHS.Comparison(
			newRHS,
			symbolic.SenseEqual,
		)

		// Add the new constraint to the problem
		problemInStandardForm.Constraints = append(
			problemInStandardForm.Constraints,
			newConstraint,
		)
	}

	// Now, let's create the new objective function by substituting the variables
	// according to the map we created above
	newObjectiveExpression := problemIn.Objective.Expression.SubstituteAccordingTo(
		mapFromInToNewVariables,
	)
	problemInStandardForm.SetObjective(
		newObjectiveExpression,
		problemIn.Objective.Sense,
	)

	fmt.Printf("The slack variables are: %v\n", slackVariables)

	// Return the new problem and the slack variables
	return problemInStandardForm, slackVariables
}

/*
SetDifferenceOfVariables
Description:

	Performs a set difference between two slices of variables.
	Returns a slice of variables that are in the first slice but not in the second.
*/
func SetDifferenceOfVariables(a, b []symbolic.Variable) []symbolic.Variable {
	// Create a map to track the variables in b
	varsInB := make(map[symbolic.Variable]struct{})
	for _, varB := range b {
		varsInB[varB] = struct{}{}
	}

	// Create a slice to hold the result
	result := []symbolic.Variable{}

	// Iterate through the first slice and add variables that are not in b
	for _, varA := range a {
		if _, found := varsInB[varA]; !found {
			result = append(result, varA)
		}
	}

	return result
}

/*
SliceMatrixAccordingToVariableSet
Description:

	Collects the columns of the input matrix A that correspond
	to the input variable set in an optimization problem.
	Returns the new matrix.
*/
func SliceMatrixAccordingToVariableSet(
	problemIn *problem.OptimizationProblem,
	matrixIn symbolic.KMatrix,
	variablesIn []symbolic.Variable,
) (symbolic.KMatrix, error) {
	// Setup
	nVariables := len(problemIn.Variables)
	dims := matrixIn.Dims()
	out := symbolic.ZerosMatrix(dims[0], len(variablesIn))
	matrixInAsDense := matrixIn.ToDense()

	// Check that the number of variables in the problem matches the number of columns in the matrix
	if nVariables != dims[1] {
		return nil, fmt.Errorf(
			"Number of variables in the problem (%d) does not match number of columns in the matrix (%d)",
			nVariables,
			dims[1],
		)
	}

	// Iterate through each variable in the problem
	for ii := 0; ii < len(variablesIn); ii++ {
		// Determine if this variable is in the set of variables
		variable := variablesIn[ii]
		idxII, err := symbolic.FindInSlice(variable, problemIn.Variables)
		if err != nil {
			// There was an issue searching through this list!
			panic(
				fmt.Sprintf(
					"Error searching through variable list %v for %v: %v",
					variablesIn,
					variable,
					err.Error(),
				),
			)
		}

		// Set the column of the output matrix that corresponds to the extracted column
		for jj := 0; jj < dims[0]; jj++ {
			fmt.Printf("Setting out(%d, %d) = matrixIn(%d, %d)\n", jj, ii, jj, idxII)
			out.Set(jj, ii, matrixInAsDense.At(jj, idxII))
		}
	}

	return symbolic.DenseToKMatrix(out), nil
}
