//ff:func feature=genmodel type=generator control=sequence
//ff:what TestGenerateWriteFile: Generate로 파일 출력 후 escrow.go 존재 확인
package genmodel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	err := Generate("testdata/escrow.openapi.yaml", tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	outFile := filepath.Join(tmpDir, "escrow.go")
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Error("output file not created")
	}
}
