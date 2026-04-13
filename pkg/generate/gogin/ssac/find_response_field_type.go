//ff:func feature=ssac-gen type=util control=iteration dimension=1 topic=response
//ff:what 응답 필드 배열에서 이름으로 타입을 검색
package ssac

import "github.com/park-jun-woo/fullend/internal/funcspec"

func findResponseFieldType(fields []funcspec.Field, fieldName string) string {
	for _, f := range fields {
		if f.Name == fieldName {
			return f.Type
		}
	}
	return ""
}
