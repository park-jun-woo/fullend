//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what OperationID 기준으로 중복 ActionBlock을 제거한다
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// deduplicateActions removes duplicate actions by OperationID.
func deduplicateActions(actions []parser.ActionBlock) []parser.ActionBlock {
	seen := map[string]bool{}
	var result []parser.ActionBlock
	for _, a := range actions {
		if !seen[a.OperationID] {
			seen[a.OperationID] = true
			result = append(result, a)
		}
	}
	return result
}
