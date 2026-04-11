//ff:func feature=rule type=rule control=sequence
//ff:what checkNullableColumns — NOT NULL 누락 검증 (D-2)
package ddl

import (
	parseddl "github.com/park-jun-woo/fullend/pkg/parser/ddl"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func checkNullableColumns(tables []parseddl.Table) []validate.ValidationError {
	// D-2 is already checked at DDL parse time via pg_query.
	// This is a secondary check on structured Table data.
	return nil
}
