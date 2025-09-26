package utils

import (
	"fmt"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

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
	matrixIn symbolic.KMatrix,
	originalSetOfVariables []symbolic.Variable,
	subsetOfVariables []symbolic.Variable,
) (symbolic.KMatrix, error) {
	// Setup
	nVariables := len(originalSetOfVariables)
	dims := matrixIn.Dims()
	out := symbolic.ZerosMatrix(dims[0], len(subsetOfVariables))
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
	for ii := 0; ii < len(subsetOfVariables); ii++ {
		// Determine if this variable is in the set of variables
		variable := subsetOfVariables[ii]
		idxII, err := symbolic.FindInSlice(variable, originalSetOfVariables)
		if err != nil {
			// There was an issue searching through this list!
			panic(
				fmt.Sprintf(
					"Error searching through variable list %v for %v: %v",
					subsetOfVariables,
					variable,
					err.Error(),
				),
			)
		}

		// Set the column of the output matrix that corresponds to the extracted column
		for jj := 0; jj < dims[0]; jj++ {
			// fmt.Printf("Setting out(%d, %d) = matrixIn(%d, %d)\n", jj, ii, jj, idxII)
			out.Set(jj, ii, matrixInAsDense.At(jj, idxII))
		}
	}

	return symbolic.DenseToKMatrix(out), nil
}

/*
SliceVectorAccordingToVariableSet
Description:

	Collects the elements of the input vector vectorIn and
	assumes that each one corresponds to a variable in the input
	`originalSetOfVariables`. We then return a new vector which
	is just the elements of the first vector that correspond to the
	variables in the input `subsetOfVariables`.
	Returns the new vector.
*/
func SliceVectorAccordingToVariableSet(
	vectorIn symbolic.KVector,
	originalSetOfVariables []symbolic.Variable,
	subsetOfVariables []symbolic.Variable,
) (symbolic.KVector, error) {
	// Setup
	out := symbolic.ZerosVector(len(subsetOfVariables))
	vectorInAsDense := vectorIn.ToVecDense()

	// Check that the number of variables in the problem matches the number of columns in the matrix
	if len(originalSetOfVariables) != vectorIn.Len() {
		return nil, fmt.Errorf(
			"Number of variables in the problem (%d) does not match number of columns in the matrix (%d)",
			len(originalSetOfVariables),
			vectorIn.Len(),
		)
	}

	// Iterate through each variable in the problem
	for ii := 0; ii < len(subsetOfVariables); ii++ {
		// Determine if this variable is in the set of variables
		variable := subsetOfVariables[ii]
		idxII, err := symbolic.FindInSlice(variable, originalSetOfVariables)
		if err != nil {
			// There was an issue searching through this list!
			panic(
				fmt.Sprintf(
					"Error searching through variable list %v for %v: %v",
					subsetOfVariables,
					variable,
					err.Error(),
				),
			)
		}

		out.SetVec(ii, vectorInAsDense.At(idxII, 0))
	}

	return symbolic.VecDenseToKVector(out), nil
}

/*
DefinePartialAssignmentVector
Description:

	This method takes a vector `vIn` of length N and uses
	it to create a new vector `vExtended` of length M (M >= N),
	The new vector is a vector representing the partial
	assignment of values in the set `originalSetOfVariables`.
	We assign the values `vIn` to the entries corresponding
	to `subsetOfVariables`; the rest are left zero.
*/
func DefinePartialAssignmentVector(
	vIn *mat.VecDense,
	subsetOfVariables []symbolic.Variable,
	originalSetOfVariables []symbolic.Variable,
) (*mat.VecDense, error) {
	// Input Checking
	// - Input Vector is the same length as the subset
	if vIn.Len() != len(subsetOfVariables) {
		return nil, fmt.Errorf(
			"The length of the input (%v) must match the length of the subset of variables(%v)!",
			vIn.Len(),
			len(subsetOfVariables),
		)
	}

	//- subsetOfVariables is a strict subset of originalSetOfVariables
	setDifference := SetDifferenceOfVariables(subsetOfVariables, originalSetOfVariables)
	if len(setDifference) != 0 {
		return nil, fmt.Errorf(
			"Subset of variables (%v) is not a strict subset of original set of variables (%v)",
			subsetOfVariables,
			originalSetOfVariables,
		)
	}

	// Setup
	M := len(originalSetOfVariables)
	out := symbolic.ZerosVector(M)

	// Assign values to out
	for ii := 0; ii < len(subsetOfVariables); ii++ {
		// Find the variable
		sII := subsetOfVariables[ii]

		// Find the index of sII
		indexOfII, err := symbolic.FindInSlice(sII, originalSetOfVariables)
		if err != nil {
			return nil, fmt.Errorf(
				"couldn't find the index of variable %v in the list of variables %v",
				sII,
				originalSetOfVariables,
			)
		}

		// Assign the value at index ii in vIn to out at index indexOfII
		out.SetVec(indexOfII, vIn.AtVec(ii))
	}

	return &out, nil
}
