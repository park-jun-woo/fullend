//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what data-each 필드가 응답에서 배열인지 검증
package validator

import "github.com/park-jun-woo/fullend/internal/stml/parser"

// validateFetchEaches checks data-each fields are arrays in the response.
func validateFetchEaches(eaches []parser.EachBlock, opID, file string, api APISymbol) []ValidationError {
	var errs []ValidationError
	for _, e := range eaches {
		if err := checkEachField(e, opID, file, api); err != nil {
			errs = append(errs, *err)
		}
	}
	return errs
}
