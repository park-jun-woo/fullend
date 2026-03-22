//ff:func feature=gen-hurl type=generator control=sequence
//ff:what Login API 호출 Hurl 스텝을 생성하고 토큰을 캡처한다

package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeLoginStep writes a Login Hurl step for a specific role and captures the token.
func writeLoginStep(buf *strings.Builder, loginOp *openapi3.Operation, loginPath string,
	captures map[string]bool, role, emailPrefix, suffix string, multiRole bool, checkEnums map[string][]string) {

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

	tokenField := findTokenJSONPath(getResponseSchema(loginOp))

	tokenVar := "token" + suffix
	buf.WriteString("[Captures]\n")
	buf.WriteString(fmt.Sprintf("%s: jsonpath \"$.%s\"\n", tokenVar, tokenField))
	captures[tokenVar] = true

	buf.WriteString("[Asserts]\n")
	buf.WriteString(fmt.Sprintf("jsonpath \"$.%s\" exists\n", tokenField))
	buf.WriteString("\n")
}
