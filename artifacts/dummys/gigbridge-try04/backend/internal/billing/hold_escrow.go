package billing

import "fmt"

// @func holdEscrow
// @description Simulates locking funds in escrow

type HoldEscrowRequest struct {
	GigID    int64
	Amount   int64
	ClientID int64
}

type HoldEscrowResponse struct {
	TransactionID int64
}

func HoldEscrow(req HoldEscrowRequest) (HoldEscrowResponse, error) {
	if req.Amount <= 0 {
		return HoldEscrowResponse{}, fmt.Errorf("escrow amount must be positive, got %d", req.Amount)
	}
	if req.GigID <= 0 {
		return HoldEscrowResponse{}, fmt.Errorf("invalid gig ID: %d", req.GigID)
	}
	txID := req.GigID*10000 + req.Amount
	return HoldEscrowResponse{TransactionID: txID}, nil
}
