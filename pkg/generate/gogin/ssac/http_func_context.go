//ff:type feature=ssac-gen type=model topic=http-handler
//ff:what HTTP 함수 생성에 필요한 분석 컨텍스트 구조체
package ssac

import "github.com/park-jun-woo/fullend/internal/ssac/validator"

type httpFuncContext struct {
	pathParams    []validator.PathParam
	pathParamSet  map[string]bool
	requestParams []typedRequestParam
	needsCU       bool
	needsQO       bool
	pkgName       string
	resolver      *FieldTypeResolver
	resultTypes   map[string]string
	varSources    map[string]varSource
}
