//ff:func feature=rule type=rule control=iteration dimension=2
//ff:what checkPageEaches — 단일 페이지에서 data-each 필드 존재 검증 (TM-7 내부)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkPageEaches(page parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, fb := range page.Fetches {
		errs = append(errs, checkFetchEaches(page.FileName, fb)...)
	}
	return errs
}
