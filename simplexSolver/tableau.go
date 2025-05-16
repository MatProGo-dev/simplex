package simplexSolver

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
	A, b, err := tableau.Problem.LinearEqualityConstraintMatrices()
	if err != nil {
		return nil, err
	}

	// Create the matrix of coefficients of the basic variables
	N, err := SliceMatrixAccordingToVariableSet(
		tableau.Problem,
		A,
		tableau.NonBasicVariables,
	)
	if err != nil {
		return nil, err
	}

	B, err := SliceMatrixAccordingToVariableSet(
		tableau.Problem,
		A,
		tableau.BasicVariables,
	)
	if err != nil {
		return nil, err
	}

	// Solve the system of equations
	x := mat.NewVecDense(len(tableau.BasicVariables), nil)

	// Compute the part that comes from the rhs (b)
	// xComponentFromb  = B^-1 * b
	var BInv *mat.Dense
	BAsDense := B.ToDense()
	fmt.Printf("BAsDense: %v\n", BAsDense)
	BInv.Inverse(&BAsDense)
	bAsVecDense := b.ToVecDense()
	x.MulVec(BInv, &bAsVecDense)

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

func ToTableau(problem *problem.OptimizationProblem) symbolic.KMatrix {
	// Verify that this is a linear program
	if !problem.IsLinear() {
		panic("The problem is not linear.")
	}

	// Transform all non-zero variables to

	return nil
}
