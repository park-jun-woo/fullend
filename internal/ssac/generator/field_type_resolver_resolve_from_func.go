//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=type-resolve
//ff:what FuncSpec에서 함수명+필드명으로 응답 필드 타입을 조회
package generator

import "strings"

func (r *FieldTypeResolver) resolveFromFunc(modelName, fieldName string) string {
	for _, spec := range r.fs {
		if strings.EqualFold(spec.Name, modelName) {
			return findResponseFieldType(spec.ResponseFields, fieldName)
		}
	}
	return ""
}
