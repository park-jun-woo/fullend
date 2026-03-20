//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what CheckClaimsRegoAllMatch: 모든 claims가 Rego와 일치할 때 에러가 없는지 테스트
package crosscheck

import (
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
