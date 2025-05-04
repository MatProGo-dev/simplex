package simplexSolver

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
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

/*
GetTestProblem3
Description:

	Returns the LP from this youtube video:
		https://www.youtube.com/watch?v=QAR8zthQypc&t=483s
	It should look like this:
		Maximize	4 x1 + 3 x2 + 5 x3
		Subject to
			x1 + 2 x2 + 2 x3 <= 4
			3 x1 + 4 x3 <= 6
			2 x1 + x2 + 4 x3 <= 8
			x1 >= 0
			x2 >= 0
			x3 >= 0
*/
func GetTestProblem3() *problem.OptimizationProblem {
	// Setup
	out := problem.NewProblem("TestProblem3")

	// Create variables
	x := out.AddVariableVectorClassic(
		3,
		0.0,
		symbolic.Infinity.Constant(),
		symbolic.Continuous,
	)

	// Create Basic Objective
	c := getKVector.From([]float64{4.0, 3.0, 5.0})
	out.SetObjective(
		c.Transpose().Multiply(x),
		problem.SenseMaximize,
	)

	// Create Constraints (using one big matrix)
	A := getKMatrix.From([][]float64{
		{1.0, 2.0, 2.0},
		{3.0, 0.0, 4.0},
		{2.0, 1.0, 4.0},
	})
	b := getKVector.From([]float64{4.0, 6.0, 8.0})
	out.Constraints = append(out.Constraints, A.Multiply(x).LessEq(b))

	// TODO(kwesi): Figure out how to add non-negativity constraints
	// // Add non-negativity constraints
	// for _, varII := range x {
	// 	out.Constraints = append(
	// 		out.Constraints,
	// 		varII.GreaterEq(0.0),
	// 	)
	// }
	return out
}
