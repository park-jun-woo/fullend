package billing

import "time"

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
	txID := time.Now().UnixNano() / 1000000
	_ = req.Amount
	_ = req.ClientID
	return HoldEscrowResponse{TransactionID: txID}, nil
}
