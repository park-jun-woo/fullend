//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what hasStateSeq — 시퀀스 목록에 @state 타입이 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func hasStateSeq(seqs []ssac.Sequence) bool {
	for _, seq := range seqs {
		if seq.Type == "state" {
			return true
		}
	}
	return false
}
