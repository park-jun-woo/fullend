package crosscheck

import (
	"testing"

	"github.com/geul-org/fullend/internal/policy"
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
