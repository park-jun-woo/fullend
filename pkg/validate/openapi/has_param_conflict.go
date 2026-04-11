//ff:func feature=rule type=util control=iteration dimension=1
//ff:what hasParamConflict — path에 중복 파라미터가 있는지 확인
package openapi

import "strings"

func hasParamConflict(path string) bool {
	seen := map[string]bool{}
	for _, seg := range strings.Split(path, "/") {
		param := extractParam(seg)
		if param == "" {
			continue
		}
		if seen[param] {
			return true
		}
		seen[param] = true
	}
	return false
}
