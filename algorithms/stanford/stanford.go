package stanford_algorithm1

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"matprogo.dev/solvers/simplex/utils"
)

type StanfordAlgorithm struct {
	ProblemInStandardForm *problem.OptimizationProblem
	IterationLimit        int
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
func (algo *StanfordAlgorithm) ComputeFeasibleBasicSolution(state StanfordAlgorithmState) (*mat.VecDense, error) {
	// Setup
	fmt.Printf("Computing feasible solution...\n")
	fmt.Printf("Problem: %v\n", algo.ProblemInStandardForm)
	fmt.Printf("Tableau: %v\n", algo)
	nBasic := state.NumberOfBasicVariables()

	// Collect the matrices of coefficients
	A, b, err := algo.ProblemInStandardForm.LinearEqualityConstraintMatrices()
	if err != nil {
		return nil, err
	}

	// Create the matrix of coefficients of the basic variables
	N, err := utils.SliceMatrixAccordingToVariableSet(
		A,
		algo.ProblemInStandardForm.Variables,
		state.NonBasicVariables(),
	)
	if err != nil {
		return nil, err
	}

	B, err := utils.SliceMatrixAccordingToVariableSet(
		A,
		algo.ProblemInStandardForm.Variables,
		state.BasicVariables,
	)
	if err != nil {
		return nil, err
	}

	// Solve the system of equations
	x := mat.NewVecDense(state.NumberOfBasicVariables(), nil)

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
	xComponentFromXNonBasic := mat.NewVecDense(state.NumberOfBasicVariables(), nil)
	BN := mat.NewDense(state.NumberOfBasicVariables(), state.NumberOfNonBasicVariables(), nil)
	NAsDense := N.ToDense()
	BN.Mul(BInv, &NAsDense)
	xComponentFromXNonBasic.MulVec(BN, state.NonBasicValues)
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
func (algo *StanfordAlgorithm) ComputeObjectiveFunctionValueWithFeasibleBasicSolution(state StanfordAlgorithmState, xBasic *mat.VecDense) (float64, error) {
	// Setup
	fmt.Printf("Computing objective function value...\n")
	allVars := algo.ProblemInStandardForm.Variables

	// Collect the linear coefficient of the objective function
	objectiveExpression := algo.ProblemInStandardForm.Objective.Expression
	objectiveAsSE, tf := objectiveExpression.(symbolic.ScalarExpression)
	if !tf {
		return 0.0, fmt.Errorf("the objective function is not a scalar expression")
	}

	// TODO(kwesi): Check that the objective function is a constant or not

	// Compute the linear coefficient of the objective function
	c := objectiveAsSE.LinearCoeff(allVars)

	// Split the coefficient into the basic and non-basic variables
	cB, err := utils.SliceVectorAccordingToVariableSet(
		symbolic.VecDenseToKVector(c),
		allVars,
		state.BasicVariables,
	)
	if err != nil {
		return 0.0, err
	}

	cN, err := utils.SliceVectorAccordingToVariableSet(
		symbolic.VecDenseToKVector(c),
		allVars,
		state.NonBasicVariables(),
	)
	if err != nil {
		return 0.0, err
	}

	// Compute the value of the objective function
	// f(x) = c_B^T * x_B + c_N^T * x_N
	z := cB.Transpose().Multiply(xBasic).Plus(
		cN.Transpose().Multiply(state.NonBasicValues),
	)

	zAsK, tf := z.(symbolic.K)
	if !tf {
		return 0.0, fmt.Errorf("the objective function is not a scalar expression")
	}

	return float64(zAsK), nil
}

func (algo *StanfordAlgorithm) Solve(initialState StanfordAlgorithmState) (problem.Solution, error) {
	// Setup
	var stateII StanfordAlgorithmState = initialState

	// Compute the feasible solution for the current choice of
	// basic variables
	for iter := 0; iter < algo.IterationLimit; iter++ {
		// Compute the feasible Solution of the Basic variables
		xBasicII, err := algo.ComputeFeasibleBasicSolution(stateII)
		if err != nil {
			return problem.Solution{}, err
		}

		// Compute the value of the objective function
		objII, err := algo.ComputeObjectiveFunctionValueWithFeasibleBasicSolution(stateII, xBasicII)
		if err != nil {
			return problem.Solution{}, err
		}
		fmt.Printf("Iteration %d: Basic Solution: %v, Objective Value: %f\n", iter, xBasicII, objII)

		// Update the state of the algorithm
		stateII = stateII

		// Check if the solution is optimal

		if iter == algo.IterationLimit {
			return problem.Solution{
				Objective: objII,
				Status:    problem.OptimizationStatus_ITERATION_LIMIT,
			}, nil
		}

	}

	return problem.Solution{}, nil
}
