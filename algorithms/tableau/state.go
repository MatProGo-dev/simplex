package tableau_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/MatProInterface.go/solution"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms"
	"github.com/MatProGo-dev/simplex/algorithms/tableau/selection"
	tableau_termination "github.com/MatProGo-dev/simplex/algorithms/tableau/termination"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
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
	return state.Tableau.NumberOfConstraints()
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

func (state *TableauAlgorithmState) CalculateNextState() (TableauAlgorithmState, error) {
	// Input Checking
	err := state.Check()
	if err != nil {
		return TableauAlgorithmState{}, err
	}

	// fmt.Println("Calculating next state from tableau:", mat.Formatted(state.Tableau.AsCompressedMatrix))

	// Select the pivot column and row (i.e., the entering and exiting variables in the tableau)
	// Here, we use Bland's Rule to select the entering variable
	// TODO(Kwesi): Make other rules available
	selectionRule := selection.BlandsRule{}
	enteringVarIdx, exitingVarIdx, err := selectionRule.SelectEnteringAndExitingVariables(*state.Tableau)
	if err != nil {
		return TableauAlgorithmState{}, VariableSelectionError{EnteringVarIndex: enteringVarIdx, ExitingVarIndex: exitingVarIdx}
	}

	// fmt.Println("Entering variable index: ", enteringVarIdx, " (", state.Tableau.Variables[enteringVarIdx], ")")
	// fmt.Println("Exiting variable index: ", exitingVarIdx, " (", state.Tableau.Variables[exitingVarIdx], ")")

	// Create the new tableau
	newTab, err := state.Tableau.Pivot(enteringVarIdx, exitingVarIdx)
	if err != nil {
		return TableauAlgorithmState{}, fmt.Errorf("TableauAlgorithmState: Failed to pivot tableau (%v)", err)
	}

	// Create the new state
	return TableauAlgorithmState{
		Tableau:        &newTab,
		IterationCount: state.IterationCount + 1,
	}, nil
}

func (state *TableauAlgorithmState) CalculateOptimalSolution() (mat.VecDense, error) {
	// Input Checking
	err := state.Check()
	if err != nil {
		return mat.VecDense{}, err
	}

	// Setup
	numVars := state.NumberOfVariables()
	solutionVec := mat.NewVecDense(numVars, nil)

	// Create a linear system of variables consisting of:
	// - The constraints
	// - The non-basic variables set to zero
	A, b := state.A(), state.B()
	numConstraints, _ := A.Dims()
	numNonBasic := state.Tableau.NumberOfNonBasicVariables()

	// Augment the A and b matrices with the non-basic variable constraints
	AAugmented := mat.NewDense(numConstraints+numNonBasic, numVars, nil)
	AAugmented.Copy(A)
	bAugmented := mat.NewVecDense(numConstraints+numNonBasic, nil)
	bAugmented.CopyVec(b)

	// Add the non-basic variable constraints
	nonBasicVars := state.GetNonBasicVariables()
	for ii, v := range nonBasicVars {
		// fmt.Println("Adding non-basic variable constraint for variable: ", v)
		// Find the index of the variable in the tableau
		vIdxInTableau, _ := symbolic.FindInSlice(v, state.Tableau.Variables)
		// fmt.Println("vIdxInTableau: ", vIdxInTableau)
		// fmt.Println("Targeted row: ", numConstraints+ii)
		AAugmented.Set(numConstraints+ii, vIdxInTableau, 1.0)
	}
	// b is already zero in the new rows, so we don't need to set anything in bAugmented

	// Solve the system of equations
	err = solutionVec.SolveVec(AAugmented, bAugmented)
	if err != nil {
		return mat.VecDense{}, fmt.Errorf("TableauAlgorithmState: Failed to solve for optimal solution (%v)", err)
	}

	return *solutionVec, nil
}

func (state *TableauAlgorithmState) CreateOptimalValuesMap(originalVariablesAsStandardFormExpressions map[symbolic.Variable]symbolic.Expression) (map[uint64]float64, error) {
	// Input Checking
	err := state.Check()
	if err != nil {
		return nil, err
	}

	// Setup
	solutionVec, err := state.CalculateOptimalSolution()
	if err != nil {
		return nil, err
	}

	// Create the map between the STANDARD FORM variables and the optimal values
	standardFormOptimalValues := map[symbolic.Variable]symbolic.Expression{}
	for ii, v := range state.Tableau.Variables {
		standardFormOptimalValues[v] = symbolic.K(solutionVec.AtVec(ii))
	}

	// Create the map between the ORIGINAL variables and the optimal values
	optimalValueMap := map[uint64]float64{}
	for origVar, expr := range originalVariablesAsStandardFormExpressions {
		// Evaluate the expression (which is equal to the original variable) using the solution map
		value := expr.SubstituteAccordingTo(standardFormOptimalValues)
		valAsK, ok := value.(symbolic.K)
		if !ok {
			return nil, fmt.Errorf("TableauAlgorithmState: Failed to evaluate optimal value for original variable %v", origVar)
		}
		// Add the value to the output map
		optimalValueMap[origVar.ID] = float64(valAsK)
	}

	return optimalValueMap, nil
}

func (state *TableauAlgorithmState) ToSolution(
	condition tableau_termination.TerminationType,
	varMap map[symbolic.Variable]symbolic.Expression,
	originalProblem *problem.OptimizationProblem,
) (simplex_solution.SimplexSolution, error) {
	// Create container for solution
	var sol simplex_solution.SimplexSolution
	var err error

	// Construct Solution Status
	sol.Status = condition.ToOptimizationStatus()

	// Construct Iteration Count
	sol.Iterations = state.IterationCount

	// Attach original problem
	sol.OriginalProblem = originalProblem

	// Construct Variable map
	sol.VariableValues, err = state.CreateOptimalValuesMap(varMap)
	if err != nil {
		return sol,
			fmt.Errorf(
				"There was an issue creating the optimal values map at termination: %v",
				err,
			)
	}

	// Construct Objective Value
	sol.Objective, err = solution.GetOptimalObjectiveValue(&sol)
	if err != nil {
		return sol,
			fmt.Errorf(
				"There was an issue getting the objective value at termination: %v",
				err,
			)
	}

	// Assemble Solution Output
	return sol, nil
}
