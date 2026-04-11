//ff:func feature=rule type=rule control=sequence
//ff:what checkSentinelRecords — FK DEFAULT 0 센티널 레코드 누락 WARNING (D-3)
package ddl

import (
	parseddl "github.com/park-jun-woo/fullend/pkg/parser/ddl"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkSentinelRecords(tables []parseddl.Table) []validate.ValidationError {
	// D-3: FK with DEFAULT 0 needs a sentinel record (id=0) in referenced table.
	// Detection requires DDL INSERT parsing which is not in current Table struct.
	return nil
}
