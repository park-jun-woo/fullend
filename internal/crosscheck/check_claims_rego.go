//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what Rego input.claims 참조가 fullend.yaml claims 값과 일치하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/policy"
	"github.com/geul-org/fullend/internal/projectconfig"
)

// CheckClaimsRego validates that Rego input.claims.xxx references match fullend.yaml claims values.
func CheckClaimsRego(policies []*policy.Policy, claims map[string]projectconfig.ClaimDef) []CrossError {
	if claims == nil {
		return nil
	}

	claimValues := buildClaimValues(claims)
	regoRefs := collectRegoClaimsRefs(policies)

	var errs []CrossError

	// Forward: Rego claims ref → fullend.yaml claims value
	for ref, file := range regoRefs {
		if !claimValues[ref] {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ Rego",
				Context: file,
				Message: fmt.Sprintf("Rego input.claims.%s — fullend.yaml claims 값에 %q 없음", ref, ref),
				Level:   "ERROR",
			})
		}
	}

	// Reverse: fullend.yaml claim keys → Rego (WARNING if unused)
	for _, def := range claims {
		if _, used := regoRefs[def.Key]; !used {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ Rego",
				Context: "fullend.yaml",
				Message: fmt.Sprintf("claims 값 %q — Rego에서 미참조", def.Key),
				Level:   "WARNING",
			})
		}
	}

	return errs
}
