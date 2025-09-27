package examples

import (
	getKMatrix "github.com/MatProGo-dev/SymbolicMath.go/get/KMatrix"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/utils"
)

func GetTableauExample1() (*utils.Tableau, error) {
	// Setup

	// Create variables
	x := symbolic.NewVariableVector(2)
	s := symbolic.NewVariableVector(4)

	// Create the condensed tableau matrix directly
	condensedKMatrix := getKMatrix.From([][]float64{
		{-15, -25, 0.0, 0.0, 0.0, 0.0, 0.0},
		{1.0, 1.0, 1.0, 0.0, 0.0, 0.0, 450.0},
		{0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 300.0},
		{4.0, 5.0, 0.0, 0.0, 1.0, 0.0, 2000.0},
		{1.0, 0.0, 0.0, 0.0, 0.0, 1.0, 350.0},
	})
	condensed := condensedKMatrix.ToDense()

	// Create the tableau
	out := utils.Tableau{
		Variables:             append(x, s...),
		BasicVariableIndicies: []int{2, 3, 4, 5},
		AsCompressedMatrix:    &condensed,
	}

	return &out, out.Check()
}
