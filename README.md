# simplex
A small library used to demonstrate the concepts of the Simplex algorithm in convex optimization.

# Installation

You can add this module to your package using:
```bash
go get github.com/MatProGo-dev/simplex
```

# Usage

```
package main

import (
	"github.com/MatProGo-dev/MatProInterface.go/problem"
	getKVector "github.com/MatProGo-dev/SymbolicMath.go/get/KVector"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/simplexSolver"
)

func BuildOptimizationProblem() problem.OptimizationProblem {
	// setup
	varCount := 2
	out := problem.NewProblem("Box LP Problem")

	// Create the variables
	x := out.AddVariableVector(varCount)

	// Create the objective
	c := getKVector.From(
		[]float64{1, 2},
	)
	out.SetObjective(
		c.Transpose().Multiply(x),
		problem.SenseMinimize,
	)

	// Create the constraints
	// - x >= 0
	out.Constraints = append(
		out.Constraints,
		x.GreaterEq(symbolic.ZerosVector(varCount)),
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
	solution, err := solver.Solve(trickyProblem)
	if err != nil {
		panic(err)
	}

	// Print the solution
	solutionMessage, _ := solution.Status.ToMessage()
	println("Solution Status: ", solutionMessage)
	println("Objective Value: ", solution.Objective)
	println("Number of Iterations: ", solution.Iterations)
	println("Variable Values: ")
	for varName, varValue := range solution.VariableValues {
		println("  ", varName, ": ", varValue)
	}
}
```

See the examples directory for more example use cases for the library.