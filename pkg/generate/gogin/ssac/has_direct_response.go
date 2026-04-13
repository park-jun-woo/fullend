//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=response
//ff:what 시퀀스에 직접 응답(@response varName) 패턴이 있는지 확인
package ssac

import ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func hasDirectResponse(seqs []ssacparser.Sequence) bool {
	for _, s := range seqs {
		if s.Type == ssacparser.SeqResponse && s.Target != "" {
			return true
		}
	}
	return false
}
