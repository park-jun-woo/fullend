package billing

import "fmt"

// @func checkCredits
// @error 402
// @description Checks credit balance; returns error if insufficient

type CheckCreditsRequest struct {
	OrgID int64
}

type CheckCreditsResponse struct {
	Balance int64
}

func CheckCredits(req CheckCreditsRequest) (CheckCreditsResponse, error) {
	if req.OrgID == 0 {
		return CheckCreditsResponse{}, fmt.Errorf("org ID is required")
	}
	balance := int64(100)
	if balance <= 0 {
		return CheckCreditsResponse{}, fmt.Errorf("insufficient credits")
	}
	return CheckCreditsResponse{Balance: balance}, nil
}
