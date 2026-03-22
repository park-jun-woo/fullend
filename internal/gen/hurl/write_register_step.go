//ff:func feature=gen-hurl type=generator control=sequence
//ff:what Register API 호출 Hurl 스텝을 생성한다

package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeRegisterStep writes a Register Hurl step for a specific role.
func writeRegisterStep(buf *strings.Builder, registerOp *openapi3.Operation, registerPath string,
	captures map[string]bool, role, emailPrefix string, multiRole bool, checkEnums map[string][]string) {

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
	writeAssertLines(buf, asserts)
	buf.WriteString("\n")
}
