//ff:func feature=gen-gogin type=util control=sequence
//ff:what converts to Go PascalCase with Go initialism handling

package gogin

import "github.com/ettle/strcase"

// ucFirst converts to Go PascalCase (uppercases the first character with Go initialism handling).
func ucFirst(s string) string {
	return strcase.ToGoPascal(s)
}
