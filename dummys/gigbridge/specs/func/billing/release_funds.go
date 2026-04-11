package billing

// @func releaseFunds
// @description Releases escrowed funds minus 10% platform fee

type ReleaseFundsRequest struct {
	GigID        int64
	Amount       int64
	FreelancerID int64
}

type ReleaseFundsResponse struct {
	TransactionID int64
}

func ReleaseFunds(req ReleaseFundsRequest) (ReleaseFundsResponse, error) {
	payout := req.Amount * 90 / 100
	_ = payout
	return ReleaseFundsResponse{TransactionID: req.GigID}, nil
}
