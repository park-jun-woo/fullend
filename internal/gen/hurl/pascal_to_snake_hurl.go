//ff:func feature=gen-hurl type=util control=sequence
//ff:what Converts PascalCase to snake_case for Hurl variable names.
package hurl

import "github.com/ettle/strcase"

// pascalToSnakeHurl converts PascalCase to snake_case for Hurl variable names.
func pascalToSnakeHurl(s string) string {
	return strcase.ToSnake(s)
}
