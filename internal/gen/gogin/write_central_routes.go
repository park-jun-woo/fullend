//ff:func feature=gen-gogin type=generator control=iteration dimension=2 topic=http-handler
//ff:what writes route registrations from OpenAPI paths for central server

package gogin

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeCentralRoutes writes route registrations from OpenAPI paths.
func writeCentralRoutes(b *strings.Builder, doc *openapi3.T, opDomains map[string]string, hasBearer bool) {
	for pathStr, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			if op.OperationID == "" {
				continue
			}
			ginPath := convertPathParamsGin(pathStr)
			handlerName := op.OperationID

			domain := opDomains[handlerName]
			target := fmt.Sprintf("s.%s", handlerName)
			if domain != "" {
				target = fmt.Sprintf("s.%s.%s", ucFirst(domain), handlerName)
			}

			routerVar := "r"
			if opHasSecurity(op) && hasBearer {
				routerVar = "auth"
			}

			ginMethod := strings.ToUpper(method)
			b.WriteString(fmt.Sprintf("\t%s.Handle(%q, %q, %s)\n", routerVar, ginMethod, ginPath, target))
		}
	}
}
