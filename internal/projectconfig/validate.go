//ff:func feature=projectconfig type=util control=sequence
//ff:what ProjectConfig의 필수 필드와 인증 설정을 검증한다
package projectconfig

import "fmt"

// allowedClaimTypes is the set of Go types allowed in claims definitions.
var allowedClaimTypes = map[string]bool{
	"string": true,
	"int64":  true,
	"bool":   true,
}

// jwtReservedKeys are standard JWT claim keys that must not be used as custom claim keys.
var jwtReservedKeys = map[string]bool{
	"exp": true, "iat": true, "sub": true, "iss": true,
	"aud": true, "nbf": true, "jti": true,
}

// Validate checks required fields.
func (c *ProjectConfig) Validate() error {
	if c.APIVersion != "fullend/v1" {
		return fmt.Errorf("fullend.yaml: apiVersion must be \"fullend/v1\", got %q", c.APIVersion)
	}
	if c.Kind != "Project" {
		return fmt.Errorf("fullend.yaml: kind must be \"Project\", got %q", c.Kind)
	}
	if c.Metadata.Name == "" {
		return fmt.Errorf("fullend.yaml: metadata.name is required")
	}
	if c.Backend.Module == "" {
		return fmt.Errorf("fullend.yaml: backend.module is required")
	}

	// Auth section validation.
	if c.Backend.Auth != nil {
		auth := c.Backend.Auth
		if auth.Type == "" {
			return fmt.Errorf("fullend.yaml: auth.type is required (supported: jwt)")
		}
		if auth.Type != "jwt" {
			return fmt.Errorf("fullend.yaml: auth.type %q is not supported (supported: jwt)", auth.Type)
		}
		if len(auth.Claims) == 0 {
			return fmt.Errorf("fullend.yaml: auth.claims must have at least 1 entry")
		}
		if err := validateClaims(auth); err != nil {
			return err
		}
	}

	return nil
}
