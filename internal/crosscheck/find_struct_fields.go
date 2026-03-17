//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what struct 목록에서 이름에 맞는 struct의 필드 이름을 추출
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

func findStructFields(structs []ssacparser.StructInfo, typeName string) []string {
	for _, st := range structs {
		if st.Name == typeName {
			return extractFieldNames(st.Fields)
		}
	}
	return nil
}
