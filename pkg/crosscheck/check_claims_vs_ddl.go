//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what checkClaimsVsDDL — fullend.yaml claims GoType 이 DDL 매핑 컬럼 타입과 일치 (X-74)

package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkClaimsVsDDL verifies that each claim's GoType (default "string") matches
// the Go type of its mapped DDL column. ERROR if mismatch.
func checkClaimsVsDDL(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.Manifest == nil || fs.Manifest.Backend.Auth == nil {
		return nil
	}
	var errs []CrossError
	for fieldName, claim := range fs.Manifest.Backend.Auth.Claims {
		ddlType, tableCol, ok := resolveClaimDDLColumn(claim.Key, g)
		if !ok {
			continue
		}
		claimType := claim.GoType
		if claimType == "" {
			claimType = "string"
		}
		if claimType != ddlType {
			errs = append(errs, CrossError{
				Rule:       "X-74",
				Context:    fmt.Sprintf("claims.%s (JWT key=%q)", fieldName, claim.Key),
				Level:      "ERROR",
				Message:    fmt.Sprintf("claim type %s mismatches DDL column %s type %s", claimType, tableCol, ddlType),
				Suggestion: fmt.Sprintf("fullend.yaml 에 %q 뒤에 :%s 추가 (예: %s:%s)", claim.Key, ddlType, claim.Key, ddlType),
			})
		}
	}
	return errs
}
