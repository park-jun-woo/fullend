//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkJWTClaimArgs — @call args에서 currentUser 소스 필드 → Config claims 검증 (X-73)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func checkJWTClaimArgs(graph *toulmin.Graph, g *rule.Ground, funcName string, args []ssac.Arg) []CrossError {
	var errs []CrossError
	for _, arg := range args {
		if arg.Source == "currentUser" && arg.Field != "" {
			errs = append(errs, evalRef(graph, g, arg.Field, funcName)...)
		}
	}
	return errs
}
