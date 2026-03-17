//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what 단일 DDL 테이블의 role CHECK enum 값을 values 맵에 추가
package crosscheck

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

func collectRoleValuesFromTable(table ssacvalidator.DDLTable, values map[string]bool) {
	vals, ok := table.CheckEnums["role"]
	if !ok {
		return
	}
	for _, v := range vals {
		values[v] = true
	}
}
