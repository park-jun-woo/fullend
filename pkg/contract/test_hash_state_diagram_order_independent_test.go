//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashStateDiagramOrderIndependent: 순서가 달라도 동일 해시를 생성하는지 테스트
package contract

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
)

func TestHashStateDiagram_OrderIndependent(t *testing.T) {
	sd1 := &statemachine.StateDiagram{
		States: []string{"draft", "published"},
		Transitions: []statemachine.Transition{
			{From: "draft", To: "published", Event: "Publish"},
			{From: "published", To: "draft", Event: "Unpublish"},
		},
	}
	sd2 := &statemachine.StateDiagram{
		States: []string{"published", "draft"},
		Transitions: []statemachine.Transition{
			{From: "published", To: "draft", Event: "Unpublish"},
			{From: "draft", To: "published", Event: "Publish"},
		},
	}
	h1 := HashStateDiagram(sd1)
	h2 := HashStateDiagram(sd2)
	if h1 != h2 {
		t.Errorf("reordered input produced different hashes: %s vs %s", h1, h2)
	}
}
