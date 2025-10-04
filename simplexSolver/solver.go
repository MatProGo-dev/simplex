package simplexSolver

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/simplex/algorithms"
	tableau_algorithm1 "github.com/MatProGo-dev/simplex/algorithms/tableau"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
)

type SimplexSolver struct {
	Name           string
	IterationLimit int
	Algorithm      algorithms.AlgorithmType
}

func New(name string) SimplexSolver {
	return SimplexSolver{
		Name:           name,
		IterationLimit: 100,
		Algorithm:      algorithms.TypeNaiveTableau,
	}
}

func (solver *SimplexSolver) CreateAlgorithm(algoType algorithms.AlgorithmType) (algorithms.AlgorithmInterface, error) {
	// Setup

	// Selection Logic
	switch algoType {
	case algorithms.TypeNaiveTableau:
		return &tableau_algorithm1.TableauAlgorithm{
			IterationLimit: solver.IterationLimit,
		}, nil
	default:
		return &tableau_algorithm1.TableauAlgorithm{}, fmt.Errorf(
			"The Solve() function was given an unknown solver type: %v",
			algoType,
		)
	}
}

func (solver *SimplexSolver) Solve(prob problem.OptimizationProblem) (simplex_solution.SimplexSolution, error) {
	// Setup

	// Choose Algorithm
	algo, err := solver.CreateAlgorithm(solver.Algorithm)
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf(
			"The Solve() function was given an unknown solver type: %v",
			solver.Algorithm,
		)
	}

	// Apply algorithm
	return algo.Solve(prob)

}
