//ff:func feature=gen-hurl type=generator
//ff:what Step generation — writes a single endpoint test step with auth, body, captures, assertions.
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeStep writes a single endpoint test step.
func writeStep(buf *strings.Builder, step scenarioStep, captures map[string]bool, doc *openapi3.T, roleMap map[string]string, checkEnums map[string][]string) {
	op := step.Operation

	buf.WriteString(fmt.Sprintf("# %s\n", step.OperationID))

	// Build URL with path parameter substitution.
	url := substitutePathParams(step.Path, captures)

	// Add query parameters for x- extensions.
	queryParams := buildXQueryParams(op)
	if queryParams != "" {
		url += "?" + queryParams
	}

	buf.WriteString(fmt.Sprintf("%s {{host}}%s\n", step.Method, url))

	// Auth header: pick the right token based on operation's required role.
	if needsAuth(op) {
		tokenVar := resolveTokenVar(step.OperationID, roleMap, captures)
		if tokenVar != "" {
			buf.WriteString(fmt.Sprintf("Authorization: Bearer {{%s}}\n", tokenVar))
		}
	}

	// Request body for POST/PUT.
	var sentValues map[string]interface{}
	if step.Method == "POST" || step.Method == "PUT" {
		reqSchema := getRequestSchema(op)
		if reqSchema != nil {
			buf.WriteString("Content-Type: application/json\n")
			var body string
			body, sentValues = generateRequestBody(reqSchema, checkEnums)
			buf.WriteString(body + "\n")
		}
	}

	buf.WriteString(fmt.Sprintf("\nHTTP %s\n", getSuccessHTTPCode(op)))

	// Captures: extract ID from POST responses.
	if step.Method == "POST" {
		respSchema := getResponseSchema(op)
		varName, jsonPath := inferCaptureField(respSchema)
		if varName != "" && !captures[varName] {
			buf.WriteString("[Captures]\n")
			buf.WriteString(fmt.Sprintf("%s: jsonpath %q\n", varName, jsonPath))
			captures[varName] = true
		}
	}

	// Assertions.
	respSchema := getResponseSchema(op)
	asserts := generateResponseAssertions(respSchema, sentValues)
	if len(asserts) > 0 {
		buf.WriteString("[Asserts]\n")
		for _, a := range asserts {
			buf.WriteString(a + "\n")
		}
	}

	buf.WriteString("\n")
}
