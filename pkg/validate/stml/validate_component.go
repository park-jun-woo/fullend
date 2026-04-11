//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateComponent — data-component 파일 존재 검증 (TM-12)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateComponent(pages []parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		errs = append(errs, checkPageComponents(page)...)
	}
	return errs
}
