//ff:func feature=rule type=command control=sequence
//ff:what runOapiCodegen — oapi-codegen 외부 도구로 types/server 생성
package generate

import (
	"os"
	"os/exec"
	"path/filepath"
)

func runOapiCodegen(specsDir, artifactsDir string) []string {
	apiPath := filepath.Join(specsDir, "api", "openapi.yaml")
	outDir := filepath.Join(artifactsDir, "backend", "internal", "api")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return []string{"cannot create api dir: " + err.Error()}
	}
	var errs []string
	if err := runOapiTypes(apiPath, outDir); err != "" {
		errs = append(errs, err)
	}
	if err := runOapiServer(apiPath, outDir); err != "" {
		errs = append(errs, err)
	}
	_ = exec.Command // keep import
	return errs
}
