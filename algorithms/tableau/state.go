package tableau_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms"
	"github.com/MatProGo-dev/simplex/utils"
	"gonum.org/v1/gonum/mat"
)

type TableauAlgorithmState struct {
	Tableau        *utils.Tableau
	IterationCount int
}

func (state *TableauAlgorithmState) A() *mat.Dense {
	return state.Tableau.A()
}

func (state *TableauAlgorithmState) ABasic() *mat.Dense {
	ABasic, err := state.Tableau.ABasic()
	if err != nil {
		panic(err)
	}

	return ABasic
}

func (state *TableauAlgorithmState) B() *mat.VecDense {
	return state.Tableau.B()
}

func (state *TableauAlgorithmState) C() *mat.VecDense {
	return state.Tableau.C()
}

func (state *TableauAlgorithmState) CBasic() *mat.VecDense {
	cBasic, err := state.Tableau.CBasic()
	if err != nil {
		panic(err)
	}
	return cBasic
}

/*
Check
Description:

	This method checks whether or not the TableauAlgorithmState is well-defined.
	Specifically, we check:
	- Iteration Count >= 0
	- Tableau has:
		+ len(AllVariables) + 2 columns
*/
func (state *TableauAlgorithmState) Check() error {
	// Check that the count is a non-negative number
	if state.IterationCount < 0 {
		return algorithms.MakeIterationCountIsNegativeError(state)
	}

	// Check that the number of columns is equal to len(AllVariables) + 1
	err := state.Tableau.Check()
	if err != nil {
		return err
	}

	// All Checks passed
	return nil
}

func (state *TableauAlgorithmState) CheckTerminationCondition() (bool, error) {
	// Collect the Reduced Cost Vector
	r, err := state.GetReducedCostVector()
	if err != nil {
		return false, fmt.Errorf(
			"TableaudAlgorithm: Failed to get reduced cost vector: %v",
			err,
		)
	}

	// If all of the elements are non-negative, then return true.
	for ii := 0; ii < r.Len(); ii++ {
		if r.AtVec(ii) < 0 {
			return false, nil
		}
	}
	return true, nil
}

// func GetStateFromInitialProblem(prob problem.OptimizationProblem) TableauAlgorithmState {
// 	// Setup

// }

/*
GetBasicVariables
Description:

	Returns the basic variables in the current state of the algorithm.
*/
func (state *TableauAlgorithmState) GetBasicVariables() []symbolic.Variable {
	return state.Tableau.BasicVariables()
}

/*
GetNonBasicVariables
Description:

	Returns the non-basic variables in the current state of the algorithm.
	This should be all variables that are not part of the basic variables.
*/
func (state *TableauAlgorithmState) GetNonBasicVariables() []symbolic.Variable {
	return state.Tableau.NonBasicVariables()
}

/*
NumberOfIterations
Description:

	Returns the iteration count.
*/
func (state *TableauAlgorithmState) NumberOfIterations() int {
	return state.IterationCount
}

/*
NumberOfVariables
Description:

	Returns the total count of all variables in this variable.
*/
func (state *TableauAlgorithmState) NumberOfVariables() int {
	return state.Tableau.NumberOfBasicVariables() + state.Tableau.NumberOfNonBasicVariables()
}

/*
NumberOfConstraints
Description:

	Returns the number of constraints as inferred by the size of the tableau.
*/
func (state *TableauAlgorithmState) NumberOfConstraints() int {
	nRows, _ := state.Tableau.AsCompressedMatrix.Dims()
	return nRows - 1
}

func (state *TableauAlgorithmState) XBasic() (*mat.VecDense, error) {
	// Checking inputs
	err := state.Check()
	if err != nil {
		return nil, err
	}

	// Algorithm
	var xB *mat.VecDense
	ABasic := state.ABasic()
	xB.MulVec(ABasic, state.B())

	return xB, nil
}

func (state *TableauAlgorithmState) GetShadowPrice() (*mat.VecDense, error) {
	// Input Checking
	err := state.Check()
	if err != nil {
		return nil, err
	}

	// Setup
	ABasic := state.ABasic()

	// Assemble the complex expression: c_B^T * A_B^(-1) * A
	var ABasicInv mat.Dense
	err = ABasicInv.Inverse(ABasic)
	if err != nil {
		return nil, fmt.Errorf("[tableau algorithm state]: Inversion failed, cannot compute reduced costs (%v)", err)
	}

	cBasic := state.CBasic()

	var prod mat.Dense
	prod.Mul(&ABasicInv, state.A())
	var shadowPriceT mat.VecDense
	shadowPriceT.MulVec(prod.T(), cBasic)

	return &shadowPriceT, nil

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
func (state *TableauAlgorithmState) GetReducedCostVector() (*mat.VecDense, error) {
	// Check the state for validity
	err := state.Check()
	if err != nil {
		return nil, err // Invalid state, cannot compute reduced costs
	}

	// Invert A_B if possible
	ABasic := state.ABasic()

	var ABasicInv mat.Dense
	err = ABasicInv.Inverse(ABasic)
	if err != nil {
		return nil, fmt.Errorf("[tableau algorithm state]: Inversion failed, cannot compute reduced costs (%v)", err)
	}

	// Compute c_B^T * A_B^(-1) * A
	cBasic := state.CBasic()

	var temp mat.Dense
	temp.Mul(&ABasicInv, state.A())
	var reducedCostT mat.VecDense
	reducedCostT.MulVec(temp.T(), cBasic)

	// Compute the final reduced cost vector
	var finalReducedCost mat.VecDense
	finalReducedCost.SubVec(state.C(), &reducedCostT)

	return &finalReducedCost, nil
}

func (state *TableauAlgorithmState) ToSolution(currentStatus problem.OptimizationStatus) problem.Solution {
	// Construct Solution Map

	// Assemble Solution Output
	return problem.Solution{
		Values:    map[uint64]float64{},
		Objective: -1.0,
		Status:    currentStatus,
	}
}
