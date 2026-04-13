//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectDDLRoleCheckSets — DDL CHECK 제약에서 role 컬럼의 enum 값 집합 수집
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func collectDDLRoleCheckSets(fs *fullend.Fullstack) []rule.StringSet {
	var sets []rule.StringSet
	for _, t := range fs.DDLTables {
		sets = append(sets, collectTableRoleEnums(t.CheckEnums)...)
	}
	return sets
}
