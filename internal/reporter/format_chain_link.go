//ff:func feature=reporter type=formatter control=selection
//ff:what 체인 링크 하나를 포맷팅된 문자열로 변환한다
package reporter

import "fmt"

// formatChainLink formats a single chain link for display.
func formatChainLink(link ChainLink, isArtifact bool) string {
	switch isArtifact {
	case true:
		loc := link.File
		if link.Summary != "" && link.Summary != "(file)" {
			loc = link.File + ":" + link.Summary
		}
		return fmt.Sprintf("  %-10s %-45s %s", link.Kind, loc, ownershipIcon(link.Ownership))
	default:
		loc := link.File
		if link.Line > 0 {
			loc = fmt.Sprintf("%s:%d", link.File, link.Line)
		}
		return fmt.Sprintf("  %-10s %-45s %s", link.Kind, loc, link.Summary)
	}
}
