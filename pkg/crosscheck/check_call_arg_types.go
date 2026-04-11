//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallArgTypes — 개별 @call의 Arg 타입 ↔ FuncRequest 필드 타입 비교
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallArgTypes(g *rule.Ground, funcName, model string, args []ssac.Arg) []CrossError {
	idx := strings.IndexByte(model, '.')
	if idx <= 0 {
		return nil
	}
	callFunc := model[idx+1:]
	var errs []CrossError
	for _, arg := range args {
		if arg.Field == "" {
			continue
		}
		targetType, ok := g.Types["Func.request."+callFunc+"."+arg.Field]
		if !ok {
			continue
		}
		_ = targetType // Type inference from SSaC args requires symbol table — deferred
	}
	return errs
}
