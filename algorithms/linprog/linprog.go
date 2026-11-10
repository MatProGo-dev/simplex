// Package linprog implements a continuous linear-programming solver based on
// dual revised simplex. It does not include SciPy's presolve, interior-point,
// mixed-integer, or Python result and option features.
package linprog

import (
	"fmt"
	"math"

	"github.com/MatProGo-dev/MatProInterface.go/problem"
	solution_status "github.com/MatProGo-dev/MatProInterface.go/solution/status"
	"github.com/MatProGo-dev/SymbolicMath.go/symbolic"
	simplex_solution "github.com/MatProGo-dev/simplex/solution"
	"gonum.org/v1/gonum/mat"
)

const feasibilityTolerance = 1e-9

// LinprogAlgorithm solves continuous linear programs with a dual revised
// simplex method.
type LinprogAlgorithm struct {
	IterationLimit int
}

type feasibilityResult struct {
	values     []float64
	feasible   bool
	iterations int
	limitHit   bool
}

// Solve converts the model to standard form, finds a dual-feasible starting
// basis, and applies dual simplex pivots until the primal solution is feasible.
func (algo *LinprogAlgorithm) Solve(prob problem.OptimizationProblem) (simplex_solution.SimplexSolution, error) {
	workingProblem := prob.Copy()
	for _, variable := range workingProblem.Variables {
		if variable.Type != symbolic.Continuous {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: variable %q is not continuous", variable.Name)
		}
		if math.IsNaN(variable.Lower) || math.IsNaN(variable.Upper) || variable.Lower > variable.Upper {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: variable %q has invalid bounds", variable.Name)
		}
		if !isInfiniteBound(variable.Lower) {
			workingProblem.Constraints = append(workingProblem.Constraints, variable.GreaterEq(variable.Lower))
		}
		if !isInfiniteBound(variable.Upper) {
			workingProblem.Constraints = append(workingProblem.Constraints, variable.LessEq(variable.Upper))
		}
	}

	standardProblem, _, originalVariableExpressions, err := workingProblem.ToLPStandardForm1()
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: could not convert problem to standard form: %w", err)
	}
	if !standardProblem.IsLinear() {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: standard-form problem is not linear")
	}

	objective, ok := standardProblem.Objective.Expression.(symbolic.ScalarExpression)
	if !ok {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: objective is not a scalar linear expression")
	}
	c := make([]float64, len(standardProblem.Variables))
	if len(c) > 0 {
		cVector := objective.LinearCoeff(standardProblem.Variables)
		copy(c, cVector.RawVector().Data)
	}
	switch standardProblem.Objective.Sense {
	case problem.SenseMinimize:
	case problem.SenseMaximize:
		for i := range c {
			c[i] = -c[i]
		}
	default:
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: unsupported objective sense %q", standardProblem.Objective.Sense)
	}
	for i, value := range c {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: objective coefficient %d is not finite", i)
		}
	}

	if len(standardProblem.Constraints) == 0 {
		for _, coefficient := range c {
			if coefficient < -feasibilityTolerance {
				return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.UNBOUNDED, 0)
			}
		}
		return makeSolution(&prob, make([]float64, len(standardProblem.Variables)), standardProblem.Variables, originalVariableExpressions, solution_status.OPTIMAL, 0)
	}

	aSymbolic, bSymbolic, err := standardProblem.LinearEqualityConstraintMatrices()
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: could not extract standard-form constraints: %w", err)
	}
	aDense := aSymbolic.ToDense()
	bVector := bSymbolic.ToVecDense()
	b := append([]float64(nil), bVector.RawVector().Data...)
	for i, value := range b {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: right-hand-side value %d is not finite", i)
		}
	}
	rows, columns := aDense.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			value := aDense.At(i, j)
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: constraint coefficient (%d, %d) is not finite", i, j)
			}
		}
	}

	a, b, consistent, err := reduceEqualityRows(&aDense, b)
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: could not reduce equality constraints: %w", err)
	}
	if !consistent {
		return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.INFEASIBLE, 0)
	}
	if len(b) == 0 {
		for _, coefficient := range c {
			if coefficient < -feasibilityTolerance {
				return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.UNBOUNDED, 0)
			}
		}
		return makeSolution(&prob, make([]float64, len(standardProblem.Variables)), standardProblem.Variables, originalVariableExpressions, solution_status.OPTIMAL, 0)
	}

	limit := algo.IterationLimit
	iterations := 0
	dualMatrix, dualRHS := makeDualFeasibilitySystem(a, c)
	dualFeasibility, err := findFeasiblePoint(dualMatrix, dualRHS, limit)
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: finding a dual-feasible point: %w", err)
	}
	iterations += dualFeasibility.iterations
	if dualFeasibility.limitHit {
		return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.ITERATION_LIMIT, iterations)
	}

	if !dualFeasibility.feasible {
		primalFeasibility, err := findFeasiblePoint(a, b, remainingIterations(limit, iterations))
		if err != nil {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: checking primal feasibility: %w", err)
		}
		iterations += primalFeasibility.iterations
		if primalFeasibility.limitHit {
			return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.ITERATION_LIMIT, iterations)
		}
		status := solution_status.UNBOUNDED
		if !primalFeasibility.feasible {
			status = solution_status.INFEASIBLE
		}
		return makeSolution(&prob, nil, nil, originalVariableExpressions, status, iterations)
	}

	rowCount, variableCount := a.Dims()
	dualPrices := make([]float64, rowCount)
	for i := range dualPrices {
		dualPrices[i] = dualFeasibility.values[i] - dualFeasibility.values[rowCount+i]
	}
	basis, err := findDualFeasibleBasis(a, c, dualPrices)
	if err != nil {
		return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: could not construct a dual-feasible basis: %w", err)
	}

	for {
		basisMatrix, err := basisMatrix(a, basis)
		if err != nil {
			return simplex_solution.SimplexSolution{}, err
		}
		basicValues, err := solve(basisMatrix, b)
		if err != nil {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: solving for basic values: %w", err)
		}
		basicCosts := make([]float64, rowCount)
		isBasic := make([]bool, variableCount)
		for i, variableIndex := range basis {
			basicCosts[i] = c[variableIndex]
			isBasic[variableIndex] = true
		}
		dualPrices, err = solveTranspose(basisMatrix, basicCosts)
		if err != nil {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: solving for dual prices: %w", err)
		}
		reducedCosts := make([]float64, variableCount)
		for j := 0; j < variableCount; j++ {
			reducedCosts[j] = c[j]
			for i := 0; i < rowCount; i++ {
				reducedCosts[j] -= a.At(i, j) * dualPrices[i]
			}
			if !isBasic[j] && reducedCosts[j] < -feasibilityTolerance {
				return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: dual-feasibility invariant lost at variable %d", j)
			}
		}

		leavingRow := -1
		for i, value := range basicValues {
			if value < -feasibilityTolerance && (leavingRow == -1 || basis[i] < basis[leavingRow]) {
				leavingRow = i
			}
		}
		if leavingRow == -1 {
			values := make([]float64, variableCount)
			for i, variableIndex := range basis {
				if basicValues[i] < 0 && basicValues[i] >= -feasibilityTolerance {
					basicValues[i] = 0
				}
				values[variableIndex] = basicValues[i]
			}
			return makeSolution(&prob, values, standardProblem.Variables, originalVariableExpressions, solution_status.OPTIMAL, iterations)
		}
		if iterations >= limit {
			return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.ITERATION_LIMIT, iterations)
		}

		rowOfInverse, err := solveTranspose(basisMatrix, unitVector(rowCount, leavingRow))
		if err != nil {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: solving dual-simplex pivot row: %w", err)
		}
		enteringVariable := -1
		bestRatio := math.Inf(1)
		for j := 0; j < variableCount; j++ {
			if isBasic[j] {
				continue
			}
			pivotCoefficient := 0.0
			for i := 0; i < rowCount; i++ {
				pivotCoefficient += rowOfInverse[i] * a.At(i, j)
			}
			if pivotCoefficient >= -feasibilityTolerance {
				continue
			}
			ratio := math.Max(0, reducedCosts[j]) / -pivotCoefficient
			if ratio < bestRatio-feasibilityTolerance ||
				(math.Abs(ratio-bestRatio) <= feasibilityTolerance && (enteringVariable == -1 || j < enteringVariable)) {
				bestRatio = ratio
				enteringVariable = j
			}
		}
		if enteringVariable == -1 {
			return makeSolution(&prob, nil, nil, originalVariableExpressions, solution_status.INFEASIBLE, iterations)
		}
		basis[leavingRow] = enteringVariable
		iterations++
	}
}

