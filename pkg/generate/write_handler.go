//ff:func feature=rule type=generator control=sequence
//ff:what writeHandler — 생성된 핸들러 코드를 파일에 기록
package generate

import (
	"os"
	"path/filepath"
	"strings"

	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
)

func writeHandler(serviceDir string, fn parsessac.ServiceFunc, body string) string {
	fileName := strings.TrimSuffix(fn.FileName, ".ssac") + ".go"
	path := filepath.Join(serviceDir, fileName)
	content := "//fullend:gen ssot=service/" + fn.FileName + "\npackage service\n\nfunc (s *Server) " + fn.Name + "(ctx context.Context, req " + fn.Name + "Request) (*" + fn.Name + "Response, error) {\n" + body + "}\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "write handler failed: " + err.Error()
	}
	return ""
}
