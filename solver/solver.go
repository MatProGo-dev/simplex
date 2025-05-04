package solver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
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

func (solver *SimplexSolver) TransformAllUnboundedVariables() {
	// Setup
	originalProblem := solver.OriginalProblem

	// For each variable, let's create two new variables
	// and set the original variable to be the difference of the two
	mapFromOriginalVariablesToNewExpressions := make(map[symbolic.Variable]symbolic.Expression)
	for ii := 0; ii < len(originalProblem.Variables); ii++ {
		// Setup
		xII := originalProblem.Variables[ii]

		// Create the two new variables
		solver.ProblemWithAllPositiveVariables.AddVariableClassic(0.0, symbolic.Infinity.Constant(), symbolic.Continuous)
		nVariables := len(solver.ProblemWithAllPositiveVariables.Variables)
		solver.ProblemWithAllPositiveVariables.Variables[nVariables-1].Name = xII.Name + " (+)"
		variablePositivePart := solver.ProblemWithAllPositiveVariables.Variables[nVariables-1]

		solver.ProblemWithAllPositiveVariables.AddVariableClassic(0.0, symbolic.Infinity.Constant(), symbolic.Continuous)
		nVariables = len(solver.ProblemWithAllPositiveVariables.Variables)
		solver.ProblemWithAllPositiveVariables.Variables[nVariables-1].Name = xII.Name + " (-)"
		variableNegativePart := solver.ProblemWithAllPositiveVariables.Variables[nVariables-1]

		// Set the original variable to be the difference of the two new variables
		mapFromOriginalVariablesToNewExpressions[xII] =
			variablePositivePart.Minus(variableNegativePart)
	}

	// Now, let's create the new constraints by replacing the variables in the
	// original constraints with the new expressions
	for _, constraint := range originalProblem.Constraints {
		// Create a new constraint
		constraintCopy := constraint

		// Create a new expression by substituting the variables according
		// to the map we created above
		oldLHS := constraint.Left()
		newLHS := oldLHS.SubstituteAccordingTo(mapFromOriginalVariablesToNewExpressions)

		oldRHS := constraint.Right()
		newRHS := oldRHS.SubstituteAccordingTo(mapFromOriginalVariablesToNewExpressions)

		newConstraint := newLHS.Comparison(
			newRHS,
			constraintCopy.ConstrSense(),
		)

		// Add the new constraint to the problem
		solver.ProblemWithAllPositiveVariables.Constraints = append(
			solver.ProblemWithAllPositiveVariables.Constraints,
			newConstraint,
		)
	}

	// Now, let's create the new objective function by substituting the variables
	// according to the map we created above
	newObjectiveExpression := solver.OriginalProblem.Objective.Expression.SubstituteAccordingTo(
		mapFromOriginalVariablesToNewExpressions,
	)
	solver.ProblemWithAllPositiveVariables.SetObjective(
		newObjectiveExpression,
		solver.OriginalProblem.Objective.Sense,
	)
}
