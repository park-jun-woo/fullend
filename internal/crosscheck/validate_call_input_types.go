//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what @call Input 키 이름과 타입이 Request 필드와 일치하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// validateCallInputTypes checks that input key names and types match request fields.
func validateCallInputTypes(
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
	if len(seq.Inputs) == 0 {
		return nil
	}

	var errs []CrossError
	reqFieldMap := make(map[string]string)
	for _, rf := range spec.RequestFields {
		reqFieldMap[rf.Name] = rf.Type
	}

	for inputKey, inputValue := range seq.Inputs {
		reqType, exists := reqFieldMap[inputKey]
		if !exists {
			errs = append(errs, CrossError{
				Rule:    "Func ↔ SSaC",
				Context: ctx,
				Message: fmt.Sprintf("@call Input 필드 %q가 %sRequest에 없음", inputKey, strcase.ToGoPascal(funcName)),
				Level:   "ERROR",
			})
			continue
		}
		allSpecs := append(fullendPkgSpecs, projectFuncSpecs...)
		valueType := resolveInputValueType(inputValue, definedVars, symbolTable, openAPIDoc, sfName, allSpecs)
		if valueType != "" && !typesCompatible(valueType, reqType) {
			errs = append(errs, CrossError{
				Rule:    "Func ↔ SSaC",
				Context: ctx,
				Message: fmt.Sprintf("@call Input %s 타입 불일치: %s(source) ≠ %s(Request)", inputKey, valueType, reqType),
				Level:   "ERROR",
			})
		}
	}
	return errs
}
