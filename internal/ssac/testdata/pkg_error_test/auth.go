package auth

// VerifyPasswordRequest is the request for VerifyPassword.
type VerifyPasswordRequest struct {
	Email    string
	Password string
}

// VerifyPasswordResponse is the response for VerifyPassword.
type VerifyPasswordResponse struct {
	Token string
}

// @error 401
func VerifyPassword(req VerifyPasswordRequest) (VerifyPasswordResponse, error) {
	return VerifyPasswordResponse{}, nil
}

// ChargeRequest is the request for Charge.
type ChargeRequest struct {
	Amount int
}

// No @error annotation — should default to 0.
func Charge(req ChargeRequest) error {
	return nil
}
