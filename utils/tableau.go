package utils

import (
	"fmt"
	"math"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

// The Tableau Representation of a linear program in standard form
// Specifically, the tableau represents the problem:
//
//	minimize 		c^T * x + d
//	subject to 		A * x = b
//					x >= 0
//
// It represents the problem by storing the following information:
// - The list of all variables in the problem (including slack variables)
// - The indicies of the basic variables in the list of all variables
// - A Matrix representing the tableau as follows
//
//	| c^T | d |
//	|  A  | b |
type Tableau struct {
	Variables             []symbolic.Variable
	BasicVariableIndicies []int      // The basic variables in order of their connection to the constraint rows
	AsCompressedMatrix    *mat.Dense // The compressed matrix contains all of the information
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

	// Return the C vector (excluding the last entry which is the constant term)
	out := mat.NewVecDense(nTableauCols-1, nil)
	out.CopyVec(topRowAsVecDense.SliceVec(0, nTableauCols-1))

	return out
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

func (tableau *Tableau) D() float64 {
	// Check that tableau is valid
	err := tableau.Check()
	if err != nil {
		panic(err)
	}

	// Return the value of d
	_, nTableauCols := tableau.AsCompressedMatrix.Dims()
	return tableau.AsCompressedMatrix.At(0, nTableauCols-1)
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

func (tableau *Tableau) NumberOfConstraints() int {
	nRows, _ := tableau.AsCompressedMatrix.Dims()
	return nRows - 1
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
	problemInStandardForm, slackVariables, err := problemIn.ToLPStandardForm2()
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
	// [ A | b ]
	A, b, err := problemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		return Tableau{}, err
	}
	Ab := symbolic.HStack(A, b)

	// Create the cost vector
	// [ -c^T | -d ]
	objectiveExpression := problemInStandardForm.Objective.Expression.(symbolic.ScalarExpression)
	c := objectiveExpression.LinearCoeff(problemInStandardForm.Variables)
	d := objectiveExpression.Constant()

	c.ScaleVec(-1.0, &c) // We negate c because we are converting from a maximization to a minimization problem

	var cExtended *mat.VecDense = mat.NewVecDense(c.Len()+1, nil)
	cExtended.CopyVec(&c)
	cExtended.SetVec(c.Len(), -d)

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
CanNotBeImproved
Description:

	Returns true if the tableau's objective row indicates that
	the objective function can not be improved by changing
	any of the non-basic variables.
*/
func (tableau *Tableau) CanNotBeImproved() bool {
	// Get the coefficients of the non-basic variables
	c := tableau.C()

	// Check if all coefficients are less than or equal to zero
	for ii := 0; ii < c.Len(); ii++ {
		if c.AtVec(ii) < 0 {
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
Pivot
Description:

	Performs a pivot operation on the tableau.
	This operation produces a new tableau where the entering variable
	becomes a basic variable and the exiting variable becomes a non-basic variable.
*/
func (tableau *Tableau) Pivot(enteringVarIdx int, exitingVarIdx int) (Tableau, error) {
	// Input Processing
	err := tableau.Check()
	if err != nil {
		return Tableau{}, fmt.Errorf("Pivot: %v", err)
	}

	// Check that the entering variable is not already a basic variable
	if foundIdx, _ := symbolic.FindInSlice(enteringVarIdx, tableau.BasicVariableIndicies); foundIdx != -1 {
		return Tableau{}, fmt.Errorf("Pivot: the entering variable is already a basic variable")
	}

	// Check that the exiting variable is a basic variable
	if foundIdx, _ := symbolic.FindInSlice(exitingVarIdx, tableau.BasicVariableIndicies); foundIdx == -1 {
		return Tableau{}, fmt.Errorf("Pivot: the exiting variable is not a basic variable")
	}

	// Setup
	nRows, nCols := tableau.AsCompressedMatrix.Dims()
	newTableauMat := mat.NewDense(nRows, nCols, nil)
	newTableauMat.Copy(tableau.AsCompressedMatrix)

	exitingConstraintIdx, err := symbolic.FindInSlice(exitingVarIdx, tableau.BasicVariableIndicies)
	if err != nil {
		return Tableau{}, fmt.Errorf("Pivot: could not find exiting variable in basic variable indicies (%v)", err)
	}

	// Perform the pivot operation
	// - "Normalize" the pivot element (i.e., the element at A[exitingVarIdx, enteringVarIdx])
	normalizingFactor := 1.0 / tableau.AsCompressedMatrix.At(exitingConstraintIdx+1, enteringVarIdx)
	newPivotRow := mat.NewVecDense(nCols, nil)
	for jj := 0; jj < nCols; jj++ {
		newPivotRow.SetVec(jj, tableau.AsCompressedMatrix.At(exitingConstraintIdx+1, jj)*normalizingFactor)
	}
	newTableauMat.SetRow(exitingConstraintIdx+1, newPivotRow.RawVector().Data)

	// - Zero out the other entries in the entering variable column
	for ii := 0; ii < nRows; ii++ {
		// Skip the pivot row (i.e., the row of the entering variable)
		if ii == exitingConstraintIdx+1 {
			continue
		}
		// Skip any row that is already zero
		if math.Abs(newTableauMat.At(ii, enteringVarIdx)) < 1e-14 {
			continue
		}

		// Otherwise, determine the factor needed to zero out this entry
		factorToZeroOut := -1.0 * newTableauMat.At(ii, enteringVarIdx)
		rowToAdd := mat.NewVecDense(nCols, nil)
		rowToAdd.ScaleVec(factorToZeroOut, newPivotRow)
		rowToAdd.AddVec(rowToAdd, newTableauMat.RowView(ii))

		// Update the row in the new tableau
		newTableauMat.SetRow(ii, rowToAdd.RawVector().Data)
	}

	// Update the list of basic variable indicies
	newBasicVariableIndicies := make([]int, len(tableau.BasicVariableIndicies))
	copy(newBasicVariableIndicies, tableau.BasicVariableIndicies)
	for ii, bvIdx := range tableau.BasicVariableIndicies {
		if bvIdx == exitingVarIdx {
			newBasicVariableIndicies[ii] = enteringVarIdx
			break
		}
	}

	// Create the new tableau
	newTableau := Tableau{
		Variables:             tableau.Variables,
		BasicVariableIndicies: newBasicVariableIndicies,
		AsCompressedMatrix:    newTableauMat,
	}

	// Check the new tableau for validity
	err = newTableau.Check()
	if err != nil {
		return Tableau{}, fmt.Errorf("Pivot: the resulting tableau is invalid (%v)", err)
	}

	return newTableau, nil
}
