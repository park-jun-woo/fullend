//ff:func feature=symbol type=loader
//ff:what openapi.yaml에서 operationId별 request/response 필드를 추출한다
package validator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// loadOpenAPI는 openapi.yaml에서 operationId별 request/response 필드를 추출한다.
func (st *SymbolTable) loadOpenAPI(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var spec openAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return fmt.Errorf("YAML 파싱 실패: %w", err)
	}

	schemas := spec.Components.Schemas

	for _, pathItem := range spec.Paths {
		for _, op := range pathItem.operations() {
			if op.OperationID == "" {
				continue
			}

			opSym := OperationSymbol{
				RequestFields: make(map[string]bool),
				XPagination:   op.XPagination,
				XSort:          op.XSort,
				XFilter:        op.XFilter,
				XInclude:       op.XInclude,
			}

			// path/query parameters
			for _, param := range op.Parameters {
				opSym.RequestFields[param.Name] = true
				if param.In == "path" {
					opSym.PathParams = append(opSym.PathParams, PathParam{
						Name:   param.Name,
						GoType: oaTypeToGo(param.Schema.Type, param.Schema.Format),
					})
				}
			}

			// request body fields
			if op.RequestBody != nil {
				if content, ok := op.RequestBody.Content["application/json"]; ok {
					fields := collectSchemaFields(content.Schema, schemas)
					for _, f := range fields {
						opSym.RequestFields[f] = true
					}
				}
			}

			st.Operations[op.OperationID] = opSym
		}
	}

	return nil
}
