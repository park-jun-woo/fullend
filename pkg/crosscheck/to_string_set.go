//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what toStringSet — 문자열 슬라이스를 StringSet으로 변환
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/rule"

func toStringSet(vals []string) rule.StringSet {
	s := make(rule.StringSet, len(vals))
	for _, v := range vals {
		s[v] = true
	}
	return s
}
