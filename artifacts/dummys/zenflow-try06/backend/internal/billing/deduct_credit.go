package billing

import "fmt"

// @func deductCredit
// @description Calculates new balance after credit deduction

type DeductCreditRequest struct {
	OrgID          int64
	Amount         int64
	CurrentBalance int64
}

type DeductCreditResponse struct {
	NewBalance  int64
	CreditsUsed int64
}

func DeductCredit(req DeductCreditRequest) (DeductCreditResponse, error) {
	if req.Amount <= 0 {
		return DeductCreditResponse{}, fmt.Errorf("deduction amount must be positive, got %d", req.Amount)
	}
	newBalance := req.CurrentBalance - req.Amount
	return DeductCreditResponse{
		NewBalance:  newBalance,
		CreditsUsed: req.Amount,
	}, nil
}
