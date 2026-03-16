//ff:func feature=gen-gogin type=generator control=iteration
//ff:what converts standalone function source to struct method with receiver, import fixes, and status replacement

package gogin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// transformSource converts a standalone function source to a struct method.
// isDomain: false → (s *Server) receiver, true → (h *Handler) receiver.
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

	// Replace model references: courseModel → {rcv}.CourseModel
	for _, m := range models {
		varName := lcFirst(m) + "Model"
		fieldName := ucFirst(varName)
		src = strings.ReplaceAll(src, varName+".", rcv+"."+fieldName+".")
	}

	// Replace func references: hashPassword( → {rcv}.HashPassword(
	for _, f := range funcs {
		fieldName := ucFirst(f)
		src = strings.ReplaceAll(src, f+"(", rcv+"."+fieldName+"(")
	}

	// Add model import when model package types are used (QueryOpts, ParseQueryOpts, CurrentUser, etc.)
	needsModel := strings.Contains(src, "model.QueryOpts") || strings.Contains(src, "model.ParseQueryOpts") || strings.Contains(src, "*model.CurrentUser")
	if needsModel {
		modelImport := fmt.Sprintf("\"%s/internal/model\"", modulePath)
		if !strings.Contains(src, modelImport) {
			src = strings.Replace(src, "import (", "import (\n\t"+modelImport, 1)
		}
	}

	// Fix state machine imports: "states/XXXstate" → "{modulePath}/internal/states/XXXstate"
	if strings.Contains(src, "\"states/") {
		src = strings.ReplaceAll(src, "\"states/", fmt.Sprintf("\"%s/internal/states/", modulePath))
	}

	// Fix authz import: "authz" → "github.com/geul-org/fullend/pkg/authz"
	if strings.Contains(src, "\t\"authz\"\n") {
		src = strings.ReplaceAll(src, "\t\"authz\"\n", "\t\"github.com/geul-org/fullend/pkg/authz\"\n")
	}

	// Fix queue import: "queue" → "github.com/geul-org/fullend/pkg/queue"
	if strings.Contains(src, "\t\"queue\"\n") {
		src = strings.ReplaceAll(src, "\t\"queue\"\n", "\t\"github.com/geul-org/fullend/pkg/queue\"\n")
	}

	// Fix config import: "config" → "github.com/geul-org/fullend/pkg/config"
	if strings.Contains(src, "\t\"config\"\n") {
		src = strings.ReplaceAll(src, "\t\"config\"\n", "\t\"github.com/geul-org/fullend/pkg/config\"\n")
	}

	// Fix auth import: pkg/auth → project internal/auth (reexport.go bridges pkg/auth utilities)
	if strings.Contains(src, "\"github.com/geul-org/fullend/pkg/auth\"") {
		src = strings.ReplaceAll(src, "\"github.com/geul-org/fullend/pkg/auth\"",
			fmt.Sprintf("\"%s/internal/auth\"", modulePath))
	}

	// Remove bare "model" import (already added as full path above)
	if strings.Contains(src, "\t\"model\"\n") {
		src = strings.ReplaceAll(src, "\t\"model\"\n", "")
	}

	// Add type assertions for @func results used as string arguments.
	for _, f := range funcs {
		callPattern := rcv + "." + ucFirst(f) + "("
		if idx := strings.Index(src, callPattern); idx > 0 {
			lineStart := strings.LastIndex(src[:idx], "\n") + 1
			assignLine := src[lineStart:idx]
			assignLine = strings.TrimSpace(assignLine)
			if commaIdx := strings.Index(assignLine, ","); commaIdx > 0 {
				varName := strings.TrimSpace(assignLine[:commaIdx])
				if varName != "_" && varName != "" {
					src = strings.ReplaceAll(src, ", "+varName+",", ", "+varName+".(string),")
					src = strings.ReplaceAll(src, ", "+varName+")", ", "+varName+".(string))")
					src = strings.ReplaceAll(src, "("+varName+",", "("+varName+".(string),")
				}
			}
		}
	}

	// Replace __RESPONSE_STATUS__ with OpenAPI success code.
	if strings.Contains(src, "__RESPONSE_STATUS__") && doc != nil && operationID != "" {
		statusConst := resolveSuccessStatus(doc, operationID)
		if statusConst != "" {
			if statusConst == "http.StatusNoContent" {
				// 204: replace entire c.JSON(__RESPONSE_STATUS__, ...) with c.Status(http.StatusNoContent)
				re := regexp.MustCompile(`c\.JSON\(__RESPONSE_STATUS__,\s*[^)]+\)`)
				src = re.ReplaceAllString(src, "c.Status(http.StatusNoContent)")
			} else {
				src = strings.ReplaceAll(src, "__RESPONSE_STATUS__", statusConst)
			}
		}
	}

	return src
}
