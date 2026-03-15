//ff:func feature=symbol type=util
//ff:what WHERE/SET 절에서 col = $N 패턴을 추출하여 $N 순서대로 반환한다
package validator

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/ettle/strcase"
)

var sqlParamRe = regexp.MustCompile(`(\w+)\s*[=<>!]+\s*\$(\d+)`)

// extractWhereSetParams는 WHERE/SET 절에서 col = $N, col > $N 패턴을 추출하여 $N 순서대로 반환한다.
func extractWhereSetParams(sql string) []string {
	matches := sqlParamRe.FindAllStringSubmatch(sql, -1)
	if len(matches) == 0 {
		return nil
	}

	type paramEntry struct {
		pos  int
		name string
	}
	var entries []paramEntry
	for _, m := range matches {
		pos, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		entries = append(entries, paramEntry{pos: pos, name: m[1]})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].pos < entries[j].pos })

	var params []string
	for _, e := range entries {
		params = append(params, strcase.ToGoPascal(e.name))
	}
	return params
}
