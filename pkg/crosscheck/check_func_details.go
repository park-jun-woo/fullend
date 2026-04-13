//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncDetails — @call 함수명 형식, input 스키마, result 매칭 (X-38, X-42~X-47)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkFuncDetails(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkFuncDetailSeqs(g, fn.Name, fn.Sequences, fs)...)
	}
	return errs
}
