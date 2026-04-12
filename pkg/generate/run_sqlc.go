//ff:func feature=rule type=command control=sequence
//ff:what runSqlc — sqlc 외부 도구로 DB 접근 코드 생성
package generate

import (
	"os/exec"
	"path/filepath"
)

func runSqlc(specsDir, artifactsDir string) []string {
	sqlcYaml := filepath.Join(specsDir, "sqlc.yaml")
	cmd := exec.Command("sqlc", "generate", "-f", sqlcYaml)
	if err := cmd.Run(); err != nil {
		// sqlc.yaml이 없으면 skip (오류 아님)
		return nil
	}
	_ = artifactsDir
	return nil
}
