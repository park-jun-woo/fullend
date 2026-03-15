//ff:func feature=gen-hurl type=generator
//ff:what Writes a Register + Login pair for a specific role.
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeAuthPair writes a Register + Login pair for a specific role.
// If multiRole is true, token is captured as token_<role>.
func writeAuthPair(buf *strings.Builder, registerOp *openapi3.Operation, registerPath string,
	loginOp *openapi3.Operation, loginPath string, captures map[string]bool, role string, multiRole bool, checkEnums map[string][]string) {

	suffix := ""
	emailPrefix := "test"
	if multiRole && role != "" {
		suffix = "_" + role
		emailPrefix = role
	}

	if registerOp != nil {
		if multiRole {
			buf.WriteString(fmt.Sprintf("# Register (%s)\n", role))
		} else {
			buf.WriteString("# Register\n")
		}
		buf.WriteString(fmt.Sprintf("POST {{host}}%s\n", registerPath))
		buf.WriteString("Content-Type: application/json\n")
		reqSchema := getRequestSchema(registerOp)
		body := generateRequestBodyWithOverrides(reqSchema, role, emailPrefix, checkEnums, captures)
		buf.WriteString(body + "\n")
		buf.WriteString(fmt.Sprintf("\nHTTP %s\n", getSuccessHTTPCode(registerOp)))

		respSchema := getResponseSchema(registerOp)
		asserts := generateResponseAssertions(respSchema, nil)
		if len(asserts) > 0 {
			buf.WriteString("[Asserts]\n")
			for _, a := range asserts {
				buf.WriteString(a + "\n")
			}
		}
		buf.WriteString("\n")
	}

	if loginOp != nil {
		if multiRole {
			buf.WriteString(fmt.Sprintf("# Login (%s)\n", role))
		} else {
			buf.WriteString("# Login\n")
		}
		buf.WriteString(fmt.Sprintf("POST {{host}}%s\n", loginPath))
		buf.WriteString("Content-Type: application/json\n")
		reqSchema := getRequestSchema(loginOp)
		body := generateLoginBodyWithEmail(reqSchema, emailPrefix, checkEnums)
		buf.WriteString(body + "\n")
		buf.WriteString(fmt.Sprintf("\nHTTP %s\n", getSuccessHTTPCode(loginOp)))

		// Find token field path from response schema (handles nested objects).
		tokenField := findTokenJSONPath(getResponseSchema(loginOp))

		tokenVar := "token" + suffix
		buf.WriteString("[Captures]\n")
		buf.WriteString(fmt.Sprintf("%s: jsonpath \"$.%s\"\n", tokenVar, tokenField))
		captures[tokenVar] = true

		buf.WriteString("[Asserts]\n")
		buf.WriteString(fmt.Sprintf("jsonpath \"$.%s\" exists\n", tokenField))
		buf.WriteString("\n")
	}
}
