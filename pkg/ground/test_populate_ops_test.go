//ff:func feature=rule type=loader control=sequence
//ff:what populateOps 검증 — OpenAPI operation → Ground.Ops
package ground

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func TestPopulateOps(t *testing.T) {
	loader := openapi3.NewLoader()
	spec := []byte(`
openapi: 3.0.3
info: {title: test, version: "1.0"}
paths:
  /users/{id}:
    get:
      operationId: GetUser
      parameters:
        - name: id
          in: path
          required: true
          schema: {type: integer, format: int64}
      responses: {"200": {description: ok}}
  /users:
    post:
      operationId: CreateUser
      requestBody:
        required: true
        content:
          application/json:
            schema: {type: object}
      responses: {"201": {description: created}}
    get:
      operationId: ListUsers
      x-pagination:
        style: cursor
        defaultLimit: 20
        maxLimit: 100
      x-sort:
        allowed: [created_at, email]
        default: created_at
        direction: desc
      x-filter:
        allowed: [status]
      responses: {"200": {description: ok}}
`)
	doc, err := loader.LoadFromData(spec)
	if err != nil {
		t.Fatal(err)
	}

	g := newGround()
	fs := &fullend.Fullstack{OpenAPIDoc: doc}
	populateOps(g, fs)

	get, ok := g.Ops["GetUser"]
	if !ok {
		t.Fatal("GetUser not populated")
	}
	if get.Method != "GET" || get.Path != "/users/{id}" {
		t.Errorf("GetUser meta: %+v", get)
	}
	if len(get.PathParams) != 1 || get.PathParams[0].Name != "id" || get.PathParams[0].GoType != "int64" {
		t.Errorf("GetUser path params: %+v", get.PathParams)
	}

	create, ok := g.Ops["CreateUser"]
	if !ok || !create.HasRequestBody {
		t.Errorf("CreateUser HasRequestBody: %+v", create)
	}

	list, ok := g.Ops["ListUsers"]
	if !ok {
		t.Fatal("ListUsers not populated")
	}
	if list.Pagination == nil || list.Pagination.Style != "cursor" || list.Pagination.DefaultLimit != 20 {
		t.Errorf("Pagination: %+v", list.Pagination)
	}
	if list.Sort == nil || len(list.Sort.Allowed) != 2 || list.Sort.Default != "created_at" || list.Sort.Direction != "desc" {
		t.Errorf("Sort: %+v", list.Sort)
	}
	if list.Filter == nil || len(list.Filter.Allowed) != 1 {
		t.Errorf("Filter: %+v", list.Filter)
	}
}
