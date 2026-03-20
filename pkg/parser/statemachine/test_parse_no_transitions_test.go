//ff:func feature=statemachine type=parser control=sequence
//ff:what 전이 없는 다이어그램에서 에러 반환 검증

package statemachine

import (
	"testing"
)

func TestParseNoTransitions(t *testing.T) {
	content := "```mermaid\nstateDiagram-v2\n    [*] --> draft\n```"
	_, err := Parse("test", content)
	if err == nil {
		t.Error("expected error for no transitions")
	}
}
