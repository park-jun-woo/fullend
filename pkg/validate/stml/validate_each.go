//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateEach — data-each 필드가 배열인지 검증 (TM-7)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateEach(pages []parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		errs = append(errs, checkPageEaches(page)...)
	}
	return errs
}