func isInfiniteBound(value float64) bool {
	return math.IsInf(value, 0) || math.Abs(value) >= float64(symbolic.Infinity)
}

func remainingIterations(limit, used int) int {
	if limit <= used {
		return 0
	}
	return limit - used
}

func makeSolution(
	originalProblem *problem.OptimizationProblem,
	values []float64,
	standardVariables []symbolic.Variable,
	originalVariableExpressions map[symbolic.Variable]symbolic.Expression,
	status solution_status.SolutionStatus,
	iterations int,
) (simplex_solution.SimplexSolution, error) {
	sol := simplex_solution.SimplexSolution{
		VariableValues:  make(map[uint64]float64),
		Status:          status,
		Iterations:      iterations,
		OriginalProblem: originalProblem,
	}
	if status != solution_status.OPTIMAL {
		return sol, nil
	}
	standardValues := make(map[symbolic.Variable]symbolic.Expression, len(standardVariables))
	for i, variable := range standardVariables {
		standardValues[variable] = symbolic.K(values[i])
	}
	for originalVariable, expression := range originalVariableExpressions {
		value, ok := expression.SubstituteAccordingTo(standardValues).(symbolic.K)
		if !ok {
			return simplex_solution.SimplexSolution{}, fmt.Errorf("linprog: could not evaluate original variable %q", originalVariable.Name)
		}
		sol.VariableValues[originalVariable.ID] = float64(value)
	}
	return sol, nil
}

