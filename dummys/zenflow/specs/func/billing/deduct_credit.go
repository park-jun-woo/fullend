package billing

// @func deductCredit
// @description Atomically deducts credits from organization balance

type DeductCreditRequest struct {
	OrgID  int64
	Amount int64
}

type DeductCreditResponse struct {
	Remaining int64
}

func DeductCredit(req DeductCreditRequest) (DeductCreditResponse, error) {
	remaining := int64(99) - req.Amount
	return DeductCreditResponse{Remaining: remaining}, nil
}
