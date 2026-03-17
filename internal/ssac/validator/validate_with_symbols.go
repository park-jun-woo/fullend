//ff:func feature=ssac-validate type=command control=iteration dimension=1
//ff:what 내부 검증 + 외부 SSOT 교차 검증을 수행한다
package validator

import (
	"github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// ValidateWithSymbols는 내부 검증 + 외부 SSOT 교차 검증을 수행한다.
func ValidateWithSymbols(funcs []parser.ServiceFunc, st *SymbolTable) []ValidationError {
	errs := Validate(funcs)
	for _, sf := range funcs {
		errs = append(errs, validateModel(sf, st)...)
		errs = append(errs, validateRequest(sf, st)...)

		errs = append(errs, validateQueryUsage(sf, st)...)
		errs = append(errs, validatePaginationType(sf, st)...)
		errs = append(errs, validateCallInputTypes(sf, st)...)
	}
	errs = append(errs, validateGoReservedWords(funcs, st)...)
	return errs
}
