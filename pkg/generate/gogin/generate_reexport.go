//ff:func feature=gen-gogin type=generator control=sequence topic=interface-derive
//ff:what creates internal/auth/reexport.go that re-exports pkg/auth utilities

package gogin

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/pkg/contract"
)

// generateReexport creates internal/auth/reexport.go that re-exports pkg/auth utilities.
func generateReexport(authDir string, claimsHash string) error {
	src := `package auth

import pkgauth "github.com/park-jun-woo/fullend/pkg/auth"

// Re-export pkg/auth utilities for unified import.
var HashPassword = pkgauth.HashPassword
var VerifyPassword = pkgauth.VerifyPassword
var GenerateResetToken = pkgauth.GenerateResetToken

type HashPasswordRequest = pkgauth.HashPasswordRequest
type HashPasswordResponse = pkgauth.HashPasswordResponse
type VerifyPasswordRequest = pkgauth.VerifyPasswordRequest
type VerifyPasswordResponse = pkgauth.VerifyPasswordResponse
type GenerateResetTokenRequest = pkgauth.GenerateResetTokenRequest
type GenerateResetTokenResponse = pkgauth.GenerateResetTokenResponse
`
	d := &contract.Directive{Ownership: "gen", SSOT: "fullend.yaml", Contract: claimsHash}
	src = injectFileDirective(src, d)
	return os.WriteFile(filepath.Join(authDir, "reexport.go"), []byte(src), 0644)
}
