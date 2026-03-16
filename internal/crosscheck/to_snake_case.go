//ff:func feature=crosscheck type=util control=sequence
//ff:what camelCase/PascalCase를 snake_case로 변환
package crosscheck

import "github.com/ettle/strcase"

// toSnakeCase converts camelCase/PascalCase to snake_case.
func toSnakeCase(s string) string {
	return strcase.ToSnake(s)
}
