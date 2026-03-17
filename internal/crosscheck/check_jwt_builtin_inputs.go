//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what jwtBuiltinFuncs 대상 @call의 input key가 claims 필드와 일치하는지 검증
package crosscheck

import (
	"github.com/park-jun-woo/fullend/internal/projectconfig"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// CheckJWTBuiltinInputs validates that @call inputs for jwt builtin functions
// use keys that match claims field names.
func CheckJWTBuiltinInputs(serviceFuncs []ssacparser.ServiceFunc, claims map[string]projectconfig.ClaimDef) []CrossError {
	if claims == nil {
		return nil
	}
	claimFields := make(map[string]bool, len(claims))
	for field := range claims {
		claimFields[field] = true
	}
	var errs []CrossError
	for _, sf := range serviceFuncs {
		errs = append(errs, checkJWTInputsForFunc(sf, claimFields, claims)...)
	}
	return errs
}
