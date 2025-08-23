package stanford_algorithm1

import (
	"fmt"

	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"matprogo.dev/solvers/simplex/utils"
)

// TODO: Change BasicVariables name to be lower case and make
// a new method BasicVariables() that returns the basic variables
// in the state. This will help with consistency and avoid confusion.
type StanfordAlgorithmState struct {
	AllVariables   []symbolic.Variable
	BasicVariables []symbolic.Variable
	NonBasicValues *mat.VecDense
	IterationCount int

	// Other fields for the problem definition
	A *mat.Dense
	B *mat.VecDense
	C *mat.VecDense
}

/*
Check
Description:

	Checks that the current state of the algorithm is valid.
	Explicitly checks the following conditions:
	1. The matrices ABasic, ANonBasic, B, CBasic, and CNonBasic are not nil.
	2. The number of basic variables matches the number of columns in A_B.
	3. The number of non-basic variables matches the number of columns in A_N.
*/
func (state *StanfordAlgorithmState) Check() error {
	// Check if all required matrices are present
	if state.A == nil {
		return fmt.Errorf("StanfordAlgorithmState: A matrix is nil.")
	}
	if state.B == nil {
		return fmt.Errorf("StanfordAlgorithmState: B vector is nil.")
	}
	if state.C == nil {
		return fmt.Errorf("StanfordAlgorithmState: C vector is nil.")
	}

	// Check dimensions of A
	nRowsInA, nColsInA := state.A.Dims()
	if nRowsInA == 0 || nColsInA == 0 {
		return fmt.Errorf("StanfordAlgorithmState: A matrix has invalid dimensions (%v, %v).", nRowsInA, nColsInA)
	}

	if nColsInA != len(state.AllVariables) {
		return fmt.Errorf(
			"StanfordAlgorithmState: Number of columns in A (%v) does not match length of AllVariables (%v).",
			nColsInA, len(state.AllVariables),
		)
	}

	// Check dimensions of B, CBasic, and CNonBasic
	// if nRowsInABasic != state.B.Len() {
	// 	return fmt.Errorf(
	// 		"StanfordAlgorithmState: Number of rows in ABasic (%v) does not match length of B vector (%v).",
	// 		nRowsInABasic, state.B.Len(),
	// 	)
	// }
	// if nColsInABasic != state.CBasic.Len() {
	// 	return fmt.Errorf(
	// 		"StanfordAlgorithmState: Number of columns in ABasic (%v) does not match length of CBasic vector (%v).",
	// 		nColsInABasic, state.CBasic.Len(),
	// 	)
	// }
	// if nColsInANonBasic != state.CNonBasic.Len() {
	// 	return fmt.Errorf(
	// 		"StanfordAlgorithmState: Number of columns in ANonBasic (%v) does not match length of CNonBasic vector (%v).",
	// 		nColsInANonBasic, state.CNonBasic.Len(),
	// 	)
	// }

	// Otherwise, return nil
	return nil
}

/*
GetBasicVariables
Description:

	Returns the basic variables in the current state of the algorithm.
*/
func (state *StanfordAlgorithmState) GetBasicVariables() []symbolic.Variable {
	return state.BasicVariables
}

/*
GetNonBasicVariables
Description:

	Returns the non-basic variables in the current state of the algorithm.
	This should be all variables that are not part of the basic variables.
*/
func (state *StanfordAlgorithmState) GetNonBasicVariables() []symbolic.Variable {
	return utils.SetDifferenceOfVariables(state.AllVariables, state.BasicVariables)
}

/*
NumberOfBasicVariables
Description:

	Returns the number of basic variables in the current state of the algorithm.
*/
func (state *StanfordAlgorithmState) NumberOfBasicVariables() int {
	return len(state.BasicVariables)
}

/*
NumberOfNonBasicVariables
Description:

	Returns the number of non-basic variables in the current state of the algorithm.
*/
func (state *StanfordAlgorithmState) NumberOfNonBasicVariables() int {
	return len(state.GetNonBasicVariables())
}

/*
ABasic
Description:

	Returns the matrix of coefficients of the basic variables (A_B) in the current state of
	the algorithm.
*/
func (state *StanfordAlgorithmState) ABasic() (*mat.Dense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return ABasic
	}

	// Slice the A Matrix according to the basic variables
	A := state.A
	ABasic, err := utils.SliceMatrixAccordingToVariableSet(
		getKMatrix.From(A),
		state.AllVariables,
		state.BasicVariables,
	)
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to slice A matrix for basic variables (%v)", err)
	}

	ABasicAsDense := ABasic.ToDense()

	return &ABasicAsDense, nil
}

/*
ANonBasic
Description:

	Returns the matrix of coefficients of the non-basic variables (A_N) in the current state of
	the algorithm.
*/
func (state *StanfordAlgorithmState) ANonBasic() (*mat.Dense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return ANonBasic
	}

	// Slice the A Matrix according to the non-basic variables
	A := state.A
	ANonBasic, err := utils.SliceMatrixAccordingToVariableSet(
		getKMatrix.From(A),
		state.AllVariables,
		state.GetNonBasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to slice A matrix for non-basic variables (%v)", err)
	}

	ANonBasicAsDense := ANonBasic.ToDense()

	return &ANonBasicAsDense, nil
}

