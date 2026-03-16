//ff:func feature=gen-hurl type=util control=sequence
//ff:what Builds query parameters from x- extensions (pagination, sort, filter).
package hurl

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// buildXQueryParams builds query parameters from x- extensions.
func buildXQueryParams(op *openapi3.Operation) string {
	var params []string

	if pag := getExtMap(op, "x-pagination"); pag != nil {
		params = append(params, "limit=2")
	}

	if sortCfg := getExtMap(op, "x-sort"); sortCfg != nil {
		def := getStr(sortCfg, "default", "")
		dir := getStr(sortCfg, "direction", "asc")
		if def != "" {
			params = append(params, "sort="+def)
			params = append(params, "direction="+dir)
		}
	}

	if filterCfg := getExtMap(op, "x-filter"); filterCfg != nil {
		allowed := getStrSlice(filterCfg, "allowed")
		if len(allowed) > 0 {
			params = append(params, allowed[0]+"=test_string")
		}
	}

	// x-include is codegen metadata only — no runtime query parameter.

	return strings.Join(params, "&")
}
