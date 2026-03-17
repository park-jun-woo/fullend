//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what 중첩된 data-fetch 블록을 재귀 검증
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func validateNestedFetches(f parser.FetchBlock, file string, st *SymbolTable, cs *CustomSymbol, frontendDir string) []ValidationError {
	var errs []ValidationError
	for _, child := range f.NestedFetches {
		errs = append(errs, validateFetchBlock(child, file, st, cs, frontendDir)...)
	}
	return errs
}
