//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=config-check
//ff:what CheckRoles: Rego에서 미사용 역할에 대한 WARNING 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func TestCheckRoles_UnusedRole(t *testing.T) {
	// roles has "admin" but Rego never uses it
	policies := []*policy.Policy{{
		File: "authz.rego",
		Rules: []policy.AllowRule{
			{Actions: []string{"CreateGig"}, Resource: "gig", RoleValue: "client", SourceLine: 10},
		},
	}}
	roles := []string{"client", "freelancer", "admin"}

	errs := CheckRoles(policies, roles)
	warnings := 0
	for _, e := range errs {
		if e.Level == "WARNING" {
			warnings++
		}
	}
	if warnings != 2 {
		t.Errorf("expected 2 warnings (freelancer, admin unused), got %d: %+v", warnings, errs)
	}
}
