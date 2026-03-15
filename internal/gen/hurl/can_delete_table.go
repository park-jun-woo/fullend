//ff:func feature=gen-hurl type=util
//ff:what Checks if a table can be safely deleted (all FK children also have DELETE endpoints).
package hurl

// canDeleteTable returns true if a table can be safely deleted in smoke tests.
// A table is deletable only if all child tables (that reference it via FK) also have DELETE endpoints.
func canDeleteTable(table string, deletableTables map[string]bool, reverseDeps map[string][]string) bool {
	children := reverseDeps[table]
	for _, child := range children {
		if !deletableTables[child] {
			return false
		}
		// Recursively check: the child must also be deletable.
		if !canDeleteTable(child, deletableTables, reverseDeps) {
			return false
		}
	}
	return true
}
