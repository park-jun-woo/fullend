//ff:func feature=crosscheck type=test control=iteration dimension=1 topic=config-check
//ff:what CheckRoles: Rego 역할값 대소문자 불일치 시 에러 검증

package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
)

func TestCheckRoles_RegoMismatch(t *testing.T) {
	// Rego has "Client" (capital C) but roles says "client"
	policies := []*policy.Policy{{
		File: "authz.rego",
		Rules: []policy.AllowRule{
			{Actions: []string{"PublishGig"}, Resource: "gig", RoleValue: "Client", SourceLine: 10},
			{Actions: []string{"SubmitProposal"}, Resource: "gig", RoleValue: "freelancer", SourceLine: 20},
		},
	}}
	roles := []string{"client", "freelancer"}

	errs := CheckRoles(policies, roles)
	foundErr := false
	for _, e := range errs {
		if e.Level != "WARNING" && contains(e.Message, "Client") {
			foundErr = true
		}
	}
	if !foundErr {
		t.Errorf("expected ERROR for 'Client' mismatch, got: %+v", errs)
	}
}
