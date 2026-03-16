//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what 필드 목록에서 이름으로 타입 조회
package crosscheck

import "github.com/geul-org/fullend/internal/funcspec"

// findFieldType finds a field's type by name in a list of fields.
func findFieldType(fields []funcspec.Field, name string) string {
	for _, f := range fields {
		if f.Name == name {
			return f.Type
		}
	}
	return ""
}
