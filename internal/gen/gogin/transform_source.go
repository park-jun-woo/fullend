//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what converts standalone function source to struct method with receiver, import fixes, and status replacement

package gogin

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// transformSource converts a standalone function source to a struct method.
// isDomain: false -> (s *Server) receiver, true -> (h *Handler) receiver.
// doc + operationID: used to replace __RESPONSE_STATUS__ with OpenAPI success code.
func transformSource(src string, models, funcs []string, modulePath string, isDomain bool, doc *openapi3.T, operationID string) string {
	rcv := "s"
	rcvType := "*Server"
	if isDomain {
		rcv = "h"
		rcvType = "*Handler"
	}

	// Add receiver to function declaration (skip if already has one).
	if idx := strings.Index(src, "\nfunc "); idx >= 0 {
		after := src[idx+len("\nfunc "):]
		if len(after) == 0 || after[0] != '(' {
			src = strings.Replace(src, "\nfunc ", "\nfunc ("+rcv+" "+rcvType+") ", 1)
		}
	}

	// Replace model references: courseModel -> {rcv}.CourseModel
	for _, m := range models {
		varName := lcFirst(m) + "Model"
		fieldName := ucFirst(varName)
		src = strings.ReplaceAll(src, varName+".", rcv+"."+fieldName+".")
	}

	// Replace func references: hashPassword( -> {rcv}.HashPassword(
	for _, f := range funcs {
		fieldName := ucFirst(f)
		src = strings.ReplaceAll(src, f+"(", rcv+"."+fieldName+"(")
	}

	src = fixImports(src, modulePath)
	src = addTypeAssertions(src, rcv, funcs)

	return replaceResponseStatus(src, doc, operationID)
}
