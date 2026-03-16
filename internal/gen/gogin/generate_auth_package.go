//ff:func feature=gen-gogin type=generator control=sequence
//ff:what creates internal/auth/ with claims-based JWT functions and reexport

package gogin

import (
	"os"
	"path/filepath"

	"github.com/geul-org/fullend/internal/projectconfig"
)

// generateAuthPackage creates internal/auth/ with claims-based JWT functions and reexport.
func generateAuthPackage(intDir, modulePath string, claims map[string]projectconfig.ClaimDef, secretEnv string) error {
	authDir := filepath.Join(intDir, "auth")
	if err := os.MkdirAll(authDir, 0755); err != nil {
		return err
	}

	fields := sortedClaimFields(claims)
	claimsHash := HashClaimDefs(claims)

	if err := generateIssueToken(authDir, claims, fields, secretEnv, claimsHash); err != nil {
		return err
	}
	if err := generateVerifyToken(authDir, claims, fields, claimsHash); err != nil {
		return err
	}
	if err := generateRefreshToken(authDir, claims, fields, claimsHash); err != nil {
		return err
	}
	if err := generateReexport(authDir, claimsHash); err != nil {
		return err
	}

	return nil
}
