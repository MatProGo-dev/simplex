package algorithms

import (
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
)

type AlgorithmInternalState struct {
	BasicVariables    []symbolic.Variable
	NonBasicVariables []symbolic.Variable
	NonBasicValues    *mat.VecDense
	IterationCount    int
}

/*
NumberOfBasicVariables
Description:

	Returns the number of basic variables in the Internal State.
*/
func (state *AlgorithmInternalState) NumberOfBasicVariables() int {
	return len(state.BasicVariables)
}

/*
NumberOfNonBasicVariables
Description:

	Returns the number of non-basic variables in the Internal State.
*/
func (state *AlgorithmInternalState) NumberOfNonBasicVariables() int {
	return len(state.NonBasicVariables)
}
