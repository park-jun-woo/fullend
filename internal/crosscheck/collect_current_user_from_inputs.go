//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what 단일 시퀀스의 입력에서 currentUser 필드 참조를 수집
package crosscheck

import (
	"strings"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func collectCurrentUserFromInputs(seq ssacparser.Sequence, loc string, result map[string][]string) {
	for _, val := range seq.Inputs {
		if strings.HasPrefix(val, "currentUser.") {
			field := strings.TrimPrefix(val, "currentUser.")
			result[field] = append(result[field], loc)
		}
	}
}
