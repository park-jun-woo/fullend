//ff:func feature=crosscheck type=rule control=sequence topic=config-check
//ff:what TestCheckJWTBuiltinInputs_NonJWTSkipped: JWT가 아닌 @call은 JWT 검증을 건너뜀 확인
package crosscheck

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestCheckJWTBuiltinInputs_NonJWTSkipped(t *testing.T) {
	claims := map[string]projectconfig.ClaimDef{
		"ID": {Key: "user_id", GoType: "int64"},
	}

	sfs := []ssacparser.ServiceFunc{{
		Name: "Handler",
		Sequences: []ssacparser.Sequence{{
			Type:  "call",
			Model: "billing.CheckCredits",
			Inputs: map[string]string{
				"Balance": "org.CreditsBalance",
			},
		}},
	}}

	errs := CheckJWTBuiltinInputs(sfs, claims)
	if len(errs) != 0 {
		t.Errorf("non-jwt call should be skipped, got: %+v", errs)
	}
}
