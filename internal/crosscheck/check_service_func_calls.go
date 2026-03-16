//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what 단일 SSaC 서비스 함수의 @call 시퀀스 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// checkServiceFuncCalls validates all @call sequences in a single service function.
func checkServiceFuncCalls(
	sf ssacparser.ServiceFunc,
	specMap map[string]*funcspec.FuncSpec,
	fullendPkgSpecs, projectFuncSpecs []funcspec.FuncSpec,
	symbolTable *ssacvalidator.SymbolTable,
	openAPIDoc *openapi3.T,
) []CrossError {
	var errs []CrossError
	definedVars := make(map[string]string)

	for seqIdx, seq := range sf.Sequences {
		if seq.Result != nil {
			definedVars[seq.Result.Var] = seq.Result.Type
		}
		if seq.Type != "call" || seq.Model == "" {
			continue
		}
		pkg, camelName, key := parseCallKey(seq.Model)
		ctx := fmt.Sprintf("%s seq[%d] @call %s", sf.Name, seqIdx, key)
		errs = append(errs, checkSingleCall(ctx, key, pkg, camelName, seq, specMap, fullendPkgSpecs, projectFuncSpecs, definedVars, symbolTable, openAPIDoc, sf.Name)...)
	}

	return errs
}
