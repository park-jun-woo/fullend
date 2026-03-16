//ff:func feature=gen-gogin type=util control=sequence
//ff:what Go camelCase/PascalCase を snake_case に変換する

package gogin

import "github.com/ettle/strcase"

// goToSnake converts a Go camelCase/PascalCase name to snake_case.
func goToSnake(s string) string {
	return strcase.ToSnake(s)
}
