//ff:func feature=crosscheck type=rule control=sequence
//ff:what SSaC 입력 파라미터가 DDL 테이블 컬럼에 존재하는지 검증
package crosscheck

import (
	"fmt"
	"strings"

	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
)

func checkParamColumn(key, value, modelName, tableName string, table ssacvalidator.DDLTable, ctx string, seqIdx int) []CrossError {
	parts := strings.SplitN(value, ".", 2)
	if parts[0] != "request" {
		return nil
	}

	colName := pascalToSnake(key)

	if strings.EqualFold(key, modelName+"ID") {
		colName = "id"
	}

	if _, ok := table.Columns[colName]; !ok {
		return []CrossError{{
			Rule:       "SSaC arg ↔ DDL",
			Context:    ctx,
			Message:    fmt.Sprintf("seq[%d] input %s (→ %s) not found in table %s", seqIdx, key, colName, tableName),
			Level:      "WARNING",
			Suggestion: fmt.Sprintf("DDL에 추가: ALTER TABLE %s ADD COLUMN %s -- TODO: 타입 지정;", tableName, colName),
		}}
	}

	return nil
}
