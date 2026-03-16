//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 파생 파라미터 배열을 Go 메서드 시그니처 문자열로 변환
package generator

import "strings"

func renderParams(params []derivedParam) string {
	var parts []string
	for _, p := range params {
		if p.Name == "" {
			continue
		}
		parts = append(parts, p.Name+" "+p.GoType)
	}
	return strings.Join(parts, ", ")
}
