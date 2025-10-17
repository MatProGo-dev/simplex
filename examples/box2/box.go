package main

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	"github.com/MatProGo-dev/MatProInterface.go/solution"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/simplexSolver"
	"gonum.org/v1/gonum/mat"
)

/*
Description:

	This function builds an optimization problem where we attempt to find
	the optimal solution to a linear programming problem that is in a feasible
	region that is a box.

	The problem will be:

	Minimize:  x1 - 2*x2
	Subject to:
		-1 <= x1 <= 1
		-1 <= x2 <= 1

	The optimal solution is x1 = -1, x2 = 1 with an objective value of -3.
*/
func BuildOptimizationProblem() problem.OptimizationProblem {
	// setup
	varCount := 2
	out := problem.NewProblem("Box LP Problem")

	// Create the variables
	x := out.AddVariableVector(varCount)

	// Create the objective
	c := getKVector.From(
		[]float64{1, -2},
	)
	out.SetObjective(
		c.Transpose().Multiply(x),
		problem.SenseMinimize,
	)

	// Create the constraints
	// - x >= -1
	out.Constraints = append(
		out.Constraints,
		x.GreaterEq(mat.NewVecDense(varCount, []float64{-1, -1})),
	)

	// - x <= 1
	out.Constraints = append(
		out.Constraints,
		x.LessEq(symbolic.OnesVector(varCount)),
	)

	return *out
}

func main() {
	// This is just a placeholder to make the package "main" valid.
	trickyProblem := BuildOptimizationProblem()

	// Use solver to solve the problem
	solver := simplexSolver.New("Simplex Solver Example")
	solver.IterationLimit = 100

	// Solve the problem
	sol, err := solver.Solve(trickyProblem)
	if err != nil {
		panic(err)
	}

	// Print the solution
	solutionMessage, _ := sol.Status.ToMessage()
	println("Solution Status: ", solutionMessage)
	optObj, err := solution.GetOptimalObjectiveValue(&sol)
	if err != nil {
		panic(err)
	}
	println("Objective Value: ", optObj)
	println("Number of Iterations: ", sol.Iterations)
	println("Variable Values: ")
	for varName, varValue := range sol.VariableValues {
		println("  ", varName, ": ", varValue)
	}
}
