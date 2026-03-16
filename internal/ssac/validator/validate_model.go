//ff:func feature=ssac-validate type=rule control=iteration dimension=2
//ff:what Model 심볼 테이블 존재 검증

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateModel은 Model이 심볼 테이블에 존재하는지 검증한다.
func validateModel(sf parser.ServiceFunc, st *SymbolTable) []ValidationError {
	var errs []ValidationError
	for i, seq := range sf.Sequences {
		if seq.Model == "" || seq.Type == parser.SeqCall {
			continue // @call은 외부 패키지이므로 교차검증 스킵
		}
		ctx := errCtx{sf.FileName, sf.Name, i}

		// Result 타입이 대문자로 시작하는지 검증 (Go exported type 규칙)
		if seq.Result != nil && seq.Result.Type != "" && seq.Result.Type[0] >= 'a' && seq.Result.Type[0] <= 'z' {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("Result 타입 %q — 대문자로 시작해야 합니다 (예: %s)", seq.Result.Type, strings.Title(seq.Result.Type))))
		}

		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 {
			continue
		}
		modelName, methodName := parts[0], parts[1]

		if seq.Package != "" {
			errs = append(errs, validatePackageModel(ctx, sf.Name, seq, modelName, methodName, st)...)
			continue
		}

		ms, ok := st.Models[modelName]
		if !ok {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 모델을 찾을 수 없습니다", modelName)))
			continue
		}
		if !ms.HasMethod(methodName) {
			errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%q 모델에 %q 메서드가 없습니다", modelName, methodName)))
		}
	}
	return errs
}
