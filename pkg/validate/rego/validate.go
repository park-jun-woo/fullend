//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — Rego 파싱 에러 검증 (P-1)
package rego

import (
	"github.com/park-jun-woo/fullend/pkg/diagnostic"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate converts rego parse diagnostics to ValidationErrors.
func Validate(diags []diagnostic.Diagnostic) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, d := range diags {
		errs = append(errs, validate.ValidationError{
			Rule: "P-1", File: d.File, Level: "ERROR",
			Message: d.Message, SeqIdx: -1,
		})
	}
	return errs
}
