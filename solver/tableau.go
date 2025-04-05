package solver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
)

func ToTableau(problem *problem.OptimizationProblem) symbolic.KMatrix {
	// Verify that this is a linear program
	if !problem.IsLinear() {
		panic("The problem is not linear.")
	}

	return nil
}
