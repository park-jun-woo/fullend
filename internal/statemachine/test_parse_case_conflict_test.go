//ff:func feature=statemachine type=parser control=sequence topic=states
//ff:what TestParseCaseConflict: 대소문자 충돌 상태명(draft/Draft)에서 에러 반환 검증
package statemachine

import (
	"testing"
)

func TestParseCaseConflict(t *testing.T) {
	content := "```mermaid\nstateDiagram-v2\n    [*] --> draft\n    Draft --> open: PublishGig\n    open --> closed: CloseGig\n```"
	_, err := Parse("test", content)
	if err == nil {
		t.Error("expected error for case-conflicting state names draft/Draft")
	}
}
