package utils

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

// The Tableau Representation of a linear program in standard form
type Tableau struct {
	Variables             []symbolic.Variable
	BasicVariableIndicies []int      // The basic variables in order of their connection to the constraint rows
	AsCompressedMatrix    *mat.Dense // The compressed matrix contains all
}

/*
A
Description:

	Extracts the A matrix (linear equality constraint matrix) from the tableau.
*/
func (tableau *Tableau) A() *mat.Dense {
	// Check that tableau is valid
	err := tableau.Check()
	if err != nil {
		panic(err)
	}

	// Setup
	nTableauRows, nTableauCols := tableau.AsCompressedMatrix.Dims()
	compressedMatrixCopy := mat.NewDense(nTableauRows, nTableauCols, nil)
	compressedMatrixCopy.Copy(tableau.AsCompressedMatrix)

	var out *mat.Dense = mat.NewDense(nTableauRows-1, nTableauCols-1, nil)

	// Extract the values that we care about.
	for ii := 0; ii < nTableauRows-1; ii++ {
		newRowII := compressedMatrixCopy.RawRowView(ii + 1)
		out.SetRow(ii, newRowII[:len(newRowII)-1])
	}

	return out
}

/*
B
Description:

	Extracts the B vector (linear equality vector) from the tableau.
*/
func (tableau *Tableau) B() *mat.VecDense {
	// Check that tableau is valid
	err := tableau.Check()
	if err != nil {
		panic(err)
	}

	// Setup
	nTableauRows, nTableauCols := tableau.AsCompressedMatrix.Dims()
	compressedMatrixCopy := mat.NewDense(nTableauRows, nTableauCols, nil)
	compressedMatrixCopy.Copy(tableau.AsCompressedMatrix)

	// Extract the values that we care about
	var lhsAsSliceOfFloats []float64
	for ii := 0; ii < nTableauRows-1; ii++ {
		lhsAsSliceOfFloats = append(
			lhsAsSliceOfFloats,
			compressedMatrixCopy.At(1+ii, nTableauCols-1),
		)

	}

	return mat.NewVecDense(nTableauRows-1, lhsAsSliceOfFloats)
}

/*
ABasic
Description:

	Returns the matrix of coefficients of the basic variables (A_B) in the current tableau of
	the algorithm.
*/
func (tableau *Tableau) ABasic() (*mat.Dense, error) {
	// Check the state for validity
	err := tableau.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return ABasic
	}

	// Slice the A Matrix according to the basic variables
	A := tableau.A()
	ABasic, err := SliceMatrixAccordingToVariableSet(
		getKMatrix.From(A),
		tableau.Variables,
		tableau.BasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to slice A matrix for basic variables (%v)", err)
	}

	ABasicAsDense := ABasic.ToDense()

	return &ABasicAsDense, nil
}

/*
ANonBasic
Description:

	Returns the matrix of coefficients of the non-basic variables (A_N) in the current tableau.
*/
func (tableau *Tableau) ANonBasic() (*mat.Dense, error) {
	// Check the tableau for validity
	err := tableau.Check()
	if err != nil {
		return nil, err // Invalid tableau, cannot return ANonBasic
	}

	// Slice the A Matrix according to the non-basic variables
	A := tableau.A()
	ANonBasic, err := SliceMatrixAccordingToVariableSet(
		getKMatrix.From(A),
		tableau.Variables,
		tableau.NonBasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to slice A matrix for non-basic variables (%v)", err)
	}

	ANonBasicAsDense := ANonBasic.ToDense()

	return &ANonBasicAsDense, nil
}

/*
C
Description:

	Extracts the C vector (linear objective vector) from the tableau
*/
func (tableau *Tableau) C() *mat.VecDense {
	// Check that tableau is valid
	err := tableau.Check()
	if err != nil {
		panic(err)
	}

	// Setup
	nTableauRows, nTableauCols := tableau.AsCompressedMatrix.Dims()
	compressedMatrixCopy := mat.NewDense(nTableauRows, nTableauCols, nil)
	compressedMatrixCopy.Copy(tableau.AsCompressedMatrix)

	// Extract the values that we care about
	topRow := compressedMatrixCopy.RowView(0)
	topRowAsVecDense, _ := topRow.(*mat.VecDense)

	return mat.NewVecDense(nTableauCols-1, topRowAsVecDense.RawVector().Data)
}

