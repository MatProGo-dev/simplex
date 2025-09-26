package utils_test

import (
	"math"
	"testing"

	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/utils"
	"gonum.org/v1/gonum/mat"
)

/*
TestDefinePartialAssignmentVector1
Description:

	This test will verify that the DefinePartialAssignmentVector
*/
func TestDefinePartialAssignmentVector1(t *testing.T) {
	// Setup
	M := 5

	// Create the original set of variables
	vv1 := symbolic.NewVariableVector(M)

	// Create a subset of the original set of variables
	subset := []symbolic.Variable{vv1[1], vv1[3]}
	N := len(subset)

	// Create the partial assignment vector
	partialAssignment := mat.NewVecDense(N, []float64{10, 20})

	// Call the function
	result, err := utils.DefinePartialAssignmentVector(
		partialAssignment,
		subset,
		vv1,
	)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check the result
	if result.Len() != M {
		t.Errorf(
			"Expected result length to be %d, but got %d",
			M,
			result.Len(),
		)
	}

	// Make sure that the first, third and fifth entries are zero
	for _, ii := range []int{0, 2, 4} {
		if math.Abs(result.AtVec(ii)) > 1e-7 {
			t.Errorf("Expected result.AtVec(%d) to be 0, but got %f", ii, result.AtVec(ii))
		}
	}

	// Make sure that the second and fourth entries are 10 and 20 respectively
	if math.Abs(result.AtVec(1)-10) > 1e-7 {
		t.Errorf("Expected result.AtVec(1) to be 10, but got %f", result.AtVec(1))
	}
	if math.Abs(result.AtVec(3)-20) > 1e-7 {
		t.Errorf("Expected result.AtVec(3) to be 20, but got %f", result.AtVec(3))
	}

}

/*
TestDefinePartialAssignmentVector2
Description:

	This test will verify that the DefinePartialAssignmentVector function
	returns an error when the input vector length does not match the subset length.
*/
func TestDefinePartialAssignmentVector2(t *testing.T) {
	// Setup
	M := 5

	// Create the original set of variables
	vv1 := symbolic.NewVariableVector(M)

	// Create a subset of the original set of variables
	subset := []symbolic.Variable{vv1[1], vv1[3]}
	N := len(subset)

	// Create the partial assignment vector with incorrect length
	partialAssignment := mat.NewVecDense(N+1, []float64{10, 20, 30})

	// Call the function
	_, err := utils.DefinePartialAssignmentVector(
		partialAssignment,
		subset,
		vv1,
	)
	// TODO(Kwesi): Check the exact error message
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}
}

/*
TestDefinePartialAssignmentVector3
Description:

	This test will verify that the DefinePartialAssignmentVector function
	returns an error when the subset is not a strict subset of the original set.
*/
func TestDefinePartialAssignmentVector3(t *testing.T) {
	// Setup
	M := 5

	// Create the original set of variables
	vv1 := symbolic.NewVariableVector(M)

	// Create a subset of the original set of variables that is not a strict subset
	subset := []symbolic.Variable{vv1[1], symbolic.NewVariable()}

	// Create the partial assignment vector
	partialAssignment := mat.NewVecDense(len(subset), []float64{10, 20})

	// Call the function
	_, err := utils.DefinePartialAssignmentVector(
		partialAssignment,
		subset,
		vv1,
	)
	// TODO(Kwesi): Check the exact error message
	if err == nil {
		t.Errorf("Expected an error, but got none")
	}

}
