package worker

// @func processAction
// @description Simulates processing a workflow action via external API

type ProcessActionRequest struct {
	ActionType      string
	PayloadTemplate string
}

type ProcessActionResponse struct {
	Result string
}

func ProcessAction(req ProcessActionRequest) (ProcessActionResponse, error) {
	result := "processed:" + req.ActionType + ":" + req.PayloadTemplate
	return ProcessActionResponse{Result: result}, nil
}
