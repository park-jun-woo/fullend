//ff:func feature=ssac-parse type=parser control=sequence
//ff:what 플랫 service/ 디렉토리에 .ssac 파일 배치 시 에러 검증

package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFlatServiceError(t *testing.T) {
	dir := t.TempDir()

	src := `package service

// @get User user = User.FindByEmail({Email: request.Email})
// @response {
//   user: user
// }
func Login() {}
`
	os.WriteFile(filepath.Join(dir, "login.ssac"), []byte(src), 0644)

	_, err := ParseDir(dir)
	if err == nil {
		t.Fatal("expected error for flat service/ file, got nil")
	}
	if !strings.Contains(err.Error(), "도메인 서브 폴더를 사용하세요") {
		t.Errorf("unexpected error: %v", err)
	}
}
