//ff:func feature=stml-validate type=rule control=iteration dimension=1
//ff:what 파싱된 PageSpec을 OpenAPI·custom.ts와 교차 검증
package validator

import (
	"fmt"
	"path/filepath"

	"github.com/geul-org/fullend/internal/stml/parser"
)

// Validate checks parsed PageSpecs against an OpenAPI spec and custom.ts files.
// projectRoot is the project root containing api/openapi.yaml and frontend/.
func Validate(pages []parser.PageSpec, projectRoot string) []ValidationError {
	openAPIPath := filepath.Join(projectRoot, "api", "openapi.yaml")
	st, err := LoadOpenAPI(openAPIPath)
	if err != nil {
		return []ValidationError{{
			File:    openAPIPath,
			Attr:    "openapi",
			Message: fmt.Sprintf("OpenAPI 파일을 읽을 수 없습니다: %v", err),
		}}
	}

	frontendDir := filepath.Join(projectRoot, "frontend")
	var errs []ValidationError

	for _, page := range pages {
		customPath := filepath.Join(frontendDir, page.Name+".custom.ts")
		cs, _ := LoadCustomTS(customPath)

		for _, f := range page.Fetches {
			errs = append(errs, validateFetchBlock(f, page.FileName, st, cs, frontendDir)...)
		}
		for _, a := range page.Actions {
			errs = append(errs, validateActionBlock(a, page.FileName, st, frontendDir)...)
		}
	}

	return errs
}
