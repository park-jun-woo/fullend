//ff:func feature=stml-parse type=parser control=sequence
//ff:what "column:direction" 문자열을 SortDecl로 파싱
package stml

import "strings"

// parseSortDecl parses "column:direction" into a SortDecl.
func parseSortDecl(v string) *SortDecl {
	parts := strings.SplitN(v, ":", 2)
	sd := &SortDecl{Column: strings.TrimSpace(parts[0]), Direction: "asc"}
	if len(parts) == 2 {
		sd.Direction = strings.TrimSpace(parts[1])
	}
	return sd
}
