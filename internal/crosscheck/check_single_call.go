//ff:func feature=crosscheck type=rule control=sequence topic=func-check
//ff:what 단일 @call 시퀀스의 spec 조회 및 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

// checkSingleCall validates a single @call sequence against its func spec.
func checkSingleCall(
	ctx, key, pkg, camelName string,
	seq ssacparser.Sequence,
	specMap map[string]*funcspec.FuncSpec,
	fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec,
	definedVars map[string]string,
	symbolTable *ssacvalidator.SymbolTable,
	openAPIDoc *openapi3.T,
	sfName string,
) []CrossError {
	spec, found := specMap[key]
	if !found {
		if jwtBuiltinFuncs[key] {
			return nil
		}
		return []CrossError{newMissingFuncError(ctx, key, pkg, camelName, seq)}
	}
	return validateCallSpec(ctx, spec, seq, camelName, fullendPkgSpecs, projectFuncSpecs, definedVars, symbolTable, openAPIDoc, sfName)
}
