//ff:func feature=stml-validate type=test-helper control=iteration dimension=1
//ff:what 임시 프로젝트 디렉토리를 생성하는 테스트 헬퍼
package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestProject(t *testing.T, openapi string, customTS map[string]string, components []string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "project")
	os.MkdirAll(filepath.Join(root, "api"), 0o755)
	os.MkdirAll(filepath.Join(root, "frontend", "components"), 0o755)
	os.WriteFile(filepath.Join(root, "api", "openapi.yaml"), []byte(openapi), 0o644)
	for name, content := range customTS { os.WriteFile(filepath.Join(root, "frontend", name), []byte(content), 0o644) }
	for _, comp := range components { os.WriteFile(filepath.Join(root, "frontend", "components", comp+".tsx"), []byte("export default {}"), 0o644) }
	return root
}
