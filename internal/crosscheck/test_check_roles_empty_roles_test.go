//ff:func feature=crosscheck type=test control=sequence topic=config-check
//ff:what CheckRoles: 역할 미정의 시 에러 없이 스킵 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func TestCheckRoles_EmptyRoles(t *testing.T) {
	// No roles defined — should return no errors (skip).
	policies := []*policy.Policy{{
		File: "authz.rego",
		Rules: []policy.AllowRule{
			{Actions: []string{"CreateGig"}, Resource: "gig", RoleValue: "client", SourceLine: 10},
		},
	}}

	errs := CheckRoles(policies, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for nil roles, got %d", len(errs))
	}
}
