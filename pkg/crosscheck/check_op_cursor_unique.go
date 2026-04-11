//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkOpCursorUnique — 단일 operation의 cursor sort default UNIQUE 검증 (X-8)
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkOpCursorUnique(g *rule.Ground, op *openapi3.Operation, path string) []CrossError {
	if op.Extensions == nil {
		return nil
	}
	rawPag, ok := op.Extensions["x-pagination"]
	if !ok {
		return nil
	}
	var pag struct{ Style string `json:"style"` }
	data, _ := json.Marshal(rawPag)
	if json.Unmarshal(data, &pag) != nil || pag.Style != "cursor" {
		return nil
	}
	rawSort, ok := op.Extensions["x-sort"]
	if !ok {
		return nil // defaults to id DESC — always UNIQUE
	}
	var xSort struct {
		Default string `json:"default"`
	}
	data, _ = json.Marshal(rawSort)
	if json.Unmarshal(data, &xSort) != nil || xSort.Default == "" {
		return nil
	}
	lookupKey := lookupKeyForPath(op)
	table := lookupKey[len("DDL.column."):]
	indexed := g.Lookup["DDL.index."+table]
	if !indexed[xSort.Default] {
		return []CrossError{{Rule: "X-8", Context: path, Level: "ERROR",
			Message: "cursor sort default " + xSort.Default + " is not a UNIQUE column"}}
	}
	return nil
}
