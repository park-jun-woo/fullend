//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what CheckClaimsRegoUnusedClaim: 미사용 claims에 대해 WARNING을 생성하는지 테스트
package crosscheck

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

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
