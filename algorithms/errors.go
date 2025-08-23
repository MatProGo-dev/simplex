package algorithms

import "fmt"

type IterationCounter interface {
	NumberOfIterations() int
}

type IterationCountIsNegativeError struct {
	counter IterationCounter
}

func MakeIterationCountIsNegativeError(counter IterationCounter) IterationCountIsNegativeError {
	return IterationCountIsNegativeError{counter}
}

func (err IterationCountIsNegativeError) Error() string {
	return fmt.Sprintf(
		"The number of iterations (%v) in the algorithm is a negative number.",
		err.counter.NumberOfIterations(),
	)
}
