package simplexSolver

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

type SimplexSolver struct {
	OriginalProblem       *problem.OptimizationProblem
	ProblemInStandardForm *problem.OptimizationProblem
	State                 SimplexSolverInternalState
	IterationLimit        int
}

func New(name string) SimplexSolver {
	// Create name for the base problem
	baseProblemName := name + " Problem"
	return SimplexSolver{
		OriginalProblem:       problem.NewProblem(baseProblemName + " (Original Problem)"),
		ProblemInStandardForm: problem.NewProblem(baseProblemName + " (In Standard Form)"),
	}
}

func For(problem *problem.OptimizationProblem) (SimplexSolver, error) {
	// Create a new solver
	solver := New(problem.Name + " Solver")

	// Set the original problem
	original := problem
	original.Name = problem.Name + " (Original Problem)"
	solver.OriginalProblem = original

	// Transform the problem into the standard form where all constraints
	// are equality constraints
	var err error
	var slackVariables []symbolic.Variable
	solver.ProblemInStandardForm, slackVariables, err = solver.OriginalProblem.ToLPStandardForm1()
	if err != nil {
		return solver, err
	}

	// Initialize the internal state
	solver.State = SimplexSolverInternalState{
		BasicVariables: slackVariables,
		NonBasicVariables: SetDifferenceOfVariables(
			solver.ProblemInStandardForm.Variables,
			slackVariables,
		),
		IterationCount: 0,
	}
	solver.State.NonBasicValues = mat.NewVecDense(solver.NumberOfNonBasicVariables(), nil)

	// Configure the solver
	solver.IterationLimit = 1000

	return solver, nil
}

/*
NumberOfBasicVariables
Description:

	Returns the number of basic variables in the tableau.
*/
func (solver *SimplexSolver) NumberOfBasicVariables() int {
	return len(solver.State.BasicVariables)
}

/*
NumberOfNonBasicVariables
Description:

	Returns the number of non-basic variables in the tableau.
*/
func (solver *SimplexSolver) NumberOfNonBasicVariables() int {
	return len(solver.State.NonBasicVariables)
}

/*
ComputeFeasibleBasicSolution
Description:

	Computes a feasible solution of the BASIC variables
	of the optimization problem:
	maximize 		c^T * x
	subject to 		A * x = b
					x >= 0
	The introduction of basic and non-basic variables allows us to
	rewrite the problem as:
		maximize 		c_B^T * x_B + c_N^T * x_N
		subject to 		A_B * x_B + A_N * x_N = b
					x_B >= 0
					x_N >= 0
	Where
		A_B is the matrix of coefficients of the basic variables,
		A_N is the matrix of coefficients of the non-basic variables,
		c_B is the vector of coefficients of the basic variables,
		c_N is the vector of coefficients of the non-basic variables,
		x_B is the vector of basic variables,
		x_N is the vector of non-basic variables,
		b is the vector of constants.

	We assume that the value of the non-basic variables is given
	(i.e., they are already saved in the state).
*/
func (solver *SimplexSolver) ComputeFeasibleBasicSolution() (*mat.VecDense, error) {
	// Setup
	fmt.Printf("Computing feasible solution...\n")
	fmt.Printf("Problem: %v\n", solver.ProblemInStandardForm)
	fmt.Printf("Tableau: %v\n", solver)
	nBasic := solver.NumberOfBasicVariables()

	// Collect the matrices of coefficients
	A, b, err := solver.ProblemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		return nil, err
	}

	// Create the matrix of coefficients of the basic variables
	N, err := SliceMatrixAccordingToVariableSet(
		A,
		solver.ProblemInStandardForm.Variables,
		solver.State.NonBasicVariables,
	)
	if err != nil {
		return nil, err
	}

	B, err := SliceMatrixAccordingToVariableSet(
		A,
		solver.ProblemInStandardForm.Variables,
		solver.State.BasicVariables,
	)
	if err != nil {
		return nil, err
	}

	// Solve the system of equations
	x := mat.NewVecDense(solver.NumberOfBasicVariables(), nil)

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
	xComponentFromXNonBasic := mat.NewVecDense(solver.NumberOfBasicVariables(), nil)
	BN := mat.NewDense(solver.NumberOfBasicVariables(), solver.NumberOfNonBasicVariables(), nil)
	NAsDense := N.ToDense()
	BN.Mul(BInv, &NAsDense)
	xComponentFromXNonBasic.MulVec(BN, solver.State.NonBasicValues)
	x.AddVec(x, xComponentFromXNonBasic)

	return x, nil
}

/*
ComputeObjectiveFunctionValueWithFeasibleBasicSolution
Description:

	Computes the value of the objective function
	of the optimization problem:
	maximize 		c^T * x
	subject to 		A * x = b
					x >= 0
	when the feasible solution of the BASIC variables is given as
	input and the value of the non-basic variables is given
	in the state of the solver.
*/
func (solver *SimplexSolver) ComputeObjectiveFunctionValueWithFeasibleBasicSolution(xBasic *mat.VecDense) (float64, error) {
	// Setup
	fmt.Printf("Computing objective function value...\n")
	allVars := solver.ProblemInStandardForm.Variables

	// Collect the linear coefficient of the objective function
	objectiveExpression := solver.ProblemInStandardForm.Objective.Expression
	objectiveAsSE, tf := objectiveExpression.(symbolic.ScalarExpression)
	if !tf {
		return 0.0, fmt.Errorf("the objective function is not a scalar expression")
	}

	// TODO(kwesi): Check that the objective function is a constant or not

	// Compute the linear coefficient of the objective function
	c := objectiveAsSE.LinearCoeff(allVars)

	// Split the coefficient into the basic and non-basic variables
	cB, err := SliceVectorAccordingToVariableSet(
		symbolic.VecDenseToKVector(c),
		allVars,
		solver.State.BasicVariables,
	)
	if err != nil {
		return 0.0, err
	}

	cN, err := SliceVectorAccordingToVariableSet(
		symbolic.VecDenseToKVector(c),
		allVars,
		solver.State.NonBasicVariables,
	)
	if err != nil {
		return 0.0, err
	}

	// Compute the value of the objective function
	// f(x) = c_B^T * x_B + c_N^T * x_N
	z := cB.Transpose().Multiply(xBasic).Plus(
		cN.Transpose().Multiply(solver.State.NonBasicValues),
	)

	zAsK, tf := z.(symbolic.K)
	if !tf {
		return 0.0, fmt.Errorf("the objective function is not a scalar expression")
	}

	return float64(zAsK), nil
}

func (solver *SimplexSolver) CurrentStateToTableau() (symbolic.KMatrix, error) {

}

func (solver *SimplexSolver) Solve() (problem.Solution, error) {
	// Setup
	solver.State.IterationCount = 0

	// Compute the feasible solution for the current choice of
	// basic variables
	for iter := 0; iter < solver.IterationLimit; iter++ {
		// Compute the feasible Solution of the Basic variables
		xBasicII, err := solver.ComputeFeasibleBasicSolution()
		if err != nil {
			return problem.Solution{}, err
		}

		// Compute the value of the objective function
		objII, err := solver.ComputeObjectiveFunctionValueWithFeasibleBasicSolution(xBasicII)
		if err != nil {
			return problem.Solution{}, err
		}
		fmt.Printf("Iteration %d: Basic Solution: %v, Objective Value: %f\n", iter, xBasicII, objII)

		// Check if the solution is optimal

	}

	return problem.Solution{}, nil

}
