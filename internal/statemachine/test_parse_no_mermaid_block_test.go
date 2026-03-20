//ff:func feature=statemachine type=parser control=sequence topic=states
//ff:what TestParseNoMermaidBlock: mermaid 블록 없는 입력에서 에러 반환 검증
package statemachine

import (
	"testing"
)

func TestParseNoMermaidBlock(t *testing.T) {
	_, err := Parse("test", "# No mermaid block here")
	if err == nil {
		t.Error("expected error for missing mermaid block")
	}
}
