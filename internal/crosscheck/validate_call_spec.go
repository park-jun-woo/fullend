//ff:func feature=crosscheck type=rule control=sequence topic=func-check
//ff:what @call 시퀀스의 spec 매칭 (HasBody, 금지 import, 파라미터, 결과, 소스 변수) 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// validateCallSpec validates a single @call sequence against its func spec.
func validateCallSpec(
	ctx string,
	spec *funcspec.FuncSpec,
	seq ssacparser.Sequence,
	funcName string,
	fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec,
	definedVars map[string]string,
	symbolTable *ssacvalidator.SymbolTable,
	openAPIDoc *openapi3.T,
	sfName string,
) []CrossError {
	var errs []CrossError

	errs = append(errs, validateCallBody(ctx, spec)...)
	errs = append(errs, validateCallImports(ctx, spec)...)
	errs = append(errs, validateCallInputCount(ctx, spec, seq)...)
	errs = append(errs, validateCallInputTypes(ctx, spec, seq, funcName, fullendPkgSpecs, projectFuncSpecs, definedVars, symbolTable, openAPIDoc, sfName)...)
	errs = append(errs, validateCallResult(ctx, spec, seq)...)
	errs = append(errs, validateCallSourceVars(ctx, seq, definedVars)...)

	return errs
}
