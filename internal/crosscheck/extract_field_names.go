//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what struct 필드 목록에서 필드 이름만 추출
package crosscheck

import ssacparser "github.com/geul-org/fullend/internal/ssac/parser"

func extractFieldNames(fields []ssacparser.StructField) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names
}
