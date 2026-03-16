//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-ddl
//ff:what SSaC @result/@param 타입이 DDL과 일치하는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

// CheckSSaCDDL validates SSaC @result types and @param types against DDL.
func CheckSSaCDDL(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable, dtoTypes map[string]bool) []CrossError {
	var errs []CrossError

	for _, fn := range funcs {
		ctx := fmt.Sprintf("%s:%s", fn.FileName, fn.Name)
		errs = append(errs, checkSSaCDDLFunc(fn, st, ctx, dtoTypes)...)
	}

	return errs
}
