//ff:func feature=symbol type=util control=sequence topic=ddl
//ff:what VARCHAR(N) 타입에서 길이 N을 추출
package validator

import (
	"regexp"
	"strconv"
)

var reVarcharLen = regexp.MustCompile(`(?i)VARCHAR\((\d+)\)`)

func extractVarcharLen(colType string) int {
	m := reVarcharLen.FindStringSubmatch(colType)
	if len(m) == 2 {
		n, _ := strconv.Atoi(m[1])
		return n
	}
	return 0
}
