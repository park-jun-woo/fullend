//ff:func feature=ssac-validate type=rule
//ff:what Model 심볼 테이블 존재 검증

package validator

import (
	"fmt"
	"sort"
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
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) < 2 {
			continue
		}
		modelName, methodName := parts[0], parts[1]

		if seq.Package != "" {
			// 패키지 접두사 모델: "pkg.Model" 키로 조회
			pkgModelKey := seq.Package + "." + modelName
			ms, ok := st.Models[pkgModelKey]
			if !ok {
				// interface가 없으면 WARNING (검증 불가)
				errs = append(errs, ValidationError{
					FileName: ctx.fileName, FuncName: ctx.funcName, SeqIndex: ctx.seqIndex,
					Tag: "@" + seq.Type, Message: fmt.Sprintf("%s.%s — 패키지 interface를 찾을 수 없습니다. 검증을 건너뜁니다", seq.Package, modelName), Level: "WARNING",
				})
				continue
			}
			if !ms.HasMethod(methodName) {
				// 사용 가능 메서드 안내
				available := make([]string, 0, len(ms.Methods))
				for m := range ms.Methods {
					available = append(available, m)
				}
				sort.Strings(available)
				errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("%s.%s — 메서드 %q 없음. 사용 가능: %s", seq.Package, modelName, methodName, strings.Join(available, ", "))))
				continue
			}
			// 파라미터 매칭 검증
			mi := ms.Methods[methodName]
			if len(mi.Params) > 0 {
				ifaceParamSet := make(map[string]bool, len(mi.Params))
				for _, p := range mi.Params {
					ifaceParamSet[p] = true
				}
				ssacKeys := make([]string, 0, len(seq.Inputs))
				for k := range seq.Inputs {
					ssacKeys = append(ssacKeys, k)
				}
				sort.Strings(ssacKeys)
				// SSaC에 있지만 interface에 없는 키
				for _, key := range ssacKeys {
					if !ifaceParamSet[key] {
						errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("SSaC ↔ Interface: %s — @model %s.%s.%s 파라미터 불일치. SSaC에 %q가 있지만 interface에 없습니다. interface 파라미터: [%s]", sf.Name, seq.Package, modelName, methodName, key, strings.Join(mi.Params, ", "))))
					}
				}
				// interface에 있지만 SSaC에 없는 파라미터
				ssacKeySet := make(map[string]bool, len(seq.Inputs))
				for k := range seq.Inputs {
					ssacKeySet[k] = true
				}
				for _, param := range mi.Params {
					if !ssacKeySet[param] {
						errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("SSaC ↔ Interface: %s — @model %s.%s.%s 파라미터 누락. interface에 %q가 필요하지만 SSaC에 없습니다. SSaC 파라미터: [%s]", sf.Name, seq.Package, modelName, methodName, param, strings.Join(ssacKeys, ", "))))
					}
				}
			}
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
