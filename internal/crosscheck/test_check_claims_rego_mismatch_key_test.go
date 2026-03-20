//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what CheckClaimsRegoMismatchKey: Rego claims 키 불일치 시 에러를 생성하는지 테스트
package crosscheck

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func TestCheckClaimsRego_MismatchKey(t *testing.T) {
	policies := []*policy.Policy{{
		File:       "authz.rego",
		ClaimsRefs: []string{"user_id", "role"},
	}}
	// user_id -> userId: Rego still references user_id
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
