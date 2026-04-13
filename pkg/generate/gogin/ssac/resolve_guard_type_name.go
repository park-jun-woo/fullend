//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what 가드 대상의 타입명을 resolver와 resultTypes에서 순차 조회
package ssac

import "strings"

func resolveGuardTypeName(target string, resolver *FieldTypeResolver, resultTypes map[string]string) string {
	typeName := ""
	if resolver != nil && strings.Contains(target, ".") {
		typeName = resolver.ResolveFieldType(target)
	}
	if typeName == "" {
		typeName = resultTypes[rootVar(target)]
	}
	return typeName
}
