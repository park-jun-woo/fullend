//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what map 키를 콤마로 합쳐 문자열로 반환
package crosscheck

import "strings"

// joinKeys returns sorted comma-joined keys of a map.
func joinKeys(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
