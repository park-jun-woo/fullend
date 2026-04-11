//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkPagePaginates — 단일 페이지에서 data-paginate x-pagination 존재 검증 (TM-9 내부)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkPagePaginates(page parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, fb := range page.Fetches {
		if fb.Paginate && !hasPaginationExt(ground, fb.OperationID) {
			errs = append(errs, validate.ValidationError{
				Rule: "TM-9", File: page.FileName, Func: fb.OperationID,
				SeqIdx: -1, Level: "ERROR",
				Message: "data-paginate used but OpenAPI has no x-pagination for " + fb.OperationID,
			})
		}
	}
	return errs
}
