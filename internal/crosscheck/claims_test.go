package crosscheck

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func TestCheckClaimsRego_AllMatch(t *testing.T) {
	policies := []*policy.Policy{{
		File:       "authz.rego",
		ClaimsRefs: []string{"user_id", "role"},
	}}
	claims := map[string]projectconfig.ClaimDef{
		"ID":   {Key: "user_id", GoType: "int64"},
		"Role": {Key: "role", GoType: "string"},
	}

	errs := CheckClaimsRego(policies, claims)
	for _, e := range errs {
		if e.Level != "WARNING" {
			t.Errorf("unexpected ERROR: %s", e.Message)
		}
	}
}

func TestCheckClaimsRego_MismatchKey(t *testing.T) {
	policies := []*policy.Policy{{
		File:       "authz.rego",
		ClaimsRefs: []string{"user_id", "role"},
	}}
	// user_id → userId: Rego still references user_id
	claims := map[string]projectconfig.ClaimDef{
		"ID":   {Key: "userId", GoType: "int64"},
		"Role": {Key: "role", GoType: "string"},
	}

	errs := CheckClaimsRego(policies, claims)
	hasError := false
	for _, e := range errs {
		if e.Level == "ERROR" && strings.Contains(e.Message, "user_id") {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected ERROR for Rego input.claims.user_id not in claims values")
	}
}

func TestCheckClaimsRego_UnusedClaim(t *testing.T) {
	policies := []*policy.Policy{{
		File:       "authz.rego",
		ClaimsRefs: []string{"user_id", "role"},
	}}
	claims := map[string]projectconfig.ClaimDef{
		"ID":    {Key: "user_id", GoType: "int64"},
		"Role":  {Key: "role", GoType: "string"},
		"Email": {Key: "email", GoType: "string"},
	}

	errs := CheckClaimsRego(policies, claims)
	hasWarning := false
	for _, e := range errs {
		if e.Level == "WARNING" && strings.Contains(e.Message, "email") {
			hasWarning = true
		}
	}
	if !hasWarning {
		t.Error("expected WARNING for unused claims value 'email'")
	}
}

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
