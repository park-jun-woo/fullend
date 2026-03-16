//ff:func feature=genmodel type=util control=sequence
//ff:what 문자열을 Go PascalCase로 변환한다
package genmodel

import "github.com/ettle/strcase"

func toPascalCase(s string) string {
	return strcase.ToGoPascal(s)
}
