package contract

import (
	"testing"

	ssacparser "github.com/geul-org/ssac/parser"

	"github.com/geul-org/fullend/internal/statemachine"
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

func TestHashModelMethod_Deterministic(t *testing.T) {
	h1 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	h2 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
}

func TestHashModelMethod_Different(t *testing.T) {
	h1 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	h2 := HashModelMethod("FindByID", []string{"int64"}, []string{"*Gig", "error"})
	if h1 == h2 {
		t.Errorf("different methods produced same hash: %s", h1)
	}
}

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

func TestHashClaims(t *testing.T) {
	claims := map[string]string{"user_id": "int64", "email": "string", "role": "string"}
	h1 := HashClaims(claims)
	h2 := HashClaims(claims)
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
	if len(h1) != 7 {
		t.Errorf("hash length = %d, want 7", len(h1))
	}
}
