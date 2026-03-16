//ff:func feature=gen-gogin type=parser control=iteration
//ff:what extracts param count, column names, and cleans up the SQL body

package gogin

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// finishQuery extracts param count, column names, and cleans up the SQL body.
func finishQuery(q *sqlcQuery, sql string, paramRe, insertColRe, updateSetRe *regexp.Regexp) {
	// Store cleaned SQL (trim trailing whitespace).
	q.SQL = strings.TrimSpace(sql)

	// Count parameters.
	matches := paramRe.FindAllString(sql, -1)
	seen := make(map[string]bool)
	for _, m := range matches {
		seen[m] = true
	}
	q.ParamCount = len(seen)

	// Extract column names for INSERT.
	if insMatch := insertColRe.FindStringSubmatch(sql); insMatch != nil {
		cols := strings.Split(insMatch[1], ",")
		for _, c := range cols {
			q.Columns = append(q.Columns, strings.TrimSpace(c))
		}
	}

	// For UPDATE: extract all column = $N patterns, ordered by $N position.
	// This maps interface params to the correct $N placeholder positions.
	if len(q.Columns) == 0 && strings.HasPrefix(strings.TrimSpace(strings.ToUpper(sql)), "UPDATE") {
		updateColNRe := regexp.MustCompile(`(\w+)\s*=\s*\$(\d+)`)
		colMatches := updateColNRe.FindAllStringSubmatch(sql, -1)
		if len(colMatches) > 0 {
			type colPos struct {
				col string
				pos int
			}
			var positions []colPos
			for _, m := range colMatches {
				pos, _ := strconv.Atoi(m[2])
				positions = append(positions, colPos{col: m[1], pos: pos})
			}
			sort.Slice(positions, func(i, j int) bool { return positions[i].pos < positions[j].pos })
			for _, cp := range positions {
				q.Columns = append(q.Columns, cp.col)
			}
		}
	}
}
