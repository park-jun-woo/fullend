//ff:func feature=policy type=util control=iteration dimension=1
//ff:what 쉼표로 구분된 액션 문자열을 파싱하여 슬라이스로 반환한다
package policy

import "strings"

// parseActionSet parses the inside of an action set: "update", "delete", "publish"
func parseActionSet(s string) []string {
	var actions []string
	parts := strings.Split(s, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\"")
		if part != "" {
			actions = append(actions, part)
		}
	}
	return actions
}
