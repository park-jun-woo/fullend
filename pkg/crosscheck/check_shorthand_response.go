//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkShorthandResponse — shorthand @response 검증 (X-19, X-20)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkShorthandResponse(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		if fn.Subscribe != nil {
			continue
		}
		errs = append(errs, checkFuncShorthandResponse(g, fn.Name, fn.Sequences)...)
	}
	return errs
}
