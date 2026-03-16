//ff:func feature=crosscheck type=util control=sequence
//ff:what 상태 필드 PascalCase를 snake_case로 변환
package crosscheck

import "github.com/ettle/strcase"

// pascalToSnakeState converts PascalCase to snake_case.
func pascalToSnakeState(s string) string {
	return strcase.ToSnake(s)
}
