package webhook

import "fmt"

// @func deliver
// @description Simulates delivering webhook payload to registered URLs

type DeliverRequest struct {
	OrgID      int64
	WorkflowID int64
	Status     string
}

type DeliverResponse struct {
	Delivered bool
}

func Deliver(req DeliverRequest) (DeliverResponse, error) {
	if req.OrgID <= 0 {
		return DeliverResponse{}, fmt.Errorf("invalid org ID: %d", req.OrgID)
	}
	if req.WorkflowID <= 0 {
		return DeliverResponse{}, fmt.Errorf("invalid workflow ID: %d", req.WorkflowID)
	}
	return DeliverResponse{Delivered: true}, nil
}
