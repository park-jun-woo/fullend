//ff:func feature=ssac-validate type=util control=sequence topic=type-resolve
//ff:what Model.Method 문자열에서 Model 부분 추출
package validator

import "strings"

func extractModel(modelMethod string) string {
	if idx := strings.IndexByte(modelMethod, '.'); idx > 0 {
		return modelMethod[:idx]
	}
	return ""
}
