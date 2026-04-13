//ff:func feature=contract type=walker control=iteration dimension=1
//ff:what 소스 파일 상단에서 파일 수준 fullend 디렉티브를 찾는다
package contract

import "strings"

// findFileLevelDirective finds a file-level //fullend: directive in source.
func findFileLevelDirective(src string) *Directive {
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
		return d
	}
	return nil
}
