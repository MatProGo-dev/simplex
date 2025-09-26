package dictionary

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/simplex/utils"
)

/*
ComputeInitialState
Description:

Constructs an initial DictionaryAlgorithmState value for the given problem based on the problem
definition.
*/
func ComputeInitialState(problemIn *problem.OptimizationProblem) DictionaryAlgorithmState {
	// Setup

	// Filter constraints based on whether or not they are simple inequality constraints
	asScalarConstraints := utils.ExtractScalarConstraints(
		problemIn.Constraints,
	)

	// Select the basic variables
	NBasicVariables := len(asScalarConstraints)

	basicVariableIndicies := []int{}
	for ii := 0; ii < NBasicVariables; ii++ {
		basicVariableIndicies = append(
			basicVariableIndicies,
			len(problemIn.Variables)-NBasicVariables+ii)
	}

	return DictionaryAlgorithmState{
		AllVariables:          problemIn.Variables,
		BasicVariableIndicies: basicVariableIndicies,
		IterationCount:        0,
		ObjectiveExpression:   problemIn.Objective.Expression,
		DictionaryConstraints: asScalarConstraints,
	}

}
