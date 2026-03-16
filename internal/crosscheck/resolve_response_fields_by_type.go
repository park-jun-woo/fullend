//ff:func feature=crosscheck type=util control=selection
//ff:what 시퀀스 타입에 따라 응답 필드 해석 (call vs model)
package crosscheck

import (
	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// resolveResponseFieldsByType resolves response fields based on sequence type.
func resolveResponseFieldsByType(seq ssacparser.Sequence, funcSpecs []funcspec.FuncSpec, st *ssacvalidator.SymbolTable) []string {
	switch seq.Type {
	case "call":
		return resolveCallResponseFields(seq, funcSpecs)
	case "get", "put", "post", "delete":
		return resolveDDLResponseFields(seq, st)
	default:
		return nil
	}
}
