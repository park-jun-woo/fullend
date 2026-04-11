//ff:func feature=rule type=util control=iteration dimension=1
//ff:what snakeToPascal — snake_case를 PascalCase로 변환
package ground

import "strings"

func snakeToPascal(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if len(p) > 0 {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}
