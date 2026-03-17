//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=response
//ff:what 시퀀스에 직접 응답(@response varName) 패턴이 있는지 확인
package generator

import "github.com/park-jun-woo/fullend/internal/ssac/parser"

func hasDirectResponse(seqs []parser.Sequence) bool {
	for _, s := range seqs {
		if s.Type == parser.SeqResponse && s.Target != "" {
			return true
		}
	}
	return false
}
