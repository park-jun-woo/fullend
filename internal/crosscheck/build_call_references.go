//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what SSaC @call 시퀀스에서 참조된 func 이름을 수집
package crosscheck

import (
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func buildCallReferences(funcs []ssacparser.ServiceFunc) map[string]bool {
	referenced := make(map[string]bool)
	for _, fn := range funcs {
		collectCallRefsFromFunc(fn, referenced)
	}
	return referenced
}
