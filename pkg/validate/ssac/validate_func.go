//ff:func feature=rule type=rule control=sequence
//ff:what validateFunc — 단일 SSaC 함수의 시퀀스별 검증 디스패치
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateFunc(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	errs = append(errs, validateRequiredFields(fn, ground)...)
	errs = append(errs, validateVariableFlow(fn, ground)...)
	errs = append(errs, validateForbiddenRefs(fn, ground)...)
	errs = append(errs, validateNameFormats(fn, ground)...)
	errs = append(errs, validateModelRefs(fn, ground)...)
	errs = append(errs, validateStaleResponse(fn)...)
	errs = append(errs, validateDeleteInputs(fn)...)
	errs = append(errs, validateErrStatus(fn)...)
	errs = append(errs, validateUnknownSeq(fn)...)
	errs = append(errs, validateFKGuard(fn)...)
	errs = append(errs, validateMethodRef(fn, ground)...)
	errs = append(errs, validateRequestRef(fn, ground)...)
	errs = append(errs, validatePagination(fn)...)
	errs = append(errs, validateSubscribeForbidden(fn, ground)...)
	if fn.Subscribe != nil {
		errs = append(errs, validateSubscribe(fn)...)
	}
	return errs
}
