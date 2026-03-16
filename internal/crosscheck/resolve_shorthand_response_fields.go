//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what shorthand @response 변수명에서 JSON 필드명 해석
package crosscheck

import (
	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// resolveShorthandResponseFields resolves the JSON field names for a shorthand @response varName.
func resolveShorthandResponseFields(
	fn ssacparser.ServiceFunc,
	funcSpecs []funcspec.FuncSpec,
	st *ssacvalidator.SymbolTable,
) []string {
	var varName string
	for _, seq := range fn.Sequences {
		if seq.Type == "response" && seq.Target != "" {
			varName = seq.Target
			break
		}
	}
	if varName == "" {
		return nil
	}

	for _, seq := range fn.Sequences {
		if seq.Result == nil || seq.Result.Var != varName {
			continue
		}
		if seq.Result.Wrapper != "" {
			return nil
		}
		return resolveResponseFieldsByType(seq, funcSpecs, st)
	}

	return nil
}
