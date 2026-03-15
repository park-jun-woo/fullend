//ff:func feature=gen-gogin type=util
//ff:what returns true if the method name indicates a list query

package gogin

import "strings"

// isListMethod returns true if the method name indicates a list query.
func isListMethod(name string) bool {
	return strings.HasPrefix(name, "List")
}
