//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=query-opts
//ff:what ServiceFunc의 시퀀스에서 query 옵션이 필요한지 확인
package generator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
	"github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func needsQueryOpts(sf parser.ServiceFunc, st *validator.SymbolTable) bool {
	for _, seq := range sf.Sequences {
		if hasQueryInput(seq.Inputs) {
			return true
		}
	}
	return false
}
