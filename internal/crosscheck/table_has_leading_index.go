//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 테이블에서 컬럼이 인덱스의 선행 컬럼인지 확인
package crosscheck

import ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"

// tableHasLeadingIndex checks if a column is a leading column in any index of the table.
func tableHasLeadingIndex(col string, table ssacvalidator.DDLTable) bool {
	for _, idx := range table.Indexes {
		if len(idx.Columns) > 0 && idx.Columns[0] == col {
			return true
		}
	}
	return false
}