func makeDualFeasibilitySystem(a mat.Matrix, c []float64) (*mat.Dense, []float64) {
	rowCount, variableCount := a.Dims()
	matrix := mat.NewDense(variableCount, 2*rowCount+variableCount, nil)
	for j := 0; j < variableCount; j++ {
		for i := 0; i < rowCount; i++ {
			matrix.Set(j, i, a.At(i, j))
			matrix.Set(j, rowCount+i, -a.At(i, j))
		}
		matrix.Set(j, 2*rowCount+j, 1)
	}
	return matrix, append([]float64(nil), c...)
}

func findDualFeasibleBasis(a mat.Matrix, c, dualPrices []float64) ([]int, error) {
	rowCount, variableCount := a.Dims()
	basis := make([]int, 0, rowCount)
	var current *mat.Dense
	currentRank := 0
	for j := 0; j < variableCount && len(basis) < rowCount; j++ {
		reducedCost := c[j]
		for i := 0; i < rowCount; i++ {
			reducedCost -= a.At(i, j) * dualPrices[i]
		}
		if math.Abs(reducedCost) > 10*feasibilityTolerance {
			continue
		}
		candidate := mat.NewDense(rowCount, len(basis)+1, nil)
		for col := 0; col < len(basis); col++ {
			for row := 0; row < rowCount; row++ {
				candidate.Set(row, col, current.At(row, col))
			}
		}
		for row := 0; row < rowCount; row++ {
			candidate.Set(row, len(basis), a.At(row, j))
		}
		rank := matrixRank(candidate)
		if rank > currentRank {
			basis = append(basis, j)
			current = candidate
			currentRank = rank
		}
	}
	if len(basis) != rowCount {
		return nil, fmt.Errorf("only found %d independent tight columns for %d equality rows", len(basis), rowCount)
	}
	return basis, nil
}

