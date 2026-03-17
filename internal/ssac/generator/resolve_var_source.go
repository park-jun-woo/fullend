//ff:func feature=ssac-gen type=util control=selection topic=type-resolve
//ff:what 시퀀스 타입에 따라 변수 출처(DDL/func)를 결정
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func resolveVarSource(seq parser.Sequence) varSource {
	switch seq.Type {
	case parser.SeqCall:
		parts := strings.SplitN(seq.Model, ".", 2)
		if len(parts) == 2 {
			return varSource{Kind: "func", ModelName: parts[1]}
		}
		return varSource{Kind: "func", ModelName: seq.Model}
	default:
		return varSource{Kind: "ddl", ModelName: seq.Result.Type}
	}
}
