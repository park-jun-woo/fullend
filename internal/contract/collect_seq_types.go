//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 시퀀스 목록에서 타입 문자열 슬라이스를 수집한다
package contract

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

// collectSeqTypes returns "@type" strings for each sequence.
func collectSeqTypes(seqs []ssacparser.Sequence) []string {
	var types []string
	for _, seq := range seqs {
		types = append(types, "@"+seq.Type)
	}
	return types
}
