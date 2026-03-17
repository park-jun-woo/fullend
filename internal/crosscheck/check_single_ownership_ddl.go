//ff:func feature=crosscheck type=rule control=sequence topic=policy-check
//ff:what 단일 @ownership 매핑의 테이블·컬럼·조인 테이블이 DDL에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkSingleOwnershipDDL(om policy.OwnershipMapping, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError

	tbl, ok := st.DDLTables[om.Table]
	if !ok {
		return []CrossError{{
			Rule:       "Policy ↔ DDL",
			Context:    fmt.Sprintf("@ownership %s", om.Resource),
			Message:    fmt.Sprintf("ownership table %q does not exist in DDL", om.Table),
			Level:      "ERROR",
			Suggestion: fmt.Sprintf("Create table %s in DDL or fix @ownership annotation", om.Table),
		}}
	}
	if _, colOk := tbl.Columns[om.Column]; !colOk {
		errs = append(errs, CrossError{
			Rule:       "Policy ↔ DDL",
			Context:    fmt.Sprintf("@ownership %s", om.Resource),
			Message:    fmt.Sprintf("ownership column %s.%s does not exist in DDL", om.Table, om.Column),
			Level:      "ERROR",
			Suggestion: fmt.Sprintf("Add column %s to table %s in DDL", om.Column, om.Table),
		})
	}

	if om.JoinTable != "" {
		errs = append(errs, checkJoinTableDDL(om, st)...)
	}

	return errs
}
