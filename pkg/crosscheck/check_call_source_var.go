//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallSourceVar — @call arg source 변수 미정의 WARNING (X-47)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

var implicitSources = map[string]bool{
	"request": true, "currentUser": true, "query": true, "message": true, "": true,
}

func checkCallSourceVar(fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		declared := collectDeclaredVars(fn.Sequences)
		errs = append(errs, checkUndeclaredSources(fn.Name, fn.Sequences, declared)...)
	}
	return errs
}
