//ff:func feature=crosscheck type=test control=sequence topic=config-check
//ff:what CheckRoles: 모든 역할 일치 시 에러 없음 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func TestCheckRoles_AllMatch(t *testing.T) {
	policies := []*policy.Policy{{
		File: "authz.rego",
		Rules: []policy.AllowRule{
			{Actions: []string{"CreateGig"}, Resource: "gig", RoleValue: "client", SourceLine: 10},
			{Actions: []string{"SubmitProposal"}, Resource: "gig", RoleValue: "freelancer", SourceLine: 20},
		},
	}}
	roles := []string{"client", "freelancer"}

	errs := CheckRoles(policies, roles)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}
