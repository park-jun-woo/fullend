//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateBindNotFound — data-bind 필드가 OpenAPI response에도 custom.ts에도 없는 경우 (TM-8)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateBindNotFound(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		for _, fb := range page.Fetches {
			errs = append(errs, checkBindNotFound(fb, page.FileName, ground)...)
		}
	}
	return errs
}
