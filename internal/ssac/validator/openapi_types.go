//ff:type feature=symbol type=model
//ff:what OpenAPI YAML 구조체 (openAPISpec, Components, Schema, PathItem, Operation, Parameter, RequestBody, Response, MediaType)
package validator

// --- OpenAPI YAML 구조체 ---

type openAPISpec struct {
	Paths      map[string]openAPIPathItem `yaml:"paths"`
	Components openAPIComponents          `yaml:"components"`
}

type openAPIComponents struct {
	Schemas map[string]openAPISchema `yaml:"schemas"`
}

type openAPISchema struct {
	Type       string                   `yaml:"type"`
	Format     string                   `yaml:"format"`
	Properties map[string]openAPISchema `yaml:"properties"`
	Ref        string                   `yaml:"$ref"`
}

type openAPIPathItem struct {
	Get    *openAPIOperation `yaml:"get"`
	Post   *openAPIOperation `yaml:"post"`
	Put    *openAPIOperation `yaml:"put"`
	Delete *openAPIOperation `yaml:"delete"`
}

func (p openAPIPathItem) operations() []*openAPIOperation {
	var ops []*openAPIOperation
	for _, op := range []*openAPIOperation{p.Get, p.Post, p.Put, p.Delete} {
		if op != nil {
			ops = append(ops, op)
		}
	}
	return ops
}

type openAPIOperation struct {
	OperationID string                     `yaml:"operationId"`
	Parameters  []openAPIParameter         `yaml:"parameters"`
	RequestBody *openAPIRequestBody        `yaml:"requestBody"`
	Responses   map[string]openAPIResponse `yaml:"responses"`
	XPagination *XPagination               `yaml:"x-pagination"`
	XSort       *XSort                     `yaml:"x-sort"`
	XFilter     *XFilter                   `yaml:"x-filter"`
	XInclude    *XInclude                  `yaml:"x-include"`
}

type openAPIParameter struct {
	Name   string          `yaml:"name"`
	In     string          `yaml:"in"`
	Schema openAPISchema   `yaml:"schema"`
}

type openAPIRequestBody struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}

type openAPIResponse struct {
	Content map[string]openAPIMediaType `yaml:"content"`
}

type openAPIMediaType struct {
	Schema openAPISchema `yaml:"schema"`
}
