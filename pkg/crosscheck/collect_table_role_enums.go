//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectTableRoleEnums — 테이블의 CHECK enum에서 role 컬럼의 값 집합 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func collectTableRoleEnums(checkEnums map[string][]string) []rule.StringSet {
	var sets []rule.StringSet
	for col, vals := range checkEnums {
		if col == "role" || col == "user_role" {
			sets = append(sets, toStringSet(vals))
		}
	}
	return sets
}
