//ff:func feature=crosscheck type=rule control=sequence
//ff:what RulesCount: 등록된 규칙 수가 기대값과 일치하는지 테스트
package crosscheck

import (
	"testing"
)

func TestRules_Count(t *testing.T) {
	if got := len(Rules()); got != 20 {
		t.Errorf("expected 20 rules, got %d", got)
	}
}
