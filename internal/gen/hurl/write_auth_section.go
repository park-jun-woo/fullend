//ff:func feature=gen-hurl type=generator
//ff:what Auth generation — Register + Login steps with role-specific tokens and FK resolution.
package hurl

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// writeAuthSection writes Register + Login steps for each required role.
// When multiple roles are needed, creates separate users per role with
// role-suffixed tokens (e.g. token_client, token_freelancer).
func writeAuthSection(buf *strings.Builder, doc *openapi3.T, captures map[string]bool, roles []string, checkEnums map[string][]string) {
	buf.WriteString("# ===== Auth =====\n\n")

	// Find Register and Login operations.
	var registerOp, loginOp *openapi3.Operation
	var registerPath, loginPath string

	for path, pi := range doc.Paths.Map() {
		for _, op := range pi.Operations() {
			if op == nil {
				continue
			}
			switch strings.ToLower(op.OperationID) {
			case "register":
				registerOp = op
				registerPath = path
			case "login":
				loginOp = op
				loginPath = path
			}
		}
	}

	// If no roles detected or only one role, use single auth flow.
	if len(roles) <= 1 {
		role := ""
		if len(roles) == 1 {
			role = roles[0]
		}
		writeAuthPair(buf, registerOp, registerPath, loginOp, loginPath, captures, role, false, checkEnums)
		return
	}

	// Multi-role: register + login for each role.
	for _, role := range roles {
		writeAuthPair(buf, registerOp, registerPath, loginOp, loginPath, captures, role, true, checkEnums)
	}
}
