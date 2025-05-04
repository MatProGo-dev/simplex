package solver_test

import (
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"

	"matprogo.dev/solvers/simplex/simplexSolver"
)

/*
TestTransformAllUnboundedVariables1
Description:

	Tests the TransformAllUnboundedVariables function of the SimplexSolver.
	Here, we create a simple problem with one variable and one constraint.
	The variable is unbounded, and the constraint is a simple equality.
	After transforming the problem, we check that the new problem has two variables
	(one for the positive part and one for the negative part) and that the
	constraint is transformed correctly.
*/
func TestTransformAllUnboundedVariables1(t *testing.T) {
	// Setup

	// Create a new problem
	p1 := problem.NewProblem("TestTransformAllUnboundedVariables1 Problem")
	p1.AddVariable()
	p1.Variables[0].Name = "x1"
	v1 := p1.Variables[0]

	// Set the objective function
	p1.Constraints = append(
		p1.Constraints,
		v1.GreaterEq(-10.0),
	)

	// Set the objective function
	p1.SetObjective(
		v1,
		problem.SenseMinimize,
	)

	// Create a new solver
	solver := simplexSolver.New(p1.Name + " Solver")
	solver.OriginalProblem = p1

	// Transform all unbounded variables
	solver.TransformAllUnboundedVariables()

	// Check that the new problem has two variables
	if len(solver.ProblemWithAllPositiveVariables.Variables) != 2 {
		t.Errorf("Expected 2 variables, but got %d", len(solver.ProblemWithAllPositiveVariables.Variables))
	}

	// Check that there are the same number of constraints in the new problem as in the original problem
	if len(solver.ProblemWithAllPositiveVariables.Constraints) != len(p1.Constraints)+1 {
		t.Errorf("Expected 3 constraints, but got %d", len(solver.ProblemWithAllPositiveVariables.Constraints))
	}

	// Check that the first constraint contains the difference of the two new variables
	if len(solver.ProblemWithAllPositiveVariables.Constraints[2].Left().Variables()) != 2 {
		t.Errorf("Expected 2 variables in the first constraint, but got %d", len(solver.ProblemWithAllPositiveVariables.Constraints[0].Left().Variables()))
		t.Errorf("Constraint: %s", solver.ProblemWithAllPositiveVariables.Constraints[2])
	}

	// Check that the new objective function contains 2 variables instead of 1
	if len(solver.ProblemWithAllPositiveVariables.Objective.Expression.Variables()) != 2 {
		t.Errorf("Expected 2 variables in the objective function, but got %d", len(solver.ProblemWithAllPositiveVariables.Objective.Expression.Variables()))
		t.Errorf("Objective: %s", solver.ProblemWithAllPositiveVariables.Objective.Expression)
	}

}
