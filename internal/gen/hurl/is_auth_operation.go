//ff:func feature=gen-hurl type=util
//ff:what Checks if an operation ID is an auth operation (register or login).
package hurl

import "strings"

func isAuthOperation(opID string) bool {
	lower := strings.ToLower(opID)
	return lower == "register" || lower == "login"
}
