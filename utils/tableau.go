package utils

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

type Tableau struct {
	// The Tableau Representation of a linear program in standard form
	BasicVariables    []symbolic.Variable
	NonBasicVariables []symbolic.Variable
	NonBasicValues    *mat.VecDense
	Problem           *problem.OptimizationProblem
}

func GetInitialTableau(problemIn *problem.OptimizationProblem) (Tableau, error) {
	// Setup

	// Transform the problem into the standard form where all constraints
	// are equality constraints
	problemInStandardForm, slackVariables, err := problemIn.ToLPStandardForm1()
	if err != nil {
		return Tableau{}, err
	}

	// Create the tableau
	tableau := Tableau{
		BasicVariables:    slackVariables, // The slack variables are the initial basic variables
		NonBasicVariables: []symbolic.Variable{},
		Problem:           problemInStandardForm,
	}

	// The non-basic variables are the original variables
	tableau.NonBasicVariables = SetDifferenceOfVariables(
		problemInStandardForm.Variables,
		slackVariables,
	)

	fmt.Printf("Basic Variables: %v\n", tableau.BasicVariables)
	fmt.Printf("Non-Basic Variables: %v\n", tableau.NonBasicVariables)

	// The non-basic values are assumed to be zero
	tableau.NonBasicValues = mat.NewVecDense(len(tableau.NonBasicVariables), nil)

	return tableau, nil
}

/*
AllObjectiveRowEntriesAreLessThanOrEqualToZero
Description:

	Checks if all entries in the objective row of the tableau
	are less than or equal to zero.

	This is a necessary condition for the tableau to be optimal.
*/
func (tableau *Tableau) AllObjectiveRowEntriesAreLessThanOrEqualToZero() bool {
	// Get the coefficients of the non-basic variables
	pLike, err := symbolic.ToPolynomialLikeScalar(tableau.Problem.Objective.Expression)
	if err != nil {
		return false
	}

	c := pLike.LinearCoeff(tableau.NonBasicVariables)

	// Check if all coefficients are less than or equal to zero
	for ii := 0; ii < c.Len(); ii++ {
		if c.AtVec(ii) > 0 {
			return false
		}
	}

	return true
}

/*
Check
Description:

	Checks the validity of the tableau.
	It ensures that the tableau has a valid problem,
	that the problem is a linear program,
	and that the basic and non-basic variables are not empty.
*/
func (tableau *Tableau) Check() error {
	// Input Processing
	if tableau.Problem == nil {
		return fmt.Errorf(
			"Check: tableau.Problem cannot be nil",
		)
	}

	// Ensure that the problem is a linear program
	if !tableau.Problem.IsLinear() {
		return fmt.Errorf(
			"Check: the problem is not a linear program",
		)
	}

	if len(tableau.BasicVariables) == 0 {
		return fmt.Errorf("Check: tableau.BasicVariables cannot be empty")
	}

	if len(tableau.NonBasicVariables) == 0 {
		return fmt.Errorf("Check: tableau.NonBasicVariables cannot be empty")
	}

	// if tableau.NonBasicValues == nil {
	// 	return fmt.Errorf("Check: tableau.NonBasicValues cannot be nil")
	// }

	return nil
}

/*
ComputeFeasibleSolution
Description:

	Computes a feasible solution of the BASIC variables
	by solving the system of equations
		A * x = b
	Where A is the matrix of coefficients of the basic variables
	and b is the vector of constants.
*/
func (tableau *Tableau) ComputeFeasibleSolution() (*mat.VecDense, error) {
	// Setup
	fmt.Printf("Computing feasible solution...\n")
	fmt.Printf("Problem: %v\n", tableau.Problem)
	fmt.Printf("Tableau: %v\n", tableau)
	nBasic := tableau.NumberOfBasicVariables()

	// Collect the matrices of coefficients
	A, b, err := tableau.Problem.LinearEqualityConstraintMatrices()
	if err != nil {
		return nil, err
	}

	// Create the matrix of coefficients of the basic variables
	N, err := SliceMatrixAccordingToVariableSet(
		A,
		tableau.Problem.Variables,
		tableau.NonBasicVariables,
	)
	if err != nil {
		return nil, err
	}

	B, err := SliceMatrixAccordingToVariableSet(
		A,
		tableau.Problem.Variables,
		tableau.BasicVariables,
	)
	if err != nil {
		return nil, err
	}

	// Solve the system of equations
	x := mat.NewVecDense(len(tableau.BasicVariables), nil)

	// Compute the part that comes from the rhs (b)
	// xComponentFromb  = B^-1 * b
	var BInv *mat.Dense = mat.NewDense(nBasic, nBasic, nil)
	BAsDense := B.ToDense()
	fmt.Printf("A: %v\n", A)
	fmt.Printf("BAsDense: %v\n", BAsDense)
	err = BInv.Inverse(&BAsDense)
	if err != nil {
		return nil, fmt.Errorf("there was an issue inverting the matrix: %v", err)
	}
	fmt.Printf("BInv: %v\n", BInv)
	bAsVecDense := b.ToVecDense()
	fmt.Printf("bAsVecDense: %v\n", bAsVecDense)
	x.MulVec(BInv, &bAsVecDense)

	fmt.Printf("x: %v\n", x)

	// Compute the part that comes from the non-basic variables
	// xComponentFromXNonBasic = B^(-1) * N * x
	xComponentFromXNonBasic := mat.NewVecDense(len(tableau.BasicVariables), nil)
	BN := mat.NewDense(tableau.NumberOfBasicVariables(), tableau.NumberOfNonBasicVariables(), nil)
	NAsDense := N.ToDense()
	BN.Mul(BInv, &NAsDense)
	xComponentFromXNonBasic.MulVec(BN, tableau.NonBasicValues)
	x.AddVec(x, xComponentFromXNonBasic)

	return x, nil
}

