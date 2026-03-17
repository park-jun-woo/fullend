//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-ddl
//ff:what param 목록에서 대소문자 무시로 매칭되는 파라미터를 찾아 반환
package crosscheck

import "strings"

func findCaseInsensitiveParam(key string, params []string) string {
	for _, p := range params {
		if strings.EqualFold(key, p) {
			return p
		}
	}
	return ""
}
