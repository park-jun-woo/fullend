package billing

import "fmt"

// @func validateCredits
// @error 402
// @description Validates that organization has sufficient credits

type ValidateCreditsRequest struct {
	CreditsBalance int64
}

type ValidateCreditsResponse struct{}

func ValidateCredits(req ValidateCreditsRequest) (ValidateCreditsResponse, error) {
	if req.CreditsBalance <= 0 {
		return ValidateCreditsResponse{}, fmt.Errorf("insufficient credits: balance is %d", req.CreditsBalance)
	}
	return ValidateCreditsResponse{}, nil
}
