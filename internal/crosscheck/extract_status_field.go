//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=states
//ff:what @state Inputs에서 상태 필드명 추출
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// extractStatusField extracts the status field name from @state Inputs.
func extractStatusField(seq ssacparser.Sequence) string {
	for _, v := range seq.Inputs {
		parts := strings.SplitN(v, ".", 2)
		if len(parts) == 2 {
			return parts[1]
		}
		break
	}
	return ""
}
