//ff:func feature=crosscheck type=rule control=sequence topic=policy-check
//ff:what @ownership 조인 테이블·컬럼이 DDL에 존재하는지 검증
package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/policy"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkJoinTableDDL(om policy.OwnershipMapping, st *ssacvalidator.SymbolTable) []CrossError {
	joinTbl, ok := st.DDLTables[om.JoinTable]
	if !ok {
		return []CrossError{{
			Rule:       "Policy ↔ DDL",
			Context:    fmt.Sprintf("@ownership %s via", om.Resource),
			Message:    fmt.Sprintf("join table %q does not exist in DDL", om.JoinTable),
			Level:      "ERROR",
			Suggestion: fmt.Sprintf("Create table %s in DDL or fix @ownership via annotation", om.JoinTable),
		}}
	}
	if _, colOk := joinTbl.Columns[om.JoinFK]; !colOk {
		return []CrossError{{
			Rule:       "Policy ↔ DDL",
			Context:    fmt.Sprintf("@ownership %s via", om.Resource),
			Message:    fmt.Sprintf("join column %s.%s does not exist in DDL", om.JoinTable, om.JoinFK),
			Level:      "ERROR",
			Suggestion: fmt.Sprintf("Add column %s to table %s in DDL", om.JoinFK, om.JoinTable),
		}}
	}
	return nil
}