func findFeasiblePoint(a mat.Matrix, rhs []float64, iterationLimit int) (feasibilityResult, error) {
	rowCount, variableCount := a.Dims()
	if len(rhs) != rowCount {
		return feasibilityResult{}, fmt.Errorf("constraint matrix has %d rows but right-hand side has length %d", rowCount, len(rhs))
	}
	if rowCount == 0 {
		return feasibilityResult{values: make([]float64, variableCount), feasible: true}, nil
	}

	augmented := mat.NewDense(rowCount, variableCount+rowCount, nil)
	normalizedRHS := append([]float64(nil), rhs...)
	for i := 0; i < rowCount; i++ {
		sign := 1.0
		if normalizedRHS[i] < 0 {
			sign = -1
			normalizedRHS[i] = -normalizedRHS[i]
		}
		for j := 0; j < variableCount; j++ {
			augmented.Set(i, j, sign*a.At(i, j))
		}
		augmented.Set(i, variableCount+i, 1)
	}

	totalVariables := variableCount + rowCount
	costs := make([]float64, totalVariables)
	for j := variableCount; j < totalVariables; j++ {
		costs[j] = -1
	}
	basis := make([]int, rowCount)
	for i := range basis {
		basis[i] = variableCount + i
	}

	pivots := 0
	for {
		basisMatrix, err := basisMatrix(augmented, basis)
		if err != nil {
			return feasibilityResult{}, err
		}
		basicValues, err := solve(basisMatrix, normalizedRHS)
		if err != nil {
			return feasibilityResult{}, fmt.Errorf("solving Phase-I basic values: %w", err)
		}
		basicCosts := make([]float64, rowCount)
		isBasic := make([]bool, totalVariables)
		for i, index := range basis {
			basicCosts[i] = costs[index]
			isBasic[index] = true
		}
		dual, err := solveTranspose(basisMatrix, basicCosts)
		if err != nil {
			return feasibilityResult{}, fmt.Errorf("solving Phase-I dual prices: %w", err)
		}

		entering := -1
		for j := 0; j < totalVariables; j++ {
			if isBasic[j] {
				continue
			}
			reducedCost := costs[j]
			for i := 0; i < rowCount; i++ {
				reducedCost -= augmented.At(i, j) * dual[i]
			}
			if reducedCost > feasibilityTolerance {
				entering = j
				break
			}
		}
		if entering == -1 {
			artificialSum := 0.0
			for i, variableIndex := range basis {
				if variableIndex >= variableCount {
					artificialSum += math.Max(0, basicValues[i])
				}
			}
			values := make([]float64, variableCount)
			for i, variableIndex := range basis {
				if variableIndex < variableCount {
					values[variableIndex] = math.Max(0, basicValues[i])
				}
			}
			return feasibilityResult{
				values:     values,
				feasible:   artificialSum <= feasibilityTolerance,
				iterations: pivots,
			}, nil
		}
		if iterationLimit <= 0 {
			return feasibilityResult{iterations: pivots, limitHit: true}, nil
		}

		column := make([]float64, rowCount)
		for i := range column {
			column[i] = augmented.At(i, entering)
		}
		direction, err := solve(basisMatrix, column)
		if err != nil {
			return feasibilityResult{}, fmt.Errorf("solving Phase-I pivot direction: %w", err)
		}
		leaving := -1
		bestRatio := math.Inf(1)
		for i, coefficient := range direction {
			if coefficient <= feasibilityTolerance {
				continue
			}
			ratio := basicValues[i] / coefficient
			if ratio < bestRatio-feasibilityTolerance ||
				(math.Abs(ratio-bestRatio) <= feasibilityTolerance && (leaving == -1 || basis[i] < basis[leaving])) {
				bestRatio = ratio
				leaving = i
			}
		}
		if leaving == -1 {
			return feasibilityResult{}, fmt.Errorf("Phase-I feasibility objective is unexpectedly unbounded")
		}
		basis[leaving] = entering
		iterationLimit--
		pivots++
	}
}

