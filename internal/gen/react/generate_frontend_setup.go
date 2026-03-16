//ff:func feature=gen-react type=generator control=sequence
//ff:what React + Vite 프로젝트 파일들을 생성한다

package react

import (
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateFrontendSetup creates React + Vite project files.
func generateFrontendSetup(artifactsDir string, doc *openapi3.T, stmlDeps map[string]string, stmlPages []string, stmlPageOps map[string]string) error {
	frontendDir := filepath.Join(artifactsDir, "frontend")
	srcDir := filepath.Join(frontendDir, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return err
	}

	if err := writePackageJSON(frontendDir, stmlDeps); err != nil {
		return err
	}
	if err := writeViteConfig(frontendDir); err != nil {
		return err
	}
	if err := writeTSConfig(frontendDir); err != nil {
		return err
	}
	if err := writeIndexHTML(frontendDir); err != nil {
		return err
	}
	if err := writeMainTSX(srcDir, stmlDeps); err != nil {
		return err
	}
	if err := writeAppTSX(srcDir, doc, stmlPages, stmlPageOps); err != nil {
		return err
	}
	if err := writeAPIClient(srcDir, doc); err != nil {
		return err
	}
	return nil
}
