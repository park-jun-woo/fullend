//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what compareDDLEnumWithOpenAPICols — 테이블의 각 CHECK enum 컬럼을 OpenAPI enum과 비교
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func compareDDLEnumWithOpenAPICols(table string, checkEnums map[string][]string, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for col, ddlVals := range checkEnums {
		errs = append(errs, compareDDLEnumWithOpenAPI(table, col, ddlVals, fs)...)
	}
	return errs
}
