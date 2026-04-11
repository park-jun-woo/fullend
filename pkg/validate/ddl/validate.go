//ff:func feature=rule type=rule control=sequence
//ff:what Validate — DDL 검증: sqlc 중복, NOT NULL, 센티널 (D-1~D-3)
package ddl

import (
	parseddl "github.com/park-jun-woo/fullend/pkg/parser/ddl"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks DDL tables for common issues.
func Validate(tables []parseddl.Table) []validate.ValidationError {
	var errs []validate.ValidationError
	errs = append(errs, checkNullableColumns(tables)...)
	errs = append(errs, checkSentinelRecords(tables)...)
	return errs
}
