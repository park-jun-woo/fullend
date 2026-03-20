//ff:func feature=manifest type=parser control=sequence
//ff:what 잘못된 YAML 파싱 시 에러 반환 검증
package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "fullend.yaml"), []byte(":\ninvalid: [yaml"), 0644)

	_, err := Load(dir)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
