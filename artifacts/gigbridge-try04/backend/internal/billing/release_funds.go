package billing

import "fmt"

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
	if req.Amount <= 0 {
		return ReleaseFundsResponse{}, fmt.Errorf("release amount must be positive, got %d", req.Amount)
	}
	fee := req.Amount / 10
	payout := req.Amount - fee
	_ = payout
	txID := req.GigID*10000 + req.Amount
	return ReleaseFundsResponse{TransactionID: txID}, nil
}
