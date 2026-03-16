//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what OpenAPI YAML 파일을 파싱하여 SymbolTable 구성
package validator

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadOpenAPI parses an OpenAPI YAML file and builds a SymbolTable.
func LoadOpenAPI(path string) (*SymbolTable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read openapi: %w", err)
	}

	var doc openAPIDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse openapi: %w", err)
	}

	st := &SymbolTable{Operations: make(map[string]APISymbol)}

	for _, pathItem := range doc.Paths {
		collectOperations(pathItem, doc.Components.Schemas, st)
	}

	return st, nil
}