/*
CBasic
Description:

	Returns the cost vector for the basic variables (c_B) in the tableau.
*/
func (tableau *Tableau) CBasic() (*mat.VecDense, error) {
	// Check the tableau for validity
	err := tableau.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return CBasic
	}

	// Slice the cost vector according to the basic variables
	cBasic, err := SliceVectorAccordingToVariableSet(
		getKVector.From(tableau.C()),
		tableau.Variables,
		tableau.BasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("[tableau algorithm state]: Failed to slice cost vector for basic variables (%v)", err)
	}

	cBasicAsVecDense := cBasic.ToVecDense()

	return &cBasicAsVecDense, nil
}

/*
Check
Description:

	This method checks whether or not the Tableau is well-defined.
	Specifically, we check:
	- BasicVariableIndicies are within the range [0, len(tableau.Variables)]
	- Tableau has:
		+ len(AllVariables) + 2 columns
*/
func (tableau *Tableau) Check() error {
	// Check the BasicVariableIndicies
	nVariables := len(tableau.Variables)
	for _, bvIndex := range tableau.BasicVariableIndicies {
		if (bvIndex < 0) || (bvIndex >= nVariables) {
			return fmt.Errorf(
				"the basic variable %v is outside of the expected range [0,%v]",
				bvIndex,
				nVariables-1,
			)
		}
	}

	// Check that the number of columns is equal to len(AllVariables) + 1
	_, nTableauCols := tableau.AsCompressedMatrix.Dims()
	if nTableauCols != len(tableau.Variables)+1 {
		return fmt.Errorf(
			"The number of columns in the tableau is %v; expected %v columns.",
			nTableauCols,
			len(tableau.Variables)+2,
		)
	}

	// All Checks passed
	return nil
}

/*
CNonBasic
Description:

	Returns the cost vector for the non-basic variables (c_N) in the tableau.
*/
func (tableau *Tableau) CNonBasic() (*mat.VecDense, error) {
	// Check the tableau for validity
	err := tableau.Check()
	if err != nil {
		return nil, err // Invalid state, cannot return CNonBasic
	}

	// Slice the cost vector according to the non-basic variables
	cNonBasic, err := SliceVectorAccordingToVariableSet(
		getKVector.From(tableau.C()),
		tableau.Variables,
		tableau.NonBasicVariables(),
	)
	if err != nil {
		return nil, fmt.Errorf("[tableau algorithm state]: Failed to slice cost vector for non-basic variables (%v)", err)
	}

	cNonBasicAsVecDense := cNonBasic.ToVecDense()

	return &cNonBasicAsVecDense, nil
}

func (tableau *Tableau) NonBasicVariableIndicies() []int {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		panic(
			fmt.Errorf("tableau provided to NonBasicVariablesIndicies() was invalid: %v", err),
		)
	}

	// Algorithm
	out := []int{}
	for ii := 0; ii < len(tableau.Variables); ii++ {
		if foundIdx, _ := symbolic.FindInSlice(ii, tableau.BasicVariableIndicies); foundIdx == -1 {
			out = append(out, ii)
		}
	}
	return out
}

func (tableau *Tableau) NonBasicVariables() []symbolic.Variable {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		panic(
			fmt.Errorf("tableau provided to NonBasicVariablesIndicies() was invalid: %v", err),
		)
	}

	// Compute output
	out := []symbolic.Variable{}
	for _, nbIndex := range tableau.NonBasicVariableIndicies() {
		out = append(out, tableau.Variables[nbIndex])
	}
	return out
}

func (tableau *Tableau) BasicVariables() []symbolic.Variable {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		panic(
			fmt.Errorf("tableau provided to NonBasicVariablesIndicies() was invalid: %v", err),
		)
	}

	// Compute output
	out := []symbolic.Variable{}
	for _, nbIndex := range tableau.BasicVariableIndicies {
		out = append(out, tableau.Variables[nbIndex])
	}
	return out
}

