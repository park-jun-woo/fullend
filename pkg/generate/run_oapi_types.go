//ff:func feature=rule type=command control=sequence
//ff:what runOapiTypes — oapi-codegen types 생성 호출
package generate

import (
	"os/exec"
	"path/filepath"
)

func runOapiTypes(apiPath, outDir string) string {
	out := filepath.Join(outDir, "types.gen.go")
	cmd := exec.Command("oapi-codegen", "-package", "api", "-generate", "types", "-o", out, apiPath)
	if err := cmd.Run(); err != nil {
		return "oapi-codegen types failed: " + err.Error()
	}
	return ""
}
