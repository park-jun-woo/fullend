//ff:func feature=gen-gogin type=util control=sequence topic=ddl
//ff:what nobodyPlaceholderFor — UNIQUE 충돌 회피용 auto-seed placeholder 문자열

package db

import "strings"

func nobodyPlaceholderFor(table, col string) string {
	if strings.Contains(col, "email") {
		return "nobody-" + table + "-" + col + "@autoseed.local"
	}
	return "nobody-" + table + "-" + col
}
