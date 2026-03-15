//ff:func feature=gen-gogin type=util
//ff:what snake_case を Go PascalCase に変換する

package gogin

import "github.com/ettle/strcase"

// snakeToGo converts a snake_case column name to a Go PascalCase field name.
func snakeToGo(s string) string {
	return strcase.ToGoPascal(s)
}
