//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=request-params
//ff:what POST/PUT 또는 2+ 파라미터일 때 JSON body 사용 여부를 판단
package generator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

func shouldUseJSONBody(seqs []parser.Sequence, st *validator.SymbolTable, rawParams []rawParam) bool {
	hasBodySeq := false
	for _, seq := range seqs {
		if seq.Type == parser.SeqPost || seq.Type == parser.SeqPut {
			hasBodySeq = true
			break
		}
	}
	return (st != nil && len(rawParams) >= 2) || (hasBodySeq && len(rawParams) >= 1)
}
