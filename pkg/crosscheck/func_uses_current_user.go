//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what funcUsesCurrentUser — 시퀀스의 Args에서 currentUser 참조 여부 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func funcUsesCurrentUser(seqs []ssac.Sequence) bool {
	for _, seq := range seqs {
		if argUsesCurrentUser(seq.Args) {
			return true
		}
	}
	return false
}
