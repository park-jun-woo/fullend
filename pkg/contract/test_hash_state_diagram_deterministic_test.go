//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashStateDiagramDeterministic: 동일 상태 다이어그램에 대해 동일 해시를 생성하는지 테스트
package contract

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
)

func TestHashStateDiagram_Deterministic(t *testing.T) {
	sd := &statemachine.StateDiagram{
		ID:     "gig",
		States: []string{"draft", "published", "cancelled"},
		Transitions: []statemachine.Transition{
			{From: "draft", To: "published", Event: "PublishGig"},
			{From: "published", To: "cancelled", Event: "CancelGig"},
		},
	}
	h1 := HashStateDiagram(sd)
	h2 := HashStateDiagram(sd)
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
}
