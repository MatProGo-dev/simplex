package simplexSolver

import (
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

type SimplexSolverInternalState struct {
	BasicVariables    []symbolic.Variable
	NonBasicVariables []symbolic.Variable
	NonBasicValues    *mat.VecDense
	IterationCount    int
}
