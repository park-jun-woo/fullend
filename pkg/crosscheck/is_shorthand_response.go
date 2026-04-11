//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what isShorthandResponse — 시퀀스에 shorthand @response가 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func isShorthandResponse(seqs []ssac.Sequence) bool {
	for _, seq := range seqs {
		if seq.Type == "response" && len(seq.Fields) == 0 {
			return true
		}
	}
	return false
}
