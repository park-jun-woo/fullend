//ff:func feature=contract type=util control=sequence
//ff:what 단일 파일의 보존 함수들을 복원하고 경고를 반환한다
package contract

import "os"

// restoreFileFuncs restores preserved functions for a single file.
func restoreFileFuncs(path string, funcs map[string]*PreservedFunc) []Warning {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	result, err := SpliceWithPreserved(string(src), funcs, path)
	if err != nil {
		return nil
	}
	os.WriteFile(path, []byte(result.Content), 0644)

	// Write .new file for contract changes.
	if len(result.Warnings) > 0 {
		newPath := path + ".new"
		os.WriteFile(newPath, src, 0644)
	}

	return result.Warnings
}
