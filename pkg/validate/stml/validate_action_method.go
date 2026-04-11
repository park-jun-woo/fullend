//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateActionMethod — data-action은 GET이 아닌 메서드만 허용 (TM-3)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateActionMethod(pages []parsestml.PageSpec) []validate.ValidationError {
	// TM-3: action operationId must map to non-GET method.
	// This requires OpenAPI method lookup — deferred to crosscheck X-15 level.
	// At STML-only validation, we just ensure operationId is not empty.
	var errs []validate.ValidationError
	for _, page := range pages {
		errs = append(errs, checkPageActionMethods(page)...)
	}
	return errs
}
