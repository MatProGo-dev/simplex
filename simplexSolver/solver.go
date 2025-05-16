package simplexSolver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
)

type SimplexSolver struct {
	OriginalProblem                 *problem.OptimizationProblem
	ProblemWithAllPositiveVariables *problem.OptimizationProblem
	ProblemInStandardForm           *problem.OptimizationProblem
}

func New(name string) SimplexSolver {
	// Create name for the base problem
	baseProblemName := name + " Problem"
	return SimplexSolver{
		OriginalProblem:                 problem.NewProblem(baseProblemName + " (Original Problem)"),
		ProblemWithAllPositiveVariables: problem.NewProblem(baseProblemName + " (With All Positive Variables)"),
		ProblemInStandardForm:           problem.NewProblem(baseProblemName + " (In Standard Form)"),
	}
}

func For(problem *problem.OptimizationProblem) SimplexSolver {
	// Create a new solver
	solver := New(problem.Name + " Solver")

	// Set the original problem
	original := problem
	original.Name = problem.Name + " (Original Problem)"
	solver.OriginalProblem = original

	// Transform the problem to have all positive variables
	solver.TransformAllUnboundedVariables()

	// TODO: Transform the problem to standard form

	return solver
}

func (solver *SimplexSolver) FindAllBasicSolutionsForRank(m int) [][]int {
	//
	return [][]int{}
}

func (solver *SimplexSolver) TransformAllUnboundedVariables() error {
	var err error
	solver.ProblemWithAllPositiveVariables, err = solver.OriginalProblem.ToProblemWithAllPositiveVariables()
	if err != nil {
		return err
	}

	// Return nil
	return nil
}
