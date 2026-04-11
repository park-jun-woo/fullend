//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallTypeMatch — @call input type ↔ FuncRequest field type (X-44)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallTypeMatch(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkCallTypeMatchSeqs(g, fn.Name, fn.Sequences)...)
	}
	return errs
}
