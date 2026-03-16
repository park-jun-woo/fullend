//ff:func feature=reporter type=formatter control=selection
//ff:what 소유권 타입에 따른 아이콘 문자열을 반환한다
package reporter

// ownershipIcon returns the display icon for an ownership type.
func ownershipIcon(ownership string) string {
	switch ownership {
	case "preserve":
		return "preserve ✎"
	case "gen":
		return "gen"
	default:
		return ""
	}
}
