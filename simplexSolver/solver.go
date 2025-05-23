package simplexSolver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

type SimplexSolverInternalState struct {
	BasicVariables        []symbolic.Variable
	NonBasicVariables     []symbolic.Variable
	CurrentNonBasicValues *mat.VecDense
	IterationCount        int
}

type SimplexSolver struct {
	OriginalProblem       *problem.OptimizationProblem
	ProblemInStandardForm *problem.OptimizationProblem
	State                 SimplexSolverInternalState
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
	solver.ProblemInStandardForm, slackVariables, err := solver.OriginalProblem.ToLPStandardForm1()
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
		CurrentNonBasicValues: mat.NewVecDense(len(solver.ProblemInStandardForm.Variables), nil),
		IterationCount:        0,
	}

	// TODO: Transform the problem to standard form

	return solver, nil
}

func (solver *SimplexSolver) GetStateAsTableau() [][]int {
	//
	return [][]int{}
}
