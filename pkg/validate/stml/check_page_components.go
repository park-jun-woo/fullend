//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what checkPageComponents — 단일 페이지에서 data-component 파일 존재 검증 (TM-12 내부)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkPageComponents(page parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, comp := range page.Children {
		if comp.Kind == "component" && comp.Component != nil && comp.Component.Name != "" {
			errs = append(errs, validate.ValidationError{
				Rule: "TM-12", File: page.FileName, Func: comp.Component.Name,
				SeqIdx: -1, Level: "WARNING",
				Message: "data-component " + comp.Component.Name + " — verify .tsx file exists",
			})
		}
	}
	return errs
}
