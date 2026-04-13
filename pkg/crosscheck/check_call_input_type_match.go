//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallInputTypeMatch — @call input type ↔ FuncRequest field type (X-44)
package crosscheck

import (

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallInputTypeMatch(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkCallInputTypeMatchSeqs(g, fn.Name, fn.Sequences)...)
	}
	return errs
}
