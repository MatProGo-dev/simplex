package simplexSolver

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"matprogo.dev/solvers/simplex/algorithms"
	"matprogo.dev/solvers/simplex/utils"
)

type SimplexSolver struct {
	OriginalProblem       *problem.OptimizationProblem
	ProblemInStandardForm *problem.OptimizationProblem
	InternalState         algorithms.AlgorithmInternalState
	config                Configuration
}

func New(name string) SimplexSolver {
	// Create name for the base problem
	baseProblemName := name + " Problem"
	return SimplexSolver{
		OriginalProblem:       problem.NewProblem(baseProblemName + " (Original Problem)"),
		ProblemInStandardForm: problem.NewProblem(baseProblemName + " (In Standard Form)"),
	}
}

func For(problem *problem.OptimizationProblem, configIn Configuration) (SimplexSolver, error) {
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

	// Initialize Internal Solver State
	solver.InternalState, err = solver.InitializeInternalState(slackVariables)
	if err != nil {
		return solver, err
	}

	// Configure the solver
	solver.config = configIn

	return solver, nil
}

// func (solver *SimplexSolver) CurrentStateToTableau() (symbolic.KMatrix, error) {

// }

func (solver *SimplexSolver) InitializeInternalState(initialSlackVariables []symbolic.Variable) (algorithms.AlgorithmInternalState, error) {
	// Setup

	// Initialize the internal state
	out := algorithms.AlgorithmInternalState{
		BasicVariables: initialSlackVariables,
		NonBasicVariables: utils.SetDifferenceOfVariables(
			solver.ProblemInStandardForm.Variables,
			initialSlackVariables,
		),
		IterationCount: 0,
	}

	nNonBasicVariables := out.NumberOfNonBasicVariables()
	out.NonBasicValues = mat.NewVecDense(
		nNonBasicVariables,
		nil,
	)

	return out, nil
}

func (solver *SimplexSolver) CreateAlgorithm(algoType algorithms.AlgorithmType) (algorithms.AlgorithmInterface, error) {
	// Setup

	// Selection Logic
	switch algoType {
	case algorithms.TypeNaive:
		return &algorithms.NaiveAlgorithm{
			ProblemInStandardForm: solver.ProblemInStandardForm,
			IterationLimit:        solver.config.IterationLimit,
		}, nil
	default:
		return &algorithms.NaiveAlgorithm{}, fmt.Errorf(
			"The Solve() function was given an unknown solver type: %v",
			algoType,
		)
	}
}

func (solver *SimplexSolver) Solve(algoType algorithms.AlgorithmType) (problem.Solution, error) {
	// Setup

	// Choose Algorithm
	algo, err := solver.CreateAlgorithm(algoType)
	if err != nil {
		return problem.Solution{}, fmt.Errorf(
			"The Solve() function was given an unknown solver type: %v",
			algoType,
		)
	}

	// Apply algorithm
	return algo.Solve(solver.InternalState)

}
