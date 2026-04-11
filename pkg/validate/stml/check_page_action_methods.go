//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkPageActionMethods — 페이지의 action에 operationId 존재 검증
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkPageActionMethods(page parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, ab := range page.Actions {
		if ab.OperationID == "" {
			errs = append(errs, validate.ValidationError{
				Rule: "TM-3", File: page.FileName, Func: "action",
				SeqIdx: -1, Level: "ERROR",
				Message: "data-action requires operationId",
			})
		}
	}
	return errs
}
