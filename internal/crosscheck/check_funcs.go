//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what SSaC @func 참조를 파싱된 func spec과 교차 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckFuncs validates SSaC @func references against parsed func specs.
func CheckFuncs(
	serviceFuncs []ssacparser.ServiceFunc,
	fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec,
	symbolTable *ssacvalidator.SymbolTable,
	openAPIDoc *openapi3.T,
) []CrossError {
	var errs []CrossError

	specMap := buildFuncSpecMap(fullendPkgSpecs, projectFuncSpecs)

	for _, sf := range serviceFuncs {
		errs = append(errs, checkServiceFuncCalls(sf, specMap, fullendPkgSpecs, projectFuncSpecs, symbolTable, openAPIDoc)...)
	}

	return errs
}
