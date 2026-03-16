//ff:func feature=crosscheck type=util control=selection topic=func-check
//ff:what 리터럴 값의 Go 타입 추론
package crosscheck

import "strings"

// inferLiteralType infers the Go type of a literal value.
func inferLiteralType(s string) string {
	switch {
	case s == "true" || s == "false":
		return "bool"
	case s == "nil":
		return ""
	case strings.Contains(s, "."):
		return "float64"
	default:
		return "int"
	}
}
