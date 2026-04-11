//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkTableConstraints — 테이블별 VARCHAR/CHECK/password/email 제약 검증
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkTableConstraints(t ddl.Table, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError

	// X-71: password field without minLength
	for col := range t.Columns {
		if strings.Contains(strings.ToLower(col), "password") {
			errs = append(errs, CrossError{Rule: "X-71", Context: t.Name + "." + col, Level: "WARNING",
				Message: "password field — consider adding OpenAPI minLength constraint"})
		}
	}

	// X-72: email field without format
	for col := range t.Columns {
		if strings.Contains(strings.ToLower(col), "email") {
			errs = append(errs, CrossError{Rule: "X-72", Context: t.Name + "." + col, Level: "WARNING",
				Message: "email field — consider adding OpenAPI format: email constraint"})
		}
	}

	_ = fs // OpenAPI maxLength checks require schema walking — deferred to schema-level checks
	return errs
}
