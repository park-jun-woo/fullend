//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCallFuncName — @call 함수명 소문자 시작 검증 (X-38)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkCallFuncName(fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkCallFuncNameSeqs(fn.Name, fn.Sequences)...)
	}
	return errs
}
