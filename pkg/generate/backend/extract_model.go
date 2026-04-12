//ff:func feature=rule type=util control=sequence
//ff:what extractModel — Model.Method에서 Model 부분 추출
package backend

import "strings"

func extractModel(modelMethod string) string {
	if idx := strings.IndexByte(modelMethod, '.'); idx > 0 {
		return modelMethod[:idx]
	}
	return modelMethod
}
