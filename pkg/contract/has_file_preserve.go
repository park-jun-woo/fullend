//ff:func feature=contract type=walker control=iteration dimension=1
//ff:what 소스 파일에 파일 수준 preserve 디렉티브가 있는지 확인한다
package contract

import "strings"

// hasFilePreserve checks if source has a file-level //fullend:preserve directive.
func hasFilePreserve(src string) bool {
	lines := strings.SplitN(src, "\n", 10)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "package ") {
			break
		}
		if !strings.HasPrefix(line, "//") {
			continue
		}
		d, err := Parse(line)
		if err != nil {
			continue
		}
		return d.Ownership == "preserve"
	}
	return false
}
