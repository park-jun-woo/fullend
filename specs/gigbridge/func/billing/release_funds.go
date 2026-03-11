package billing

// @func releaseFunds
// @description Calculates 10% platform fee and releases 90% to freelancer

type ReleaseFundsRequest struct {
	GigID        int64
	Amount       int
	FreelancerID int64
}

type ReleaseFundsResponse struct {
	TransactionID int64
	Payout        int
	Fee           int
}

func ReleaseFunds(req ReleaseFundsRequest) (ReleaseFundsResponse, error) {
	fee := req.Amount / 10
	payout := req.Amount - fee
	txID := req.GigID*1000 + int64(payout)
	return ReleaseFundsResponse{TransactionID: txID, Payout: payout, Fee: fee}, nil
}
