//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkFetchEaches — 단일 fetch 블록에서 data-each 필드 존재 검증 (TM-7 내부)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkFetchEaches(fileName string, fb parsestml.FetchBlock) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, each := range fb.Eaches {
		if each.Field == "" {
			errs = append(errs, validate.ValidationError{
				Rule: "TM-7", File: fileName, Func: fb.OperationID,
				SeqIdx: -1, Level: "ERROR",
				Message: "data-each requires a field name",
			})
		}
	}
	return errs
}
