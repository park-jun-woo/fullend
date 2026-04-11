//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what roleInAnySets — role 값이 하나라도 포함된 StringSet이 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func roleInAnySets(role string, sets []rule.StringSet) bool {
	for _, s := range sets {
		if s[role] {
			return true
		}
	}
	return false
}
