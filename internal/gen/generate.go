//ff:func feature=genapi type=command
//ff:what parsed SSOT에서 backend + frontend + hurl 전체 코드를 생성한다
package gen

import (
	"github.com/geul-org/fullend/internal/gen/hurl"
	"github.com/geul-org/fullend/internal/gen/react"
	"github.com/geul-org/fullend/internal/genapi"
)

// Generate creates all artifacts from parsed SSOTs.
func Generate(parsed *genapi.ParsedSSOTs, cfg *genapi.GenConfig, stmlOut *genapi.STMLGenOutput) error {
	// 1. Backend code generation.
	backend := selectBackend(parsed.Config)
	if err := backend.Generate(parsed, cfg); err != nil {
		return err
	}
	// 2. React frontend (OpenAPI contract-based, backend-independent).
	if err := react.Generate(parsed, cfg, stmlOut); err != nil {
		return err
	}
	// 3. Hurl smoke tests (OpenAPI contract-based, backend-independent).
	if err := hurl.Generate(parsed, cfg); err != nil {
		return err
	}
	return nil
}
