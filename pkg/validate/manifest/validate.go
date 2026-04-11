//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — fullend.yaml 로드 검증 (C-1)
package manifest

import (
	"github.com/park-jun-woo/fullend/pkg/diagnostic"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks if fullend.yaml loaded successfully.
func Validate(diags []diagnostic.Diagnostic) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, d := range diags {
		errs = append(errs, validate.ValidationError{
			Rule: "C-1", File: d.File, Level: "ERROR",
			Message: "fullend.yaml load error: " + d.Message, SeqIdx: -1,
		})
	}
	return errs
}
