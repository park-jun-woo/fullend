//ff:func feature=orchestrator type=util control=iteration dimension=2
//ff:what counts total columns across all DDL tables

package orchestrator

import ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"

// countDDLColumns counts total columns across all DDL tables.
func countDDLColumns(tables map[string]ssacvalidator.DDLTable) int {
	cols := 0
	for _, t := range tables {
		cols += len(t.Columns)
	}
	return cols
}
