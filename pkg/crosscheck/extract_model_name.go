//ff:func feature=crosscheck type=util control=sequence
//ff:what extractModelName — Model.Method에서 Model 부분 추출
package crosscheck

import "strings"

func extractModelName(model string) string {
	if idx := strings.IndexByte(model, '.'); idx > 0 {
		return model[:idx]
	}
	return ""
}
