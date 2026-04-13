//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 코드 생성 후 보존 콘텐츠를 복원한다
package contract

import "os"

// RestorePreserved restores all preserved content after code generation.
func RestorePreserved(snap *PreserveSnapshot) []Warning {
	var allWarnings []Warning

	// Restore file-level preserves.
	for path, content := range snap.FilePreserves {
		if _, err := os.Stat(path); err == nil {
			os.WriteFile(path, []byte(content), 0644)
		}
	}

	// Restore function-level preserves.
	for path, funcs := range snap.FuncPreserves {
		warnings := restoreFileFuncs(path, funcs)
		allWarnings = append(allWarnings, warnings...)
	}

	return allWarnings
}
