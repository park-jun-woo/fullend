//ff:func feature=ssac-validate type=util control=iteration dimension=2 topic=type-resolve
//ff:what struct에 지정된 필드가 존재하는지 확인한다
package validator

import "github.com/geul-org/fullend/internal/ssac/parser"

// hasStructField는 struct에 지정된 필드가 존재하는지 확인한다.
func hasStructField(structs []parser.StructInfo, typeName, fieldName string) bool {
	for _, si := range structs {
		if si.Name != typeName {
			continue
		}
		for _, f := range si.Fields {
			if f.Name == fieldName {
				return true
			}
		}
		return false
	}
	return false
}
