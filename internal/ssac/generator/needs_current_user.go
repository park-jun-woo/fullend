//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 시퀀스에서 currentUser 참조가 필요한지 확인
package generator

import "github.com/geul-org/fullend/internal/ssac/parser"

func needsCurrentUser(seqs []parser.Sequence) bool {
	for _, seq := range seqs {
		if seq.Type == parser.SeqAuth {
			return true
		}
		if hasCurrentUserArg(seq) {
			return true
		}
		if hasCurrentUserInput(seq) {
			return true
		}
	}
	return false
}
