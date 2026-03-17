//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what input key와 param 목록을 대소문자 무시 비교하여 불일치 반환
package crosscheck

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func matchInputKeysToParams(ctx string, seqIdx int, seq ssacparser.Sequence, params []string) []CrossError {
	paramSet := make(map[string]bool, len(params))
	for _, p := range params {
		paramSet[p] = true
	}
	var errs []CrossError
	for key := range seq.Inputs {
		if paramSet[key] {
			continue
		}
		if matched := findCaseInsensitiveParam(key, params); matched != "" {
			errs = append(errs, CrossError{
				Rule:       "SSaC input key case",
				Context:    fmt.Sprintf("%s seq[%d]", ctx, seqIdx),
				Message:    fmt.Sprintf("input key %q와 sqlc 파라미터 %q — 대소문자 불일치 (Go initialism 확인 필요)", key, matched),
				Level:      "WARNING",
				Suggestion: fmt.Sprintf("input key를 %q로 변경하세요", matched),
			})
		}
	}
	return errs
}
