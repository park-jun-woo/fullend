//ff:func feature=crosscheck type=util control=sequence
//ff:what PascalCase를 snake_case로 변환
package crosscheck

import "github.com/ettle/strcase"

// pascalToSnake converts PascalCase to snake_case.
func pascalToSnake(s string) string {
	return strcase.ToSnake(s)
}
