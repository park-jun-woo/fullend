//ff:func feature=crosscheck type=rule control=iteration dimension=2 topic=openapi-ddl
//ff:what OpenAPI 필드별 제약을 DDL VARCHAR/CHECK와 대조하여 누락·불일치 검출
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

func checkFieldConstraints(st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError

	for opID, rs := range st.RequestSchemas {
		for fieldName, fc := range rs.Fields {
			col := toSnakeCase(fieldName)
			varcharLen, checkEnums, found := findDDLColumnConstraints(st, col)
			errs = append(errs, checkSingleFieldConstraint(opID, fieldName, col, fc, varcharLen, checkEnums, found)...)
		}
	}
	return errs
}