/*
GetInitialTableauFrom
Description:

	This function computes the initial tableau of an
*/
func GetInitialTableauFrom(problemIn *problem.OptimizationProblem) (Tableau, error) {
	// Input Processing
	if problemIn == nil {
		return Tableau{}, fmt.Errorf(
			"Check: tableau.Problem cannot be nil",
		)
	}

	// Ensure that the problem is a linear program
	if !problemIn.IsLinear() {
		return Tableau{}, fmt.Errorf(
			"Check: the problem is not a linear program",
		)
	}

	// Transform the problem into the standard form where all constraints
	// are equality constraints
	problemInStandardForm, slackVariables, err := problemIn.ToLPStandardForm1()
	if err != nil {
		return Tableau{}, err
	}

	// Transform SlackVariables object into indicies
	var slackVariableIndicies []int
	for _, slackVar := range slackVariables {
		foundIdx, _ := symbolic.FindInSlice(slackVar, problemInStandardForm.Variables)
		slackVariableIndicies = append(slackVariableIndicies, foundIdx)
	}

	// Create the matrix
	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		return Tableau{}, err
	}
	Ab := symbolic.HStack(A, b)

	objectiveExpression := problemInStandardForm.Objective.Expression.(symbolic.ScalarExpression)
	c := objectiveExpression.LinearCoeff(problemInStandardForm.Variables)

	var cExtended *mat.VecDense = mat.NewVecDense(c.Len()+1, nil)
	cExtended.CopyVec(&c)
	cExtended.SetVec(c.Len(), 0.0)

	tableauMatCondensed := symbolic.VStack(
		symbolic.VecDenseToKVector(*cExtended).Transpose(),
		Ab,
	)
	tableauMatCondensedAsDense := tableauMatCondensed.(symbolic.KMatrix).ToDense()

	// Create the tableau
	return Tableau{
		AsCompressedMatrix:    &tableauMatCondensedAsDense,
		Variables:             problemInStandardForm.Variables,
		BasicVariableIndicies: slackVariableIndicies,
	}, nil
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
	c := tableau.AsCompressedMatrix.RowView(0)

	// Check if all coefficients are less than or equal to zero
	for ii := 0; ii < c.Len(); ii++ {
		if c.AtVec(ii) > 0 {
			return false
		}
	}

	return true
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
func (tableau *Tableau) ComputeFeasibleSolution(xNonBasic *mat.VecDense) (*mat.VecDense, error) {
	// Setup
	fmt.Printf("Computing feasible solution...\n")
	fmt.Printf("Tableau: %v\n", tableau)
	nBasic := tableau.NumberOfBasicVariables()

	// Collect the matrices of coefficients
	A, b := tableau.A(), tableau.B()

	// Create the matrix of coefficients of the basic variables
	N, err := tableau.ANonBasic()
	if err != nil {
		return nil, err
	}

	B, err := tableau.ABasic()
	if err != nil {
		return nil, err
	}

	// Solve the system of equations
	x := mat.NewVecDense(len(tableau.BasicVariables()), nil)

	// Compute the part that comes from the rhs (b)
	// xComponentFromb  = B^-1 * b
	var BInv *mat.Dense = mat.NewDense(nBasic, nBasic, nil)

	fmt.Printf("A: %v\n", A)
	fmt.Printf("BAsDense: %v\n", B)
	err = BInv.Inverse(B)
	if err != nil {
		return nil, fmt.Errorf("there was an issue inverting the matrix: %v", err)
	}
	fmt.Printf("BInv: %v\n", BInv)

	x.MulVec(BInv, b)

	fmt.Printf("x: %v\n", x)
	fmt.Printf("N: %v\n", N)
	fmt.Printf("vNonBasic: %v\n", tableau.NonBasicVariables())

	// Compute the part that comes from the non-basic variables
	// xComponentFromXNonBasic = B^(-1) * N * x
	xComponentFromXNonBasic := mat.NewVecDense(len(tableau.BasicVariables()), nil)
	BN := mat.NewDense(tableau.NumberOfBasicVariables(), tableau.NumberOfNonBasicVariables(), nil)
	BN.Mul(BInv, N)
	nr, nc := BN.Dims()
	fmt.Println("BN.Dims() =", nr, nc)
	xComponentFromXNonBasic.MulVec(BN, xNonBasic)
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
	objectiveExpression := getKVector.From(tableau.C()).Transpose().Multiply(
		symbolic.VariableVector(tableau.Variables),
	)

	objectiveSE := objectiveExpression.(symbolic.ScalarExpression)

	c := objectiveSE.LinearCoeff(tableau.BasicVariables())

	return &c, nil
}

/*
NumberOfBasicVariables
Description:

	Returns the number of basic variables in the tableau.
*/
func (tableau *Tableau) NumberOfBasicVariables() int {
	return len(tableau.BasicVariables())
}

/*
NumberOfNonBasicVariables
Description:

	Returns the number of non-basic variables in the tableau.
*/
func (tableau *Tableau) NumberOfNonBasicVariables() int {
	return len(tableau.NonBasicVariables())
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
	objectiveExpression := getKVector.From(tableau.C()).Transpose().Multiply(
		symbolic.VariableVector(tableau.Variables),
	)
	objectiveSE, _ := objectiveExpression.(symbolic.ScalarExpression)

	c := objectiveSE.LinearCoeff(tableau.NonBasicVariables())

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
