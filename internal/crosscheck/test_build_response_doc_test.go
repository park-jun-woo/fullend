//ff:func feature=crosscheck type=test-helper control=iteration dimension=1 topic=ssac-openapi
//ff:what 테스트용 OpenAPI 응답 문서 빌더

package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

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