/*
CBasic
Description:

	Returns the cost vector for the basic variables (c_B) in the current state of
	the algorithm.
*/
func (state *StanfordAlgorithmState) CBasic() (*mat.VecDense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return CBasic
	}

	// Slice the cost vector according to the basic variables
	cBasic, err := utils.SliceVectorAccordingToVariableSet(
		getKVector.From(state.C),
		state.AllVariables,
		state.BasicVariables,
	)
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to slice cost vector for basic variables (%v)", err)
	}

	cBasicAsVecDense := cBasic.ToVecDense()

	return &cBasicAsVecDense, nil
}

/*
CNonBasic
Description:

	Returns the cost vector for the non-basic variables (c_N) in the current state of
	the algorithm.
*/
func (state *StanfordAlgorithmState) CNonBasic() (*mat.VecDense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return CNonBasic
	}

	// Slice the cost vector according to the non-basic variables
	cNonBasic, err := utils.SliceVectorAccordingToVariableSet(
		getKVector.From(state.C),
		state.AllVariables,
		state.GetNonBasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to slice cost vector for non-basic variables (%v)", err)
	}

	cNonBasicAsVecDense := cNonBasic.ToVecDense()

	return &cNonBasicAsVecDense, nil
}

/*
ReducedCostVector
Description:

	Computes the reduced cost vector for the current state of the algorithm.
	This is computed as:
		c_N - c_B^T * A_N^(-1) * A_B
	where:
		c_N is the cost vector for the non-basic variables,
		c_B is the cost vector for the basic variables,
		A_N is the matrix of coefficients for the non-basic variables,
		A_B is the matrix of coefficients for the basic variables.

Returns:
  - A pointer to the reduced cost vector (mat.VecDense) if successful.
  - An error if the state is invalid or if the computation fails.
*/
func (state *StanfordAlgorithmState) GetReducedCostVector() (*mat.VecDense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot compute reduced costs
	}

	// Invert A_B if possible
	ABasic, err := state.ABasic()
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to get ANonBasic matrix (%v)", err)

	}

	var ABasicInv mat.Dense
	err = ABasicInv.Inverse(ABasic)
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Inversion failed, cannot compute reduced costs (%v)", err)
	}

	// Compute c_B^T * A_B^(-1) * A
	cBasic, err := state.CBasic()
	if err != nil {
		return nil, fmt.Errorf("StanfordAlgorithmState: Failed to get CBasic vector (%v)", err)
	}

	var temp mat.Dense
	temp.Mul(&ABasicInv, state.A)
	var reducedCostT mat.VecDense
	reducedCostT.MulVec(temp.T(), cBasic)

	// Compute the final reduced cost vector
	var finalReducedCost mat.VecDense
	finalReducedCost.SubVec(state.C, &reducedCostT)

	return &finalReducedCost, nil
}

/*
ComputeMinimumRatioTest
Description:

	Computes the minimum ratio test for the current state of the algorithm.
	Returns the outgoing variable xO with its associated increase, theta.
*/
func (state *StanfordAlgorithmState) ComputeMinimumRatioTest(enteringVarIndex int) (symbolic.Variable, float64, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return symbolic.Variable{}, 0.0, err // Invalid state, cannot compute minimum ratio test
	}

	// Compute (ABasic)^(-1) * A's e-th column
	var ABasicInv mat.Dense
	ABasic, err := state.ABasic()
	if err != nil {
		return symbolic.Variable{}, 0.0, fmt.Errorf("StanfordAlgorithmState: Failed to get ABasic matrix (%v)", err)
	}
	ABasicInv.Inverse(ABasic)

	var ABasicAe mat.VecDense
	ABasicAe.MulVec(&ABasicInv, state.A.ColView(enteringVarIndex))

	// Compute the (ABasic)^(-1) * b's vector
	var ABasicB mat.VecDense
	ABasicB.MulVec(&ABasicInv, state.B)

	// Find the entering variable (most negative reduced cost)
	theta := 0.0
	outgoingVarIndex := -1
	for ii := 0; ii < ABasicB.Len(); ii++ {
		// Check to see if ABasicAe is positive
		if ABasicAe.AtVec(ii) <= 0.0 {
			continue // Skip this variable, it cannot be the outgoing variable
		}
		// Compute the ratio
		ratio := ABasicB.AtVec(ii) / ABasicAe.AtVec(ii)
		if ratio < theta {
			theta = ratio
			outgoingVarIndex = ii
		}
	}

	if outgoingVarIndex == -1 {
		return symbolic.Variable{}, 0.0, fmt.Errorf("StanfordAlgorithmState: No entering variable found (all reduced costs are non-negative)")
	}

	outgoingVariable := state.GetBasicVariables()[outgoingVarIndex]

	return outgoingVariable, theta, nil
}
