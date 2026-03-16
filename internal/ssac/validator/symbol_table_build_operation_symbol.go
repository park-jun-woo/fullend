//ff:func feature=symbol type=util control=iteration dimension=1
//ff:what 단일 OpenAPI operation에서 OperationSymbol을 구성한다
package validator

// buildOperationSymbol은 단일 OpenAPI operation에서 OperationSymbol을 구성한다.
func (st *SymbolTable) buildOperationSymbol(op *openAPIOperation, schemas map[string]openAPISchema) OperationSymbol {
	opSym := OperationSymbol{
		RequestFields: make(map[string]bool),
		XPagination:   op.XPagination,
		XSort:         op.XSort,
		XFilter:       op.XFilter,
		XInclude:      op.XInclude,
	}

	// path/query parameters
	for _, param := range op.Parameters {
		opSym.RequestFields[param.Name] = true
		if param.In != "path" {
			continue
		}
		opSym.PathParams = append(opSym.PathParams, PathParam{
			Name:   param.Name,
			GoType: oaTypeToGo(param.Schema.Type, param.Schema.Format),
		})
	}

	// request body fields
	if op.RequestBody == nil {
		return opSym
	}
	content, ok := op.RequestBody.Content["application/json"]
	if !ok {
		return opSym
	}
	for _, f := range collectSchemaFields(content.Schema, schemas) {
		opSym.RequestFields[f] = true
	}
	return opSym
}
