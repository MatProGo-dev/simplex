package main

import (
	"github.com/MatProGo-dev/MatProInterface.go/solution"
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
	sol, err := solver.Solve(*problem5)
	if err != nil {
		panic(err)
	}

	// Print the solution
	solutionMessage, _ := sol.Status.ToMessage()
	println("Solution Status: ", solutionMessage)
	optVal, _ := solution.GetOptimalObjectiveValue(&sol)
	println("Objective Value: ", optVal)
	println("Number of Iterations: ", sol.Iterations)
	println("Variable Values: ")
	for varName, varValue := range sol.VariableValues {
		println("  ", varName, ": ", varValue)
	}
}
