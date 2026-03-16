//ff:func feature=orchestrator type=util control=iteration
//ff:what sortedStringKeys returns sorted keys from a map[string]bool.

package orchestrator

import "sort"

func sortedStringKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
