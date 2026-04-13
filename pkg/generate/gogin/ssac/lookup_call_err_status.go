//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what SymbolTable에서 @call 대상 함수의 @error 어노테이션 값을 조회
package ssac

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// lookupCallErrStatus는 SymbolTable에서 @call 대상 함수의 @error 어노테이션 값을 조회한다.
func lookupCallErrStatus(st *validator.SymbolTable, model string) int {
	if st == nil {
		return 0
	}
	parts := strings.SplitN(model, ".", 2)
	if len(parts) < 2 {
		return 0
	}
	pkgName, funcName := parts[0], parts[1]
	for modelKey, ms := range st.Models {
		if !strings.HasPrefix(modelKey, pkgName+".") {
			continue
		}
		if mi, ok := ms.Methods[funcName]; ok {
			return mi.ErrStatus
		}
	}
	return 0
}
