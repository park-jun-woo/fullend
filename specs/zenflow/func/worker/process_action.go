package worker

import "fmt"

// @func processAction
// @description Simulates processing all actions for a workflow

type ProcessActionRequest struct {
	WorkflowID int64
}

type ProcessActionResponse struct {
	ProcessedCount int64
}

func ProcessAction(req ProcessActionRequest) (ProcessActionResponse, error) {
	if req.WorkflowID <= 0 {
		return ProcessActionResponse{}, fmt.Errorf("invalid workflow ID")
	}
	processed := int64(1)
	return ProcessActionResponse{ProcessedCount: processed}, nil
}
