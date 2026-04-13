//ff:func feature=sqlc-parse type=util control=iteration dimension=1
//ff:what extractWhereSetParams — WHERE/SET 절의 col = $N 패턴을 $N 순서로 반환
package sqlc

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/ettle/strcase"
)

var paramRe = regexp.MustCompile(`(\w+)\s*[=<>!]+\s*\$(\d+)`)

func extractWhereSetParams(sql string) []string {
	matches := paramRe.FindAllStringSubmatch(sql, -1)
	if len(matches) == 0 {
		return nil
	}
	type entry struct {
		pos  int
		name string
	}
	var entries []entry
	for _, m := range matches {
		pos, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		entries = append(entries, entry{pos: pos, name: m[1]})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].pos < entries[j].pos })
	var params []string
	for _, e := range entries {
		params = append(params, strcase.ToGoPascal(e.name))
	}
	return params
}
