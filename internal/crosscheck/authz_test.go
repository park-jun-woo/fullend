package crosscheck

import (
	"testing"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

func TestCheckAuthzValidFields(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "AcceptProposal",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"UserID":     "currentUser.ID",
						"ResourceID": "gig.ClientID",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "")
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestCheckAuthzInvalidField(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "DoSomething",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"UserID":    "currentUser.ID",
						"BadField":  "gig.ClientID",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if errs[0].Level != "ERROR" {
		t.Fatalf("expected ERROR level, got %s", errs[0].Level)
	}
}

func TestCheckAuthzCustomPackageSkips(t *testing.T) {
	funcs := []ssacparser.ServiceFunc{
		{
			Name: "DoSomething",
			Sequences: []ssacparser.Sequence{
				{
					Type: "auth",
					Inputs: map[string]string{
						"CustomField": "value",
					},
				},
			},
		},
	}

	errs := CheckAuthz(funcs, "github.com/custom/authz")
	if len(errs) != 0 {
		t.Fatalf("expected no errors with custom package, got %d", len(errs))
	}
}
