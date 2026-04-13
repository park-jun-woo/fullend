//ff:func feature=rule type=loader control=sequence
//ff:what populateRequestSchemas 검증 — RequestConstraints → Ground.ReqSchemas
package ground

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/fullend"
	oapiparser "github.com/park-jun-woo/fullend/pkg/parser/openapi"
)

func TestPopulateRequestSchemas(t *testing.T) {
	g := newGround()
	minL := 1
	maxL := 100
	fs := &fullend.Fullstack{
		RequestConstraints: map[string]map[string]oapiparser.FieldConstraint{
			"CreateUser": {
				"email":    {Type: "string", Format: "email", Required: true, MinLength: &minL, MaxLength: &maxL},
				"password": {Type: "string", Required: true, MinLength: &minL},
				"role":     {Type: "string", Enum: []string{"admin", "user"}},
			},
		},
	}
	populateRequestSchemas(g, fs)

	rs, ok := g.ReqSchemas["CreateUser"]
	if !ok {
		t.Fatal("CreateUser schema not populated")
	}
	email, ok := rs.Fields["email"]
	if !ok {
		t.Fatal("email field missing")
	}
	if !email.Required || email.Format != "email" || email.MinLength == nil || *email.MinLength != 1 || email.MaxLength == nil || *email.MaxLength != 100 {
		t.Errorf("email constraint: %+v", email)
	}
	role, ok := rs.Fields["role"]
	if !ok || len(role.Enum) != 2 {
		t.Errorf("role enum: %+v", role)
	}
}
