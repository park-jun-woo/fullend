package worker

import "fmt"

// @func processAction
// @description Simulates processing workflow actions in sequence

type ProcessActionRequest struct {
	WorkflowID int64
}

type ProcessActionResponse struct {
}

func ProcessAction(req ProcessActionRequest) (ProcessActionResponse, error) {
	if req.WorkflowID == 0 {
		return ProcessActionResponse{}, fmt.Errorf("workflow ID is required")
	}
	return ProcessActionResponse{}, nil
}
