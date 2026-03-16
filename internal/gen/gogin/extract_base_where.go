//ff:func feature=gen-gogin type=parser control=iteration dimension=1
//ff:what extracts the WHERE clause from a SQL query string, stripping ORDER BY

package gogin

import (
	"fmt"
	"strings"
)

// extractBaseWhere extracts the WHERE clause from a SQL query string, stripping ORDER BY.
// e.g. "SELECT * FROM courses WHERE published = TRUE ORDER BY created_at DESC;" → "published = TRUE"
// e.g. "SELECT * FROM enrollments WHERE user_id = $1 ORDER BY created_at DESC;" → "user_id = $1"
func extractBaseWhere(sql string) (where string, paramCount int) {
	upper := strings.ToUpper(sql)

	// Find WHERE clause.
	whereIdx := strings.Index(upper, "WHERE ")
	if whereIdx < 0 {
		return "", 0
	}
	rest := sql[whereIdx+6:]

	// Strip ORDER BY and everything after.
	if orderIdx := strings.Index(strings.ToUpper(rest), "ORDER BY"); orderIdx >= 0 {
		rest = rest[:orderIdx]
	}
	// Strip LIMIT.
	if limitIdx := strings.Index(strings.ToUpper(rest), "LIMIT"); limitIdx >= 0 {
		rest = rest[:limitIdx]
	}
	// Strip trailing semicolon and whitespace.
	rest = strings.TrimRight(rest, "; \t\n")

	where = strings.TrimSpace(rest)

	// Count $N params in baseWhere.
	for i := 1; i <= 20; i++ {
		if strings.Contains(where, fmt.Sprintf("$%d", i)) {
			paramCount = i
		}
	}

	return where, paramCount
}
