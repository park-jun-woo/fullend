//ff:func feature=manifest type=util control=iteration dimension=1 topic=ddl
//ff:what splitAndTrim — 구분자로 split 후 각 요소 TrimSpace

package ddl

import "strings"

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
