//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what SSaC currentUser 필드 참조가 fullend.yaml claims에 정의되어 있는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/projectconfig"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckClaims validates that all currentUser field references in SSaC specs
// are defined in fullend.yaml backend.auth.claims.
func CheckClaims(serviceFuncs []ssacparser.ServiceFunc, claims map[string]projectconfig.ClaimDef) []CrossError {
	var errs []CrossError

	usedFields := collectCurrentUserFields(serviceFuncs)
	hasAuth := hasAuthSequence(serviceFuncs)

	// If currentUser is used but no claims config exists → ERROR.
	if (len(usedFields) > 0 || hasAuth) && claims == nil {
		errs = append(errs, CrossError{
			Rule:    "Claims ↔ SSaC",
			Context: "fullend.yaml",
			Message: "currentUser를 사용하지만 backend.auth.claims가 정의되지 않았습니다",
			Level:   "ERROR",
		})
		return errs
	}

	if claims == nil {
		return errs
	}

	for field, locations := range usedFields {
		if _, ok := claims[field]; !ok {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ SSaC",
				Context: strings.Join(locations, ", "),
				Message: fmt.Sprintf("currentUser.%s — claims에 미정의", field),
				Level:   "ERROR",
			})
		}
	}

	return errs
}
