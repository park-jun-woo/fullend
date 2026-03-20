//ff:func feature=crosscheck type=test-helper control=sequence topic=ssac-openapi
//ff:what 테스트용 OpenAPI 에러 상태 문서 빌더

package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

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
