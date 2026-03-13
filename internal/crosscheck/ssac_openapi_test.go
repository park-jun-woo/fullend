package crosscheck

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	ssacparser "github.com/geul-org/ssac/parser"
	ssacvalidator "github.com/geul-org/ssac/validator"
)

func TestCheckErrStatus_DefaultDefined(t *testing.T) {
	// @empty with default 404, OpenAPI has 404 response → no error.
	doc := buildErrStatusDoc("GetGig", "404")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "gig.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:   "empty",
			Target: "gig",
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}

func TestCheckErrStatus_DefaultMissing(t *testing.T) {
	// @empty with default 404, OpenAPI has no 404 response → error.
	doc := buildErrStatusDoc("GetGig", "200")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetGig",
		FileName: "gig.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:   "empty",
			Target: "gig",
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !contains(errs[0].Message, "404") {
		t.Errorf("expected 404 in message, got: %s", errs[0].Message)
	}
}

func TestCheckErrStatus_CustomDefined(t *testing.T) {
	// @empty with custom 402, OpenAPI has 402 response → no error.
	doc := buildErrStatusDoc("ExecuteWorkflow", "402")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "ExecuteWorkflow",
		FileName: "workflow.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:      "empty",
			Target:    "org",
			ErrStatus: 402,
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors, got %d: %+v", len(errs), errs)
	}
}

func TestCheckErrStatus_CustomMissing(t *testing.T) {
	// @empty with custom 402, OpenAPI has no 402 response → error.
	doc := buildErrStatusDoc("ExecuteWorkflow", "404")

	funcs := []ssacparser.ServiceFunc{{
		Name:     "ExecuteWorkflow",
		FileName: "workflow.ssac",
		Sequences: []ssacparser.Sequence{{
			Type:      "empty",
			Target:    "org",
			ErrStatus: 402,
		}},
	}}

	errs := checkErrStatus(funcs, doc)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %+v", len(errs), errs)
	}
	if !contains(errs[0].Message, "402") {
		t.Errorf("expected 402 in message, got: %s", errs[0].Message)
	}
}

// --- Shorthand @response tests ---

func TestCheckResponseFields_ShorthandCallMismatch(t *testing.T) {
	// @response token + funcspec AccessToken(json:access_token) + OpenAPI AccessToken → ERROR
	doc := buildResponseDoc("Login", map[string]string{"AccessToken": "string"})
	funcSpecs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "issueToken",
		ResponseFields: []funcspec.Field{
			{Name: "AccessToken", Type: "string", JSONName: "access_token"},
		},
	}}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "Login",
		FileName: "login.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "call", Result: &ssacparser.Result{Type: "auth.IssueTokenResponse", Var: "token"}},
			{Type: "response", Target: "token"},
		},
	}}

	errs := checkResponseFields(funcs, nil, doc, funcSpecs)
	// Should detect: "access_token" not in OpenAPI {AccessToken}
	foundErr := false
	for _, e := range errs {
		if e.Level != "WARNING" && contains(e.Message, "access_token") {
			foundErr = true
		}
	}
	if !foundErr {
		t.Errorf("expected ERROR for access_token mismatch, got: %+v", errs)
	}
}

func TestCheckResponseFields_ShorthandCallMatch(t *testing.T) {
	// @response token + funcspec AccessToken(json:access_token) + OpenAPI access_token → pass
	doc := buildResponseDoc("Login", map[string]string{"access_token": "string"})
	funcSpecs := []funcspec.FuncSpec{{
		Package: "auth",
		Name:    "issueToken",
		ResponseFields: []funcspec.Field{
			{Name: "AccessToken", Type: "string", JSONName: "access_token"},
		},
	}}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "Login",
		FileName: "login.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "call", Result: &ssacparser.Result{Type: "auth.IssueTokenResponse", Var: "token"}},
			{Type: "response", Target: "token"},
		},
	}}

	errs := checkResponseFields(funcs, nil, doc, funcSpecs)
	for _, e := range errs {
		if e.Level != "WARNING" {
			t.Errorf("unexpected ERROR: %+v", e)
		}
	}
}

func TestCheckResponseFields_ShorthandDDLMatch(t *testing.T) {
	// @response user + DDL columns [id, email, name] + OpenAPI [id, email, name] → pass
	doc := buildResponseDoc("GetUser", map[string]string{"id": "integer", "email": "string", "name": "string"})
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"users": {Columns: map[string]string{"id": "int64", "email": "string", "name": "string"}},
		},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "GetUser",
		FileName: "get_user.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "get", Result: &ssacparser.Result{Type: "User", Var: "user"}},
			{Type: "response", Target: "user"},
		},
	}}

	errs := checkResponseFields(funcs, st, doc, nil)
	for _, e := range errs {
		if e.Level != "WARNING" {
			t.Errorf("unexpected ERROR: %+v", e)
		}
	}
}

func TestCheckResponseFields_ShorthandWrapperSkip(t *testing.T) {
	// @response gigPage with Page wrapper → should be skipped (no errors)
	doc := buildResponseDoc("ListGigs", map[string]string{"items": "array", "total": "integer"})
	st := &ssacvalidator.SymbolTable{
		DDLTables: map[string]ssacvalidator.DDLTable{
			"gigs": {Columns: map[string]string{"id": "int64", "title": "string"}},
		},
	}
	funcs := []ssacparser.ServiceFunc{{
		Name:     "ListGigs",
		FileName: "list_gigs.ssac",
		Sequences: []ssacparser.Sequence{
			{Type: "get", Result: &ssacparser.Result{Type: "Gig", Var: "gigPage", Wrapper: "Page"}},
			{Type: "response", Target: "gigPage"},
		},
	}}

	errs := checkResponseFields(funcs, st, doc, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for wrapper type, got %d: %+v", len(errs), errs)
	}
}

func buildResponseDoc(opID string, props map[string]string) *openapi3.T {
	schemaProps := make(openapi3.Schemas)
	for name, typ := range props {
		schemaProps[name] = &openapi3.SchemaRef{Value: &openapi3.Schema{Type: &openapi3.Types{typ}}}
	}
	schema := &openapi3.Schema{
		Type:       &openapi3.Types{"object"},
		Properties: schemaProps,
	}
	ct := openapi3.NewContentWithJSONSchema(&openapi3.Schema{
		Type:       schema.Type,
		Properties: schema.Properties,
	})
	resp := openapi3.NewResponse().WithDescription("ok")
	resp.Content = ct

	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: resp})

	op := &openapi3.Operation{
		OperationID: opID,
		Responses:   responses,
	}

	paths := openapi3.NewPaths()
	paths.Set("/test", &openapi3.PathItem{Post: op})

	return &openapi3.T{Paths: paths}
}

func buildErrStatusDoc(opID string, responseCode string) *openapi3.T {
	resp := openapi3.NewResponse().WithDescription("response")
	responses := openapi3.NewResponses()
	responses.Set(responseCode, &openapi3.ResponseRef{Value: resp})

	op := &openapi3.Operation{
		OperationID: opID,
		Responses:   responses,
	}

	paths := openapi3.NewPaths()
	paths.Set("/test", &openapi3.PathItem{Post: op})

	return &openapi3.T{Paths: paths}
}
