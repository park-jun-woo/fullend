//ff:func feature=crosscheck type=rule control=sequence
//ff:what CheckClaimsRegoNilClaims: nil claims에 대해 에러 없이 처리하는지 테스트
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func TestCheckClaimsRego_NilClaims(t *testing.T) {
	policies := []*policy.Policy{{
		File:       "authz.rego",
		ClaimsRefs: []string{"user_id"},
	}}

	errs := CheckClaimsRego(policies, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for nil claims, got %d", len(errs))
	}
}