func reduceEqualityRows(a *mat.Dense, rhs []float64) (*mat.Dense, []float64, bool, error) {
	rowCount, variableCount := a.Dims()
	if len(rhs) != rowCount {
		return nil, nil, false, fmt.Errorf("constraint matrix has %d rows but right-hand side has length %d", rowCount, len(rhs))
	}
	if rowCount == 0 {
		return &mat.Dense{}, nil, true, nil
	}
	if variableCount == 0 {
		for _, value := range rhs {
			if math.Abs(value) > feasibilityTolerance {
				return nil, nil, false, nil
			}
		}
		return &mat.Dense{}, nil, true, nil
	}
	augmented := mat.NewDense(rowCount, variableCount+1, nil)
	for i := 0; i < rowCount; i++ {
		for j := 0; j < variableCount; j++ {
			augmented.Set(i, j, a.At(i, j))
		}
		augmented.Set(i, variableCount, rhs[i])
	}
	if matrixRank(augmented) > matrixRank(a) {
		return nil, nil, false, nil
	}

	selected := make([]int, 0, rowCount)
	rank := 0
	for i := 0; i < rowCount; i++ {
		candidate := mat.NewDense(len(selected)+1, variableCount, nil)
		for row, selectedIndex := range selected {
			for col := 0; col < variableCount; col++ {
				candidate.Set(row, col, a.At(selectedIndex, col))
			}
		}
		for col := 0; col < variableCount; col++ {
			candidate.Set(len(selected), col, a.At(i, col))
		}
		candidateRank := matrixRank(candidate)
		if candidateRank > rank {
			selected = append(selected, i)
			rank = candidateRank
		}
	}
	if len(selected) == 0 {
		return &mat.Dense{}, nil, true, nil
	}

	reduced := mat.NewDense(len(selected), variableCount, nil)
	reducedRHS := make([]float64, len(selected))
	for row, sourceRow := range selected {
		for col := 0; col < variableCount; col++ {
			reduced.Set(row, col, a.At(sourceRow, col))
		}
		reducedRHS[row] = rhs[sourceRow]
	}
	return reduced, reducedRHS, true, nil
}

func matrixRank(matrix mat.Matrix) int {
	rows, cols := matrix.Dims()
	if rows == 0 || cols == 0 {
		return 0
	}
	data := make([][]float64, rows)
	scale := 0.0
	for i := 0; i < rows; i++ {
		data[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			data[i][j] = matrix.At(i, j)
			scale = math.Max(scale, math.Abs(data[i][j]))
		}
	}
	if scale == 0 {
		return 0
	}
	tolerance := feasibilityTolerance * scale
	rank := 0
	for col := 0; col < cols && rank < rows; col++ {
		pivot := rank
		for row := rank + 1; row < rows; row++ {
			if math.Abs(data[row][col]) > math.Abs(data[pivot][col]) {
				pivot = row
			}
		}
		if math.Abs(data[pivot][col]) <= tolerance {
			continue
		}
		data[rank], data[pivot] = data[pivot], data[rank]
		pivotValue := data[rank][col]
		for row := rank + 1; row < rows; row++ {
			factor := data[row][col] / pivotValue
			for j := col; j < cols; j++ {
				data[row][j] -= factor * data[rank][j]
			}
		}
		rank++
	}
	return rank
}

func basisMatrix(a mat.Matrix, basis []int) (*mat.Dense, error) {
	rows, _ := a.Dims()
	if len(basis) != rows {
		return nil, fmt.Errorf("basis has %d columns for %d rows", len(basis), rows)
	}
	if rows == 0 {
		return &mat.Dense{}, nil
	}
	matrix := mat.NewDense(rows, rows, nil)
	for col, variableIndex := range basis {
		for row := 0; row < rows; row++ {
			matrix.Set(row, col, a.At(row, variableIndex))
		}
	}
	return matrix, nil
}

func solve(matrix *mat.Dense, rhs []float64) ([]float64, error) {
	rows, cols := matrix.Dims()
	if rows != cols || len(rhs) != rows {
		return nil, fmt.Errorf("cannot solve %dx%d system with right-hand side length %d", rows, cols, len(rhs))
	}
	if rows == 0 {
		return []float64{}, nil
	}
	right := mat.NewDense(rows, 1, append([]float64(nil), rhs...))
	var result mat.Dense
	if err := result.Solve(matrix, right); err != nil {
		return nil, err
	}
	return append([]float64(nil), result.RawMatrix().Data...), nil
}

func solveTranspose(matrix *mat.Dense, rhs []float64) ([]float64, error) {
	rows, cols := matrix.Dims()
	if rows != cols || len(rhs) != rows {
		return nil, fmt.Errorf("cannot solve transposed %dx%d system with right-hand side length %d", rows, cols, len(rhs))
	}
	if rows == 0 {
		return []float64{}, nil
	}
	right := mat.NewDense(rows, 1, append([]float64(nil), rhs...))
	var result mat.Dense
	if err := result.Solve(matrix.T(), right); err != nil {
		return nil, err
	}
	return append([]float64(nil), result.RawMatrix().Data...), nil
}

func unitVector(size, index int) []float64 {
	vector := make([]float64, size)
	vector[index] = 1
	return vector
}
