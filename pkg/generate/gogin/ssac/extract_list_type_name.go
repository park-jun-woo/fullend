//ff:func feature=ssac-gen type=util control=sequence topic=string-convert
//ff:what Result 타입에서 리스트 요소 타입명을 추출
package ssac

import "strings"

func extractListTypeName(usage modelUsage) string {
	typeName := "interface{}"
	if usage.Result != nil {
		typeName = usage.Result.Type
		if strings.HasPrefix(typeName, "[]") {
			typeName = typeName[2:]
		}
	}
	return typeName
}
