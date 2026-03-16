//ff:func feature=symbol type=loader control=iteration dimension=2
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
			opSym := st.buildOperationSymbol(op, schemas)
			st.Operations[op.OperationID] = opSym
		}
	}

	return nil
}
