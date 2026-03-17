package workflow

// @func resolveRootID
// @description Resolves the root workflow ID for version grouping

type ResolveRootIDRequest struct {
	WorkflowID     int64
	RootWorkflowID int64
}

type ResolveRootIDResponse struct {
	ResolvedRootID int64
}

func ResolveRootID(req ResolveRootIDRequest) (ResolveRootIDResponse, error) {
	if req.RootWorkflowID == 0 {
		return ResolveRootIDResponse{ResolvedRootID: req.WorkflowID}, nil
	}
	return ResolveRootIDResponse{ResolvedRootID: req.RootWorkflowID}, nil
}
