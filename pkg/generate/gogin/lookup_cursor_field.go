//ff:func feature=gen-gogin type=util control=iteration dimension=1 topic=query-opts
//ff:what finds cursor field name by matching method name to cursor specs

package gogin

import "strings"

// lookupCursorField finds the cursor field name for a method from cursorSpecs.
func lookupCursorField(cursorSpecs map[string]string, methodName string) string {
	for opID, field := range cursorSpecs {
		if opID == methodName || strings.EqualFold(opID, methodName) {
			return field
		}
	}
	return "ID"
}
