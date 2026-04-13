//ff:func feature=manifest type=util control=iteration dimension=1 topic=ddl
//ff:what findTableCaseInsensitive — 테이블 맵에서 이름 대소문자 무시하고 조회

package ddl

import "strings"

func findTableCaseInsensitive(name string, tables map[string]*Table) *Table {
	if t, ok := tables[name]; ok {
		return t
	}
	lower := strings.ToLower(name)
	for k, v := range tables {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return nil
}
