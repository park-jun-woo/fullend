//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-bind 필드가 응답 스키마 또는 custom.ts에 존재하는지 검증
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

// validateFetchBinds checks data-bind fields against response schema and custom.ts.
func validateFetchBinds(binds []parser.FieldBind, opID, file string, api APISymbol, cs *CustomSymbol) []ValidationError {
	var errs []ValidationError
	for _, b := range binds {
		if err := checkBindField(b, opID, file, api, cs); err != nil {
			errs = append(errs, *err)
		}
	}
	return errs
}
