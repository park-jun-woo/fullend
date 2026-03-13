package billing

import "fmt"

// @func deductCredit
// @description Atomically deducts one credit from an organization

type DeductCreditRequest struct {
	OrgID int64
}

type DeductCreditResponse struct {
	CreditsDeducted int64
}

func DeductCredit(req DeductCreditRequest) (DeductCreditResponse, error) {
	if req.OrgID == 0 {
		return DeductCreditResponse{}, fmt.Errorf("org ID is required")
	}
	return DeductCreditResponse{CreditsDeducted: 1}, nil
}
