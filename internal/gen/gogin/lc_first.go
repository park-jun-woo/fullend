//ff:func feature=gen-gogin type=util control=sequence
//ff:what converts to camelCase by lowercasing the first character with Go initialism handling

package gogin

import "github.com/ettle/strcase"

// lcFirst converts to camelCase (lowercases the first character with Go initialism handling).
func lcFirst(s string) string {
	return strcase.ToGoCamel(s)
}
