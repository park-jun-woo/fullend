//ff:func feature=stml-validate type=parser control=sequence
//ff:what 단일 OpenAPI 오퍼레이션에서 APISymbol 구성
package validator

import "strings"

// buildAPISymbol constructs an APISymbol from a single OpenAPI operation.
func buildAPISymbol(method string, op openAPIOperation, schemas map[string]openAPISchema) APISymbol {
	api := APISymbol{
		Method:         strings.ToUpper(method),
		RequestFields:  make(map[string]string),
		ResponseFields: make(map[string]FieldSymbol),
	}

	collectParams(op, &api)
	collectRequestFields(op, schemas, api.RequestFields)
	collectResponseFields(op, schemas, api.ResponseFields)
	applyExtensions(op, &api)

	return api
}
