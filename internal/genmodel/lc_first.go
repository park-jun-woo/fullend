//ff:func feature=genmodel type=util control=sequence
//ff:what 첫 글자를 소문자로 변환한다
package genmodel

import "github.com/ettle/strcase"

func lcFirst(s string) string {
	return strcase.ToGoCamel(s)
}
