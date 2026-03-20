//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashServiceFuncDifferent: 다른 입력에 대해 다른 해시를 생성하는지 테스트
package contract

import (
	"testing"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

func TestHashServiceFunc_Different(t *testing.T) {
	sf1 := ssacparser.ServiceFunc{
		Name: "CreateGig",
		Sequences: []ssacparser.Sequence{
			{Type: "post", Args: []ssacparser.Arg{{Source: "request", Field: "Title"}}},
			{Type: "response", Fields: map[string]string{"gig": "gig"}},
		},
	}
	sf2 := ssacparser.ServiceFunc{
		Name: "FindGigByID",
		Sequences: []ssacparser.Sequence{
			{Type: "get", Args: []ssacparser.Arg{{Source: "request", Field: "GigID"}}},
			{Type: "response", Fields: map[string]string{"gig": "gig"}},
		},
	}
	h1 := HashServiceFunc(sf1)
	h2 := HashServiceFunc(sf2)
	if h1 == h2 {
		t.Errorf("different inputs produced same hash: %s", h1)
	}
}
