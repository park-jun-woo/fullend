//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateActionFields — data-field → OpenAPI request 필드 검증 (TM-5)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateActionFields(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		for _, ab := range page.Actions {
			errs = append(errs, checkActionFieldRefs(ab, page.FileName, ground)...)
		}
	}
	return errs
}
