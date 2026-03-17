package billing

import "fmt"

// @func checkCredits
// @error 402
// @description Checks if organization has sufficient credits

type CheckCreditsRequest struct {
	OrgID   int64
	Balance int64
}

type CheckCreditsResponse struct {
	Available bool
}

func CheckCredits(req CheckCreditsRequest) (CheckCreditsResponse, error) {
	if req.Balance <= 0 {
		return CheckCreditsResponse{Available: false}, fmt.Errorf("insufficient credits: organization %d has %d credits", req.OrgID, req.Balance)
	}
	return CheckCreditsResponse{Available: true}, nil
}
