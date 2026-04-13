//ff:func feature=crosscheck type=rule control=iteration dimension=3 topic=ssac-ddl
//ff:what checkDefaultFKSeed — DEFAULT N REFERENCES 가리키는 row 가 seed 에 존재하는지 (X-79)

package crosscheck

import (
	"fmt"
	"strconv"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// checkDefaultFKSeed verifies that a column's DEFAULT numeric value satisfies
// its FK reference — i.e., the referenced table has a seed row with id=<N>.
// WARNING — runtime INSERT will fail FK constraint if seed is missing.
func checkDefaultFKSeed(g *rule.Ground, fs *fullend.Fullstack) []CrossError {
	_ = g
	tables := tableMapByName(fs.DDLTables)
	var errs []CrossError
	for _, t := range fs.DDLTables {
		if len(t.Defaults) == 0 || len(t.ForeignKeys) == 0 {
			continue
		}
		fkMap := foreignKeyMap(t.ForeignKeys)
		for col, defVal := range t.Defaults {
			fk, ok := fkMap[col]
			if !ok {
				continue
			}
			n, err := strconv.ParseInt(defVal, 10, 64)
			if err != nil {
				continue // DEFAULT 가 비수치면 이 규칙 범위 밖
			}
			refT, has := tables[fk.RefTable]
			if !has {
				continue
			}
			if !seedHasID(refT.Seeds, fk.RefColumn, strconv.FormatInt(n, 10)) {
				errs = append(errs, CrossError{
					Rule:       "X-79",
					Context:    fmt.Sprintf("%s.%s DEFAULT %d REFERENCES %s(%s)", t.Name, col, n, fk.RefTable, fk.RefColumn),
					Level:      "WARNING",
					Message:    fmt.Sprintf("DEFAULT %d 이 참조하는 %s(%s=%d) 행이 seed 에 없음 — 런타임 INSERT 시 FK 위반 가능", n, fk.RefTable, fk.RefColumn, n),
					Suggestion: fmt.Sprintf("%s.sql 에 INSERT INTO %s (%s, ...) VALUES (%d, ...) 추가 또는 DEFAULT 제거", fk.RefTable, fk.RefTable, fk.RefColumn, n),
				})
			}
		}
	}
	return errs
}

