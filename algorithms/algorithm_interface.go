package algorithms

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
)

type AlgorithmInterface interface {
	// Solves the provided optimization problem.
	Solve(initialState AlgorithmInternalState) (problem.Solution, error)
}
