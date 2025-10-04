package main

import (
	"github.com/MatProGo-dev/simplex/simplexSolver"
	"github.com/MatProGo-dev/simplex/utils/examples"
)

func main() {
	// Get the example tableau
	problem5 := examples.GetTestProblem5()

	// // Print the problem
	// println(problem5.String())

	// Use solver to solve the problem
	solver := simplexSolver.New("Simplex Solver Example")
	solver.IterationLimit = 100

	// Solve the problem
	solution, err := solver.Solve(*problem5)
	if err != nil {
		panic(err)
	}

	// Print the solution
	solutionMessage, _ := solution.Status.ToMessage()
	println("Solution Status: ", solutionMessage)
	println("Objective Value: ", solution.Objective)
	println("Variable Values: ")
	for varName, varValue := range solution.VariableValues {
		println("  ", varName, ": ", varValue)
	}
}
