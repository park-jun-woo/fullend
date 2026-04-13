//ff:func feature=gen-gogin type=generator control=sequence
//ff:what generates auth package and middleware if claims are configured

package gogin

import (
	"fmt"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

// generateAuthIfNeeded generates auth package and middleware only when claims are configured.
func generateAuthIfNeeded(intDir, modulePath string, claims map[string]manifest.ClaimDef, secretEnv string) error {
	if len(claims) == 0 {
		return nil
	}
	if err := generateAuthPackage(intDir, modulePath, claims, secretEnv); err != nil {
		return fmt.Errorf("auth package (domain): %w", err)
	}
	if err := generateMiddleware(intDir, modulePath, claims); err != nil {
		return fmt.Errorf("middleware (domain): %w", err)
	}
	return nil
}
