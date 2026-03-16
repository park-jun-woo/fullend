//ff:type feature=stml-validate type=model
//ff:what OpenAPI 오퍼레이션 YAML 구조체
package validator

type openAPIOperation struct {
	OperationID string                     `yaml:"operationId"`
	Parameters  []openAPIParam             `yaml:"parameters"`
	RequestBody openAPIRequestBody         `yaml:"requestBody"`
	Responses   map[string]openAPIResponse `yaml:"responses"`
	XPagination *yamlPaginationExt         `yaml:"x-pagination"`
	XSort       *yamlSortExt               `yaml:"x-sort"`
	XFilter     *yamlFilterExt             `yaml:"x-filter"`
}
