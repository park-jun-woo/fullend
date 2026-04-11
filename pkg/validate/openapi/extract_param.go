//ff:func feature=rule type=util control=sequence
//ff:what extractParam — path segment에서 {param} 추출, 없으면 빈 문자열
package openapi

import "strings"

func extractParam(seg string) string {
	if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
		return seg[1 : len(seg)-1]
	}
	return ""
}
