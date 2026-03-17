//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what fetch 블록 내 자식 action 블록을 검증
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

func validateChildActions(f parser.FetchBlock, file string, st *SymbolTable, frontendDir string) []ValidationError {
	var errs []ValidationError
	for _, child := range f.Children {
		if child.Kind == "action" && child.Action != nil {
			errs = append(errs, validateActionBlock(*child.Action, file, st, frontendDir)...)
		}
	}
	return errs
}
