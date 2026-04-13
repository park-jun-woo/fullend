//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashServiceFuncDeterministic: 동일 입력에 대해 동일 해시를 생성하는지 테스트
package contract

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func TestHashServiceFunc_Deterministic(t *testing.T) {
	sf := ssacparser.ServiceFunc{
		Name: "CreateGig",
		Sequences: []ssacparser.Sequence{
			{Type: "post", Args: []ssacparser.Arg{{Source: "request", Field: "Title"}, {Source: "request", Field: "Budget"}}},
			{Type: "response", Fields: map[string]string{"gig": "gig"}},
		},
	}
	h1 := HashServiceFunc(sf)
	h2 := HashServiceFunc(sf)
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
	if len(h1) != 7 {
		t.Errorf("hash length = %d, want 7", len(h1))
	}
}
