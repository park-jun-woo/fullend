//ff:func feature=gen-gogin type=util control=sequence
//ff:what formats call arg names as a SQL-style arg string

package gogin

import "strings"

// formatCallArgs formats call arg names as a SQL-style arg string.
func formatCallArgs(callArgNames []string) string {
	if len(callArgNames) == 0 {
		return ""
	}
	return ",\n\t\t" + strings.Join(callArgNames, ", ")
}
