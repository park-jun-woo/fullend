//ff:func feature=gen-hurl type=util control=sequence
//ff:what Returns sort order for auth operations (register before login).
package hurl

import "strings"

func authOrder(opID string) int {
	if strings.ToLower(opID) == "register" {
		return 0
	}
	return 1
}
