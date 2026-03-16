//ff:func feature=symbol type=parser control=iteration dimension=1 topic=ddl
//ff:what CHECK (col IN (...)) 절에서 컬럼명과 허용 값을 파싱
package validator

import (
	"regexp"
	"strings"
)

var reCheckEnum = regexp.MustCompile(`(?i)CHECK\s*\(\s*(\w+)\s+IN\s*\(([^)]+)\)\s*\)`)

func parseCheckEnum(line string) (string, []string) {
	m := reCheckEnum.FindStringSubmatch(line)
	if len(m) < 3 {
		return "", nil
	}
	col := m[1]
	rawVals := m[2]
	var vals []string
	for _, v := range strings.Split(rawVals, ",") {
		v = strings.TrimSpace(v)
		v = strings.Trim(v, "'\"")
		if v != "" {
			vals = append(vals, v)
		}
	}
	return col, vals
}
