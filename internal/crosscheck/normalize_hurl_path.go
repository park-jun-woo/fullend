//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what Hurl URL 경로를 정규화된 세그먼트로 변환
package crosscheck

import (
	"regexp"
	"strings"
)

// normalizeHurlPath converts a Hurl URL path to normalized segments.
// {{variable}} -> ":param", pure numeric literals (e.g. "999999") -> ":param"
func normalizeHurlPath(path string) []string {
	path = strings.TrimSpace(path)
	if idx := strings.Index(path, "?"); idx >= 0 {
		path = path[:idx]
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var segs []string
	reVar := regexp.MustCompile(`\{\{.+?\}\}`)
	reNumeric := regexp.MustCompile(`^\d+$`)
	for _, p := range parts {
		if p == "" {
			continue
		}
		if reVar.MatchString(p) || reNumeric.MatchString(p) {
			segs = append(segs, ":param")
		} else {
			segs = append(segs, p)
		}
	}
	return segs
}
