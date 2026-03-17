package worker

// @func processActions
// @description Simulates executing workflow actions in sequence

type ProcessActionsRequest struct {
	WorkflowID  int64
	ActionCount int64
}

type ProcessActionsResponse struct {
	CreditsUsed int64
	Success     bool
}

func ProcessActions(req ProcessActionsRequest) (ProcessActionsResponse, error) {
	creditsUsed := int64(1)
	if req.ActionCount > 5 {
		creditsUsed = 2
	}
	return ProcessActionsResponse{
		CreditsUsed: creditsUsed,
		Success:     true,
	}, nil
}
