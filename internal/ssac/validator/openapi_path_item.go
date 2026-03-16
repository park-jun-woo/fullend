//ff:type feature=symbol type=model topic=openapi
//ff:what openAPIPathItem 타입 정의
package validator

type openAPIPathItem struct {
	Get    *openAPIOperation `yaml:"get"`
	Post   *openAPIOperation `yaml:"post"`
	Put    *openAPIOperation `yaml:"put"`
	Delete *openAPIOperation `yaml:"delete"`
}
