//ff:func feature=ssac-gen type=util control=sequence
//ff:what 심볼 테이블에서 모델 메서드의 파라미터 순서를 조회
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/ssac/validator"
)

// lookupParamOrder는 심볼 테이블에서 모델 메서드의 파라미터 순서를 조회한다.
func lookupParamOrder(model string, st *validator.SymbolTable) []string {
	parts := strings.SplitN(model, ".", 2)
	if len(parts) < 2 {
		return nil
	}
	ms, ok := st.Models[parts[0]]
	if !ok {
		return nil
	}
	mi, ok := ms.Methods[parts[1]]
	if !ok {
		return nil
	}
	return mi.Params
}
