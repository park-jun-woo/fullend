//ff:func feature=crosscheck type=util control=sequence
//ff:what snake_case를 Go PascalCase로 변환
package crosscheck

import "github.com/ettle/strcase"

// snakeToPascal converts snake_case to PascalCase with Go acronym handling.
func snakeToPascal(s string) string {
	return strcase.ToGoPascal(s)
}
