//ff:func feature=orchestrator type=rule control=iteration dimension=2
//ff:what 경로 파라미터 이름 충돌 감지 — 같은 세그먼트 위치에서 다른 파라미터 이름 검출
package orchestrator

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// checkPathParamConflicts detects path param name conflicts at the same segment position.
// e.g. /gigs/{ID} and /gigs/{GigID}/proposals conflict because segment[1] has both {ID} and {GigID}.
func checkPathParamConflicts(doc *openapi3.T) []string {
	if doc == nil || doc.Paths == nil {
		return nil
	}

	// Group: "prefix" → map[paramName][]fullPath
	// prefix is the path up to but not including the param segment, plus position index.
	// e.g. "/gigs/{ID}" → prefix="/gigs/", position=1
	type segKey struct {
		prefix   string
		position int
	}
	paramAt := make(map[segKey]map[string][]string) // segKey → paramName → []paths

	for path := range doc.Paths.Map() {
		segments := strings.Split(strings.Trim(path, "/"), "/")
		for i, seg := range segments {
			if !strings.HasPrefix(seg, "{") || !strings.HasSuffix(seg, "}") {
				continue
			}
			paramName := seg[1 : len(seg)-1]
			key := segKey{prefix: strings.Join(segments[:i], "/"), position: i}
			if paramAt[key] == nil {
				paramAt[key] = make(map[string][]string)
			}
			paramAt[key][paramName] = append(paramAt[key][paramName], path)
		}
	}

	var errs []string
	for key, names := range paramAt {
		if len(names) <= 1 {
			continue
		}
		var nameList []string
		for n := range names {
			nameList = append(nameList, "{"+n+"}")
		}
		errs = append(errs, fmt.Sprintf(
			"path param 충돌: segment[%d] (prefix=/%s/)에 %s가 혼재 — 이름을 통일하세요",
			key.position, key.prefix, strings.Join(nameList, ", "),
		))
	}
	return errs
}
