package linprog_test

import (
	"math"
	"testing"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	"github.com/MatProGo-dev/simplex/algorithms"
	linprog_algorithm "github.com/MatProGo-dev/simplex/algorithms/linprog"
	"github.com/MatProGo-dev/simplex/simplexSolver"
	"github.com/MatProGo-dev/simplex/utils/examples"
)

func TestLinprogSolvesBoundedMaximization(t *testing.T) {
	prob := problem.NewProblem("bounded")
	x := prob.AddVariableVector(2)
	if err := prob.SetObjective(x.AtVec(0).Plus(x.AtVec(1)), problem.SenseMaximize); err != nil {
		t.Fatal(err)
	}
	prob.Constraints = append(prob.Constraints, x.AtVec(0).Plus(x.AtVec(1)).LessEq(4.0))

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.GetOptimalValue()-4) > 1e-7 {
		t.Fatalf("expected objective 4, got %v", solution.GetOptimalValue())
	}
	if math.Abs(solution.VariableValues[prob.Variables[0].ID]+solution.VariableValues[prob.Variables[1].ID]-4) > 1e-7 {
		t.Fatalf("returned variables do not satisfy the active constraint: %v", solution.VariableValues)
	}
}

func TestLinprogSolvesRepositoryExample(t *testing.T) {
	prob := examples.GetTestProblem5()
	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.VariableValues[prob.Variables[0].ID]-125) > 1e-7 {
		t.Fatalf("expected x1=125, got %v", solution.VariableValues[prob.Variables[0].ID])
	}
	if math.Abs(solution.VariableValues[prob.Variables[1].ID]-300) > 1e-7 {
		t.Fatalf("expected x2=300, got %v", solution.VariableValues[prob.Variables[1].ID])
	}
}

func TestLinprogReportsInfeasibleBoundsAndConstraints(t *testing.T) {
	prob := problem.NewProblem("infeasible")
	x := prob.AddVariableClassic(0, 1, symbolic.Continuous)
	if err := prob.SetObjective(x, problem.SenseMinimize); err != nil {
		t.Fatal(err)
	}
	prob.Constraints = append(prob.Constraints, x.GreaterEq(2.0))

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.INFEASIBLE {
		t.Fatalf("expected infeasible status, got %v", solution.Status)
	}
}

func TestLinprogUsesDualSimplexToRestorePrimalFeasibility(t *testing.T) {
	prob := problem.NewProblem("dual simplex pivot")
	x := prob.AddVariableVector(1)
	if err := prob.SetObjective(x.AtVec(0), problem.SenseMinimize); err != nil {
		t.Fatal(err)
	}
	prob.Constraints = append(prob.Constraints, x.AtVec(0).GreaterEq(2.0))

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.VariableValues[prob.Variables[0].ID]-2) > 1e-7 {
		t.Fatalf("expected x=2, got %v", solution.VariableValues[prob.Variables[0].ID])
	}
}

func TestLinprogHandlesEqualitiesAndRedundantRows(t *testing.T) {
	prob := problem.NewProblem("equality constraints")
	x := prob.AddVariableVector(2)
	if err := prob.SetObjective(x.AtVec(0).Plus(x.AtVec(1).Multiply(2)), problem.SenseMaximize); err != nil {
		t.Fatal(err)
	}
	equality := x.AtVec(0).Plus(x.AtVec(1)).Comparison(3.0, symbolic.SenseEqual)
	prob.Constraints = append(
		prob.Constraints,
		equality,
		equality,
		x.AtVec(0).LessEq(2.0),
		x.AtVec(0).GreaterEq(0.0),
		x.AtVec(1).GreaterEq(0.0),
	)

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.GetOptimalValue()-6) > 1e-7 {
		t.Fatalf("expected objective 6, got %v", solution.GetOptimalValue())
	}
}

func TestLinprogReportsUnboundedObjective(t *testing.T) {
	prob := problem.NewProblem("unbounded")
	x := prob.AddVariableVector(1)
	if err := prob.SetObjective(x.AtVec(0), problem.SenseMaximize); err != nil {
		t.Fatal(err)
	}

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.UNBOUNDED {
		t.Fatalf("expected unbounded status, got %v", solution.Status)
	}
}

func TestLinprogHonorsVariableUpperBounds(t *testing.T) {
	prob := problem.NewProblem("finite bound")
	x := prob.AddVariableClassic(0, 3, symbolic.Continuous)
	if err := prob.SetObjective(x.Multiply(-1), problem.SenseMinimize); err != nil {
		t.Fatal(err)
	}

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.VariableValues[x.ID]-3) > 1e-7 {
		t.Fatalf("expected x=3 at its upper bound, got %v", solution.VariableValues[x.ID])
	}
}

func TestLinprogHandlesFreeVariables(t *testing.T) {
	prob := problem.NewProblem("free variable")
	x := prob.AddRealVariable()
	if err := prob.SetObjective(x, problem.SenseMaximize); err != nil {
		t.Fatal(err)
	}
	prob.Constraints = append(prob.Constraints, x.LessEq(3.0), x.GreaterEq(-2.0))

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 100}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
	if math.Abs(solution.VariableValues[x.ID]-3) > 1e-7 {
		t.Fatalf("expected x=3, got %v", solution.VariableValues[x.ID])
	}
}

func TestLinprogReportsIterationLimit(t *testing.T) {
	prob := problem.NewProblem("iteration limit")
	x := prob.AddVariableVector(2)
	if err := prob.SetObjective(x.AtVec(0).Plus(x.AtVec(1)), problem.SenseMaximize); err != nil {
		t.Fatal(err)
	}

	prob.Constraints = append(prob.Constraints, x.AtVec(0).Plus(x.AtVec(1)).LessEq(4.0))

	solution, err := (&linprog_algorithm.LinprogAlgorithm{IterationLimit: 0}).Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.ITERATION_LIMIT {
		t.Fatalf("expected iteration-limit status, got %v", solution.Status)
	}
}

func TestLinprogCanBeSelectedThroughSimplexSolver(t *testing.T) {
	prob := problem.NewProblem("solver selection")
	x := prob.AddVariableVector(1)
	if err := prob.SetObjective(x.AtVec(0), problem.SenseMinimize); err != nil {
		t.Fatal(err)
	}
	prob.Constraints = append(prob.Constraints, x.AtVec(0).GreaterEq(2.0))

	solver := simplexSolver.New("linprog")
	solver.Algorithm = algorithms.TypeLinprog
	solution, err := solver.Solve(*prob)
	if err != nil {
		t.Fatal(err)
	}
	if solution.Status != solution_status.OPTIMAL {
		t.Fatalf("expected optimal status, got %v", solution.Status)
	}
}
