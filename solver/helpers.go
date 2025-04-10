package solver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
)

func GetTestProblem1() *problem.OptimizationProblem {
	// Setup
	out := problem.NewProblem("TestProblem1")

	// Create variables
	x := out.AddVariableVector(2)

	// Create Basic Objective
	c := getKVector.From([]float64{1.0, 2.0})
	out.SetObjective(
		c.Transpose().Multiply(x),
		problem.SenseMinimize,
	)

	// Create Constraints (using one big matrix)
	A := getKMatrix.From([][]float64{
		{-1.0, 0.0},
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, -1.0},
	})

	b := getKVector.From([]float64{0.0, 1.0, 1.5, 1.0, 0.0})
	out.Constraints = append(out.Constraints, A.Multiply(x).LessEq(b))

	return out
}

func GetTestProblem2() *problem.OptimizationProblem {
	// Setup
	out := problem.NewProblem("TestProblem2")

	// Create variables
	x := out.AddVariableVector(2)

	// Create Basic Objective
	c := getKVector.From([]float64{1.0, 2.0})
	out.SetObjective(
		c.Transpose().Multiply(x),
		problem.SenseMinimize,
	)

	// Create Constraints (using individual constraints)
	// 1. x1 >= 0
	out.Constraints = append(
		out.Constraints,
		x.AtVec(0).GreaterEq(0.0),
	)

	// 2. x2 <= 1
	out.Constraints = append(
		out.Constraints,
		x.AtVec(1).LessEq(1.0),
	)

	// 3. x1 + x2 <= 1.5
	out.Constraints = append(
		out.Constraints,
		x.AtVec(0).Plus(x.AtVec(1)).LessEq(1.5),
	)

	// 4. x1 >= 1
	out.Constraints = append(
		out.Constraints,
		x.AtVec(0).LessEq(1.0),
	)

	// 5. x2 >= 0
	out.Constraints = append(
		out.Constraints,
		x.AtVec(1).GreaterEq(0.0),
	)

	return out
}
