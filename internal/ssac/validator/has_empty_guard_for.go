//ff:func feature=ssac-validate type=util control=iteration dimension=1 topic=type-resolve
//ff:what 시퀀스 슬라이스에서 지정한 변수에 대한 @empty 가드가 있는지 확인한다
package validator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

// hasEmptyGuardFor는 시퀀스 슬라이스에서 지정한 변수에 대한 @empty 가드가 있는지 확인한다.
func hasEmptyGuardFor(seqs []parser.Sequence, varName string) bool {
	for _, s := range seqs {
		if s.Type == parser.SeqEmpty && rootVar(s.Target) == varName {
			return true
		}
		// @exists는 "있으면 탈출" — 통과 후 변수는 nil 확정이므로 이후 필드 접근 없음
		if s.Type == parser.SeqExists && rootVar(s.Target) == varName {
			return true
		}
	}
	return false
}
