//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallArgTypeMatch — 단일 @call의 arg 타입 ↔ FuncRequest 필드 타입 비교 (X-44)
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkCallArgTypeMatch(g *rule.Ground, funcName string, seq ssac.Sequence) []CrossError {
	idx := strings.IndexByte(seq.Model, '.')
	if idx <= 0 {
		return nil
	}
	callFunc := seq.Model[idx+1:]

	graph := toulmin.NewGraph("call-type-" + callFunc)
	graph.Rule(rule.TypeMatch).With(&rule.TypeMatchSpec{
		BaseSpec:  rule.BaseSpec{Rule: "X-44", Level: "ERROR", Message: "@call input type mismatch with func request field"},
		LookupKey: "Func.request." + callFunc,
	})

	var errs []CrossError
	for _, arg := range seq.Args {
		if arg.Source == "" || arg.Field == "" {
			continue
		}
		// Resolve source type from variable tracking
		sourceType := g.Types["SSaC.var."+funcName+"."+arg.Source]
		if sourceType == "" {
			continue
		}
		ctx := toulmin.NewContext()
		ctx.Set("ground", g)
		ctx.Set("claim", &rule.TypeClaim{Name: arg.Field, SourceType: sourceType})
		results, _ := graph.Evaluate(ctx)
		errs = append(errs, toErrors(results, funcName+"/"+seq.Model+"."+arg.Field)...)
	}
	return errs
}
