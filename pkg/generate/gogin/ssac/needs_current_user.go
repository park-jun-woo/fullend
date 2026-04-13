//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=currentuser
//ff:what 시퀀스에서 currentUser 참조가 필요한지 확인
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func needsCurrentUser(seqs []ssacparser.Sequence) bool {
	for _, seq := range seqs {
		if seq.Type == ssacparser.SeqAuth {
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
