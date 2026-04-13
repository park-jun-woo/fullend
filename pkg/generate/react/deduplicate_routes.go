//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what 중복 경로를 제거하고 정렬된 라우트 목록을 반환한다

package react

import "sort"

// deduplicateRoutes removes duplicate routes by path (keeps first) and sorts by path.
func deduplicateRoutes(routes []route) []route {
	seen := make(map[string]bool)
	var unique []route
	for _, r := range routes {
		if seen[r.path] {
			continue
		}
		seen[r.path] = true
		unique = append(unique, r)
	}
	sort.Slice(unique, func(i, j int) bool {
		return unique[i].path < unique[j].path
	})
	return unique
}
