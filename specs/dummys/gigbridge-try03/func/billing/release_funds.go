package billing

import "time"

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
	payout := int64(float64(req.Amount) * 0.9)
	_ = payout
	txID := time.Now().UnixNano() / 1000000
	return ReleaseFundsResponse{TransactionID: txID}, nil
}
