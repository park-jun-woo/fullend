//ff:func feature=ssac-validate type=rule control=iteration dimension=1 topic=type-resolve
//ff:what 패키지 접두사 모델의 존재 및 파라미터 매칭을 검증한다
package validator

import (
	"fmt"
	"sort"
	"strings"

	"github.com/geul-org/fullend/internal/ssac/parser"
)

// validatePackageModel은 패키지 접두사 모델의 존재 및 파라미터 매칭을 검증한다.
func validatePackageModel(ctx errCtx, sfName string, seq parser.Sequence, modelName, methodName string, st *SymbolTable) []ValidationError {
	pkgModelKey := seq.Package + "." + modelName
	ms, ok := st.Models[pkgModelKey]
	if !ok {
		return []ValidationError{{
			FileName: ctx.fileName, FuncName: ctx.funcName, SeqIndex: ctx.seqIndex,
			Tag: "@" + seq.Type, Message: fmt.Sprintf("%s.%s — 패키지 interface를 찾을 수 없습니다. 검증을 건너뜁니다", seq.Package, modelName), Level: "WARNING",
		}}
	}
	if !ms.HasMethod(methodName) {
		available := make([]string, 0, len(ms.Methods))
		for m := range ms.Methods {
			available = append(available, m)
		}
		sort.Strings(available)
		return []ValidationError{ctx.err("@"+seq.Type, fmt.Sprintf("%s.%s — 메서드 %q 없음. 사용 가능: %s", seq.Package, modelName, methodName, strings.Join(available, ", ")))}
	}
	mi := ms.Methods[methodName]
	if len(mi.Params) == 0 {
		return nil
	}

	ifaceParamSet := make(map[string]bool, len(mi.Params))
	for _, p := range mi.Params {
		ifaceParamSet[p] = true
	}
	ssacKeys := make([]string, 0, len(seq.Inputs))
	for k := range seq.Inputs {
		ssacKeys = append(ssacKeys, k)
	}
	sort.Strings(ssacKeys)

	var errs []ValidationError
	// SSaC에 있지만 interface에 없는 키
	for _, key := range ssacKeys {
		if ifaceParamSet[key] {
			continue
		}
		errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("SSaC ↔ Interface: %s — @model %s.%s.%s 파라미터 불일치. SSaC에 %q가 있지만 interface에 없습니다. interface 파라미터: [%s]", sfName, seq.Package, modelName, methodName, key, strings.Join(mi.Params, ", "))))
	}
	// interface에 있지만 SSaC에 없는 파라미터
	ssacKeySet := make(map[string]bool, len(seq.Inputs))
	for k := range seq.Inputs {
		ssacKeySet[k] = true
	}
	for _, param := range mi.Params {
		if ssacKeySet[param] {
			continue
		}
		errs = append(errs, ctx.err("@"+seq.Type, fmt.Sprintf("SSaC ↔ Interface: %s — @model %s.%s.%s 파라미터 누락. interface에 %q가 필요하지만 SSaC에 없습니다. SSaC 파라미터: [%s]", sfName, seq.Package, modelName, methodName, param, strings.Join(ssacKeys, ", "))))
	}
	return errs
}
