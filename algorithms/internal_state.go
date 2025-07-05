package algorithms

import (
	"fmt"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"gonum.org/v1/gonum/mat"
	"matprogo.dev/solvers/simplex/utils"
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

/*
ToTableau
Description:

	Converts the internal state to a tableau representation.
*/
func (state *AlgorithmInternalState) ToTableau(standardForm *problem.OptimizationProblem) (utils.Tableau, error) {
	// Input Processing
	if standardForm == nil {
		return utils.Tableau{}, fmt.Errorf(
			"ToTableau: standardForm cannot be nil",
		)
	}

	// Ensure that the problem is a linear program
	if !standardForm.IsLinear() {
		return utils.Tableau{}, fmt.Errorf(
			"ToTableau: the problem is not a linear program",
		)
	}

	return utils.Tableau{
		BasicVariables:    state.BasicVariables,
		NonBasicVariables: state.NonBasicVariables,
		NonBasicValues:    state.NonBasicValues,
		Problem:           standardForm,
	}, nil
}
