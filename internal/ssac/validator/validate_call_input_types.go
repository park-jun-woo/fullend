//ff:func feature=ssac-validate type=rule control=iteration dimension=2 topic=type-resolve
//ff:what @call inputs 필드 타입과 func Request struct 필드 타입 비교

package validator

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validateCallInputTypes는 @call inputs의 필드 타입을 func Request struct 필드 타입과 비교한다.
func validateCallInputTypes(sf parser.ServiceFunc, st *SymbolTable) []ValidationError {
	if st == nil {
		return nil
	}

	// result 변수 → 모델명 매핑 (DDL 타입 추적용)
	resultModels := map[string]string{}
	for _, seq := range sf.Sequences {
		if seq.Result == nil || seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		resultModels[seq.Result.Var] = parts[0]
	}

	var errs []ValidationError
	for i, seq := range sf.Sequences {
		if seq.Type != parser.SeqCall || seq.Model == "" {
			continue
		}
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 {
			continue
		}
		pkgName, funcName := parts[0], parts[1]

		// pkg.Model 키로 심볼 테이블에서 ParamTypes 조회
		// @call은 func이므로 모델 키를 찾아야 함
		var paramTypes map[string]string
		for modelKey, ms := range st.Models {
			if !strings.HasPrefix(modelKey, pkgName+".") {
				continue
			}
			if mi, ok := ms.Methods[funcName]; ok && mi.ParamTypes != nil {
				paramTypes = mi.ParamTypes
				break
			}
		}
		if paramTypes == nil {
			continue // Request struct가 파싱되지 않았으면 스킵
		}

		ctx := errCtx{sf.FileName, sf.Name, i}
		for key, val := range seq.Inputs {
			expectedType, ok := paramTypes[key]
			if !ok {
				continue // 필드가 Request struct에 없으면 다른 검증에서 처리
			}
			actualType := resolveCallInputType(val, resultModels, st)
			if actualType == "" || actualType == expectedType {
				continue
			}
			errs = append(errs, ctx.err("@call", fmt.Sprintf("입력 %q의 타입 불일치: %s은 %s, Request 필드는 %s", key, val, actualType, expectedType)))
		}
	}
	return errs
}
