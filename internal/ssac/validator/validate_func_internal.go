//ff:func feature=ssac-validate type=command control=sequence
//ff:what 개별 ServiceFunc의 내부 정합성을 검증하는 디스패처
package validator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func validateFunc(sf parser.ServiceFunc) []ValidationError {
	var errs []ValidationError
	errs = append(errs, validateRequiredFields(sf)...)
	errs = append(errs, validateVariableFlow(sf)...)
	errs = append(errs, validateStaleResponse(sf)...)
	errs = append(errs, validateReservedSourceConflict(sf)...)
	errs = append(errs, validateSubscribeRules(sf)...)
	errs = append(errs, validateFKReferenceGuard(sf)...)
	errs = append(errs, validateErrStatus(sf)...)
	return errs
}
