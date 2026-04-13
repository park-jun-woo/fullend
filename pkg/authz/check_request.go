//ff:type feature=pkg-authz type=model
//ff:what 인가 검사 요청 구조체
package authz

// CheckRequest holds the inputs for an authorization check.
//
// Claims map (when non-nil) passes arbitrary JWT claims to OPA under input.claims.
// When nil, Check() falls back to {"user_id": UserID, "role": Role} for backward compat.
// Generator is expected to populate Claims from CurrentUser fields per manifest.auth.claims config.
type CheckRequest struct {
	Action     string
	Resource   string
	UserID     int64
	Role       string
	ResourceID int64
	Claims     map[string]any
}