func (tableau *Tableau) BasicVariableContributionToObjective() (*mat.VecDense, error) {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		return nil, fmt.Errorf("BasicVariableContributionToObjective: %v", err)
	}

	// Setup

	// Collect the coefficients of the basic variables
	pLike, err := symbolic.ToPolynomialLikeScalar(tableau.Problem.Objective.Expression)
	if err != nil {
		return nil, fmt.Errorf("BasicVariableContributionToObjective: %v", err)
	}

	c := pLike.LinearCoeff(tableau.BasicVariables)

	return &c, nil
}

/*
NumberOfBasicVariables
Description:

	Returns the number of basic variables in the tableau.
*/
func (tableau *Tableau) NumberOfBasicVariables() int {
	return len(tableau.BasicVariables)
}

/*
NumberOfNonBasicVariables
Description:

	Returns the number of non-basic variables in the tableau.
*/
func (tableau *Tableau) NumberOfNonBasicVariables() int {
	return len(tableau.NonBasicVariables)
}

/*
ToDense
Description:

	Converts the tableau to a dense matrix representation.
*/
func (tableau *Tableau) AsDense() (*mat.Dense, error) {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		return nil, fmt.Errorf("AsDense: %v", err)
	}

	// Setup
	allVars := append(tableau.BasicVariables, tableau.NonBasicVariables...)
	nVars := len(allVars)

	A, b, err := tableau.Problem.LinearEqualityConstraintMatrices()
	if err != nil {
		return nil, fmt.Errorf("AsDense: %v", err)
	}
	AAsDense := A.ToDense()
	bAsVecDense := b.ToVecDense()
	nRowsA, _ := AAsDense.Dims()

	// Create the dense matrix representation
	// [ 1 , -c^T, 0 ]
	// [ 0 , A,    b ]
	// where c is the vector of coefficients of the non-basic variables
	// and A is the matrix of coefficients of the basic variables.

	out := mat.NewDense(
		1+nRowsA, // 1 for the objective row, and nRowsA for the rest of the rows
		2+nVars,  // 1 for the first element, 1 for the constraint coefficient, and nVars for the basic and non-basic variables
		nil,
	)

	out.Set(0, 0, 1) // The first element is 1

	// Set the rest of the first row to be the coefficients
	// of both the basic and non-basic variables

	pLike, err := symbolic.ToPolynomialLikeScalar(tableau.Problem.Objective.Expression)
	if err != nil {
		return nil, fmt.Errorf("AsDense: %v", err)
	}

	c := pLike.LinearCoeff(allVars)
	for ii := 0; ii < c.Len(); ii++ {
		out.Set(0, ii+1, -c.AtVec(ii)) // The first row contains the negative coefficients
	}

	// Set the entries corresponding to the A matrix
	for ii := 0; ii < nRowsA; ii++ {
		for jj := 0; jj < nVars; jj++ {
			out.Set(ii+1, jj+1, AAsDense.At(ii, jj))
		}
	}

	// Set the entries corresponding to the b vector
	for ii := 0; ii < nRowsA; ii++ {
		out.Set(ii+1, nVars+1, bAsVecDense.AtVec(ii)) // The first column
	}

	return out, nil
}

/*
SelectPivotColumn
Description:

	Selects the pivot column based on the tableau.
	This is the column from the non-basic variables that leads to the
	largest increase in the objective function value.
	It returns the index of the pivot column in the tableau.
*/
func (tableau *Tableau) SelectPivotColumn() (int, error) {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		return -1, fmt.Errorf("SelectPivotColumn: %v", err)
	}

	// Setup
	nNonBasic := tableau.NumberOfNonBasicVariables()
	if nNonBasic == 0 {
		return -1, fmt.Errorf("SelectPivotColumn: there are no non-basic variables")
	}

	// Get the coefficients of the non-basic variables
	pLike, err := symbolic.ToPolynomialLikeScalar(tableau.Problem.Objective.Expression)
	if err != nil {
		return -1, fmt.Errorf("SelectPivotColumn: %v", err)
	}

	c := pLike.LinearCoeff(tableau.NonBasicVariables)

	// Find the index of the pivot column
	pivotColIndex := -1
	maxValue := 0.0
	for ii := 0; ii < nNonBasic; ii++ {
		if c.AtVec(ii) > maxValue {
			maxValue = c.AtVec(ii)
			pivotColIndex = ii
		}
	}

	return pivotColIndex, nil
}
