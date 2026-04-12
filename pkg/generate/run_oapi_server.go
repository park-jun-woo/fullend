//ff:func feature=rule type=command control=sequence
//ff:what runOapiServer — oapi-codegen std-http-server 생성 호출
package generate

import (
	"os/exec"
	"path/filepath"
)

func runOapiServer(apiPath, outDir string) string {
	out := filepath.Join(outDir, "server.gen.go")
	cmd := exec.Command("oapi-codegen", "-package", "api", "-generate", "std-http-server", "-o", out, apiPath)
	if err := cmd.Run(); err != nil {
		return "oapi-codegen server failed: " + err.Error()
	}
	return ""
}
