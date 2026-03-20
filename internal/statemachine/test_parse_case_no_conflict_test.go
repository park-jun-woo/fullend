//ff:func feature=statemachine type=parser control=sequence topic=states
//ff:what TestParseCaseNoConflict: 대소문자 일관된 상태명에서 정상 파싱 검증
package statemachine

import (
	"testing"
)

func TestParseCaseNoConflict(t *testing.T) {
	content := "```mermaid\nstateDiagram-v2\n    [*] --> draft\n    draft --> open: PublishGig\n    open --> closed: CloseGig\n```"
	_, err := Parse("test", content)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
