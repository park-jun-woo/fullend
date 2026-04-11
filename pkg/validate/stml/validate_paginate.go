//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validatePaginate — data-paginate 사용 시 x-pagination 존재 검증 (TM-9)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validatePaginate(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		errs = append(errs, checkPagePaginates(page, ground)...)
	}
	return errs
}
