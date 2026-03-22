//ff:func feature=gen-hurl type=generator control=sequence
//ff:what Writes a Register + Login pair for a specific role.
package hurl

import (
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
		writeRegisterStep(buf, registerOp, registerPath, captures, role, emailPrefix, multiRole, checkEnums)
	}

	if loginOp != nil {
		writeLoginStep(buf, loginOp, loginPath, captures, role, emailPrefix, suffix, multiRole, checkEnums)
	}
}
