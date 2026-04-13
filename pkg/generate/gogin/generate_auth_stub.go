//ff:func feature=gen-gogin type=generator control=iteration dimension=1
//ff:what creates model/auth.go with CurrentUser type and Authorizer interface from claims config

package gogin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

// generateAuthStub creates model/auth.go with CurrentUser type and Authorizer interface.
// CurrentUser fields are derived from fullend.yaml claims config.
func generateAuthStub(intDir string, modulePath string, claims map[string]manifest.ClaimDef) error {
	modelDir := filepath.Join(intDir, "model")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString("package model\n\n")

	// Generate CurrentUser from claims config — claims are required when auth is present.
	// `authz:"<key>"` 태그는 authz.ClaimsFromStruct 가 JWT claim key 조회에 사용.
	b.WriteString("// CurrentUser is the authenticated user extracted by JWT middleware.\n")
	b.WriteString("type CurrentUser struct {\n")
	fields := sortedClaimFields(claims)
	for _, field := range fields {
		def := claims[field]
		b.WriteString(fmt.Sprintf("\t%s %s `authz:%q`\n", field, def.GoType, def.Key))
	}
	b.WriteString("}\n\n")

	b.WriteString("// Authorizer checks permissions.\n")
	b.WriteString("type Authorizer interface {\n")
	b.WriteString("\tCheck(user *CurrentUser, action, resource string, input interface{}) error\n")
	b.WriteString("}\n")

	return os.WriteFile(filepath.Join(modelDir, "auth.go"), []byte(b.String()), 0644)
}
