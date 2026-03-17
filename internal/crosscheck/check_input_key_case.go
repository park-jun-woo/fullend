//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what SSaC input key가 sqlc 메서드 파라미터명과 대소문자까지 일치하는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// CheckInputKeyCase validates that SSaC input keys exactly match sqlc method parameter names (case-sensitive).
func CheckInputKeyCase(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError
	for _, fn := range funcs {
		ctx := fmt.Sprintf("%s:%s", fn.FileName, fn.Name)
		for i, seq := range fn.Sequences {
			if seq.Model == "" || seq.Type == "call" || seq.Package != "" {
				continue
			}
			parts := strings.SplitN(seq.Model, ".", 2)
			if len(parts) < 2 {
				continue
			}
			modelName, methodName := parts[0], parts[1]
			ms, ok := st.Models[modelName]
			if !ok {
				continue
			}
			mi, exists := ms.Methods[methodName]
			if !exists || len(mi.Params) == 0 {
				continue
			}
			paramSet := make(map[string]bool, len(mi.Params))
			for _, p := range mi.Params {
				paramSet[p] = true
			}
			for key := range seq.Inputs {
				if paramSet[key] {
					continue
				}
				// Check case-insensitive match.
				for _, p := range mi.Params {
					if strings.EqualFold(key, p) {
						errs = append(errs, CrossError{
							Rule:       "SSaC input key case",
							Context:    fmt.Sprintf("%s seq[%d]", ctx, i),
							Message:    fmt.Sprintf("input key %q와 sqlc 파라미터 %q — 대소문자 불일치 (Go initialism 확인 필요)", key, p),
							Level:      "WARNING",
							Suggestion: fmt.Sprintf("input key를 %q로 변경하세요", p),
						})
						break
					}
				}
			}
		}
	}
	return errs
}
