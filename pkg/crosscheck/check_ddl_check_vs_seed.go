//ff:func feature=crosscheck type=rule control=iteration dimension=3 topic=ssac-ddl
//ff:what checkDDLCheckVsSeed — DDL CHECK IN 목록과 INSERT seed 값 일치 (X-78)

package crosscheck

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkDDLCheckVsSeed verifies that each INSERT seed row's column values satisfy
// the column's CHECK IN constraints (if any). ERROR if violation — psql -f fails.
func checkDDLCheckVsSeed(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	_ = g
	var errs []CrossError
	for _, t := range fs.DDLTables {
		if len(t.CheckEnums) == 0 || len(t.Seeds) == 0 {
			continue
		}
		for i, row := range t.Seeds {
			for col, enumVals := range t.CheckEnums {
				val, present := row[col]
				if !present {
					continue
				}
				if !contains(enumVals, val) {
					errs = append(errs, CrossError{
						Rule:       "X-78",
						Context:    fmt.Sprintf("%s.sql seed[%d].%s=%q", t.Name, i, col, val),
						Level:      "ERROR",
						Message:    fmt.Sprintf("INSERT seed 값 %q 가 CHECK %s 허용 목록 %v 에 없음", val, col, enumVals),
						Suggestion: fmt.Sprintf("INSERT 값을 %v 중 하나로 변경 또는 CHECK 제약 수정", enumVals),
					})
				}
			}
		}
	}
	return errs
}
