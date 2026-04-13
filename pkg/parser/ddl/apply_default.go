//ff:func feature=manifest type=util control=sequence topic=ddl
//ff:what applyDefault — 컬럼 정의 라인에서 DEFAULT 값 추출

package ddl

import "regexp"

var reDefault = regexp.MustCompile(`(?i)DEFAULT\s+('[^']*'|[^\s,)]+)`)

func applyDefault(line, colName string, t *Table) {
	m := reDefault.FindStringSubmatch(line)
	if m == nil {
		return
	}
	if t.Defaults == nil {
		t.Defaults = make(map[string]string)
	}
	t.Defaults[colName] = stripSQLQuotes(m[1])
}
