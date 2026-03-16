//ff:func feature=stml-validate type=rule control=sequence
//ff:what 컴포넌트 TSX 파일이 존재하는지 검증
package validator

import (
	"os"
	"path/filepath"
)

func validateComponent(name, file, frontendDir string) []ValidationError {
	compPath := filepath.Join(frontendDir, "components", name+".tsx")
	if _, err := os.Stat(compPath); os.IsNotExist(err) {
		relPath := filepath.Join("frontend", "components", name+".tsx")
		return []ValidationError{errComponentNotFound(file, name, relPath)}
	}
	return nil
}
