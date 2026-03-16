//ff:func feature=orchestrator type=util control=sequence
//ff:what returns true for lines that are not column definitions

package orchestrator

import "strings"

// isSkippableDDLLine returns true for lines that are not column definitions.
func isSkippableDDLLine(trimmed string) bool {
	if trimmed == "" || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, ")") {
		return true
	}
	upper := strings.ToUpper(trimmed)
	if strings.HasPrefix(upper, "CREATE") || strings.HasPrefix(upper, "INSERT") || strings.HasPrefix(upper, "ON ") || strings.HasPrefix(upper, "VALUES") {
		return true
	}
	if strings.HasPrefix(upper, "PRIMARY KEY") || strings.HasPrefix(upper, "UNIQUE") || strings.HasPrefix(upper, "CHECK") || strings.HasPrefix(upper, "FOREIGN KEY") || strings.HasPrefix(upper, "CONSTRAINT") {
		return true
	}
	return false
}
