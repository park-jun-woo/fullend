//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validatePage — 단일 STML 페이지의 fetch/action operationId 검증
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validatePage(graph *toulmin.Graph, ground *rule.Ground, page parsestml.PageSpec) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, fb := range page.Fetches {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", fb.OperationID)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, page.FileName, fb.OperationID)...)
	}
	for _, ab := range page.Actions {
		ctx := toulmin.NewContext()
		ctx.Set("ground", ground)
		ctx.Set("claim", ab.OperationID)
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toSTMLErrors(results, page.FileName, ab.OperationID)...)
	}
	return errs
}
