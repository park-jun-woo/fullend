//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what go.mod 파일이 fullend 모듈인지 확인한다

package orchestrator

import "strings"

// isFullendGoMod checks if go.mod data declares the fullend module.
func isFullendGoMod(data []byte) bool {
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "module github.com/park-jun-woo/fullend" {
			return true
		}
	}
	return false
}
