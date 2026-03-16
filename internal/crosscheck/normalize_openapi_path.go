//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=scenario-check
//ff:what OpenAPI 경로를 정규화된 세그먼트로 변환
package crosscheck

import (
	"regexp"
	"strings"
)

// normalizeOpenAPIPath converts an OpenAPI path to normalized segments.
// {param} -> ":param"
func normalizeOpenAPIPath(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var segs []string
	reParam := regexp.MustCompile(`^\{.+\}$`)
	for _, p := range parts {
		if p == "" {
			continue
		}
		if reParam.MatchString(p) {
			segs = append(segs, ":param")
		} else {
			segs = append(segs, p)
		}
	}
	return segs
}
