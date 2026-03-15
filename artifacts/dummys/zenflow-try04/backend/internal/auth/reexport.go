//fullend:gen ssot=fullend.yaml contract=78bda02
package auth

import pkgauth "github.com/geul-org/fullend/pkg/auth"

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
