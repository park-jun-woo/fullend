//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=func-check
//ff:what 단일 SSaC 서비스 함수의 @call 시퀀스 검증
package crosscheck

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
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
		// Check original function name starts with uppercase (Go exported).
		parts := strings.SplitN(seq.Model, ".", 2)
		rawFuncName := parts[0]
		if len(parts) == 2 {
			rawFuncName = parts[1]
		}
		if len(rawFuncName) > 0 && unicode.IsLower(rune(rawFuncName[0])) {
			errs = append(errs, CrossError{
				Rule:    "SSaC @call naming",
				Context: fmt.Sprintf("%s seq[%d] @call %s", sf.Name, seqIdx, seq.Model),
				Message: fmt.Sprintf("function name %q starts with lowercase — Go exported functions must start with uppercase", rawFuncName),
				Level:   "ERROR",
			})
		}

		pkg, camelName, key := parseCallKey(seq.Model)
		ctx := fmt.Sprintf("%s seq[%d] @call %s", sf.Name, seqIdx, key)
		errs = append(errs, checkSingleCall(ctx, key, pkg, camelName, seq, specMap, fullendPkgSpecs, projectFuncSpecs, definedVars, symbolTable, openAPIDoc, sf.Name)...)
	}

	return errs
}
