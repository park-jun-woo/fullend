//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkOpCursor — 단일 operation의 cursor pagination + x-sort 검증
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func checkOpCursor(op *openapi3.Operation, path string) []CrossError {
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
		return nil
	}
	var xSort struct{ Allowed []string `json:"allowed"` }
	data, _ = json.Marshal(rawSort)
	if json.Unmarshal(data, &xSort) != nil {
		return nil
	}
	var errs []CrossError
	if len(xSort.Allowed) > 1 {
		errs = append(errs, CrossError{Rule: "X-7", Context: path, Level: "ERROR",
			Message: "cursor pagination with multiple x-sort allowed — runtime sort switching breaks cursor"})
	}
	return errs
}
