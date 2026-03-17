//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what DDL CHECK 제약에서 role 컬럼의 허용 값을 수집
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

// collectDDLRoleValues collects allowed values from DDL CHECK constraints on "role" columns.
func collectDDLRoleValues(st *ssacvalidator.SymbolTable) map[string]bool {
	values := make(map[string]bool)
	for _, table := range st.DDLTables {
		if vals, ok := table.CheckEnums["role"]; ok {
			for _, v := range vals {
				values[v] = true
			}
		}
	}
	return values
}
