//ff:type feature=symbol type=model
//ff:what openAPIOperation 타입 정의
package validator

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
