//ff:func feature=rule type=util control=sequence
//ff:what extractMethod — Model.Method에서 Method 부분 추출
package backend

import "strings"

func extractMethod(modelMethod string) string {
	if idx := strings.IndexByte(modelMethod, '.'); idx > 0 {
		return modelMethod[idx+1:]
	}
	return ""
}
