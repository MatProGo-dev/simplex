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
		state.GetNonBasicVariables(),
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
		state.GetNonBasicVariables(),
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

/*
ComputeSolutionFromState
Description:

	Computes the many components of the solution from the current state of the algorithm.
	It computes:
	- The value of the objective function
	- The values of the basic variables
	- The values of the non-basic variables
*/
func (algo *StanfordAlgorithm) ComputeSolutionFromState(state StanfordAlgorithmState) (problem.Solution, error) {
	// Setup
	fmt.Printf("Computing solution from state...\n")
	solution := problem.Solution{
		Status: problem.OptimizationStatus_OPTIMAL,
	}

	// Compute the feasible solution of the basic variables
	xBasic, err := algo.ComputeFeasibleBasicSolution(state)
	if err != nil {
		return solution, fmt.Errorf("StanfordAlgorithm: Failed to compute feasible basic solution (%v)", err)
	}

	// Compute the value of the objective function
	objValue, err := algo.ComputeObjectiveFunctionValueWithFeasibleBasicSolution(state, xBasic)
	if err != nil {
		return solution, fmt.Errorf("StanfordAlgorithm: Failed to compute objective function value (%v)", err)
	}
	solution.Objective = objValue

	// Set the values of the basic variables
	for ii, bv := range state.BasicVariables {
		solution.Values[bv.ID] = xBasic.AtVec(ii)
	}

	// Set the values of the non-basic variables
	for jj, nv := range state.GetNonBasicVariables() {
		solution.Values[nv.ID] = state.NonBasicValues.AtVec(jj)
	}

	return solution, nil
}

func (algo *StanfordAlgorithm) Solve(initialState StanfordAlgorithmState) (problem.Solution, error) {
	// Setup
	var stateII StanfordAlgorithmState = initialState

	// Compute the feasible solution for the current choice of
	// basic variables
	for iter := 0; iter < algo.IterationLimit; iter++ {
		// Test for Termination
		r, err := stateII.GetReducedCostVector()
		if err != nil {
			return problem.Solution{}, fmt.Errorf(
				"StanfordAlgorithm: Failed to get reduced cost vector (%v) at iteration #%v",
				err,
				iter,
			)
		}

		// Find the entering variable (most negative reduced cost)
		minReducedCost := 0.0
		enteringVarIndex := -1
		for ii := 0; ii < r.Len(); ii++ {
			if r.AtVec(ii) < minReducedCost {
				minReducedCost = r.AtVec(ii)
				enteringVarIndex = ii
			}
		}

		if enteringVarIndex == -1 {
			// No entering variable found, the solution is optimal
			fmt.Printf("Optimal solution found at iteration %d\n", iter)
			return algo.ComputeSolutionFromState(stateII)
		}

		// Check to see if vector (ABasic)^(-1) * A 's e-th vector contains a positive entry
		var ABasicInv mat.Dense
		ABasic, err := stateII.ABasic()
		if err != nil {
			return problem.Solution{}, fmt.Errorf("StanfordAlgorithm: Failed to get ABasic matrix (%v)", err)
		}
		ABasicInv.Inverse(ABasic)
		var ABasicAe mat.VecDense
		ABasicAe.MulVec(&ABasicInv, stateII.A.ColView(enteringVarIndex))

		objectiveIsUnboundedBelow := true
		for ii := 0; ii < ABasicAe.Len(); ii++ {
			if ABasicAe.AtVec(ii) > 0.0 {
				// Found a positive entry, we can proceed with the pivot operation
				fmt.Printf("Found a positive entry in ABasic^(-1) * A at index %d\n", ii)
				objectiveIsUnboundedBelow = false
			}
		}
		if objectiveIsUnboundedBelow {
			return problem.Solution{}, fmt.Errorf(
				"StanfordAlgorithm: Objective function is unbounded below, no solution exists at iteration %d",
				iter,
			)
		}

		// Compute the minimum ratio test

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
