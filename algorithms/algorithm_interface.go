package algorithms

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
)

type AlgorithmInterface interface {
	// Solves the provided optimization problem.
	Solve(prob problem.OptimizationProblem) (simplex_solution.SimplexSolution, error)
}
