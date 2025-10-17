package tableau_algorithm1

type VariableSelectionError struct {
	EnteringVarIndex int
	ExitingVarIndex  int
}

func (e VariableSelectionError) Error() string {
	if e.ExitingVarIndex == -1 {
		return "VariableSelectionError: No exiting variable found, problem can not be improved."
	}

	if e.EnteringVarIndex == -1 {
		return "VariableSelectionError: No entering variable found, current solution is optimal."
	}

	return "VariableSelectionError: Unknown variable selection error."
}
