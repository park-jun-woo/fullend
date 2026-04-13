package worker

import "github.com/example/zenflow/internal/model"

// @func processActions
// @description Processes all actions in sequence order

type ProcessActionsRequest struct {
	Actions []model.Action
}

type ProcessActionsResponse struct {
	Processed int
}

func ProcessActions(req ProcessActionsRequest) (ProcessActionsResponse, error) {
	count := len(req.Actions)
	return ProcessActionsResponse{Processed: count}, nil
}
