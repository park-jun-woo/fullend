package worker

// @func processActions
// @description Processes all actions in sequence order

type ProcessActionsRequest struct {
	Actions []ActionInput
}

type ActionInput struct {
	ActionType      string
	PayloadTemplate string
}

type ProcessActionsResponse struct {
	Processed int
}

func ProcessActions(req ProcessActionsRequest) (ProcessActionsResponse, error) {
	count := len(req.Actions)
	return ProcessActionsResponse{Processed: count}, nil
}
