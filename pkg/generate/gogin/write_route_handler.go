//ff:func feature=gen-gogin type=generator control=iteration dimension=1 topic=http-handler
//ff:what writes inline handler with path param extraction for a single route

package gogin

import (
	"fmt"
	"strings"
)

// writeRouteHandler writes an inline handler that extracts path parameters.
func writeRouteHandler(b *strings.Builder, pattern, handlerName string, pathParams []pathParamInfo) {
	b.WriteString(fmt.Sprintf("\tmux.HandleFunc(\"%s\", func(w http.ResponseWriter, r *http.Request) {\n", pattern))
	for _, pp := range pathParams {
		lcName := lcFirst(pp.GoName)
		if pp.IsInt {
			b.WriteString(fmt.Sprintf("\t\t%sStr := r.PathValue(\"%s\")\n", lcName, pp.Name))
			b.WriteString(fmt.Sprintf("\t\t%s, err := strconv.ParseInt(%sStr, 10, 64)\n", lcName, lcName))
			b.WriteString("\t\tif err != nil {\n")
			b.WriteString("\t\t\thttp.Error(w, \"invalid path parameter\", http.StatusBadRequest)\n")
			b.WriteString("\t\t\treturn\n")
			b.WriteString("\t\t}\n")
		} else {
			b.WriteString(fmt.Sprintf("\t\t%s := r.PathValue(\"%s\")\n", lcName, pp.Name))
		}
	}
	var args []string
	args = append(args, "w", "r")
	for _, pp := range pathParams {
		args = append(args, lcFirst(pp.GoName))
	}
	b.WriteString(fmt.Sprintf("\t\ts.%s(%s)\n", handlerName, strings.Join(args, ", ")))
	b.WriteString("\t})\n")
}
