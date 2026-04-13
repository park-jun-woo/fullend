//ff:func feature=contract type=util control=iteration dimension=1
//ff:what FuncStatus 목록을 SSOT 파일과 대조하여 상태를 갱신한다
package contract

import (
	"os"
	"path/filepath"
)

// Verify checks each FuncStatus against the current SSOT files.
// Updates Status to "broken" (hash mismatch) or "orphan" (SSOT deleted).
func Verify(specsDir string, funcs []FuncStatus) []FuncStatus {
	result := make([]FuncStatus, len(funcs))
	copy(result, funcs)

	for i := range result {
		ssotPath := filepath.Join(specsDir, result[i].Directive.SSOT)

		// Check if SSOT file exists.
		if _, err := os.Stat(ssotPath); os.IsNotExist(err) {
			result[i].Status = "orphan"
			result[i].Detail = "SSOT 삭제됨"
			continue
		}

		// Note: contract hash re-computation requires parsing the SSOT file,
		// which depends on the SSOT type. For now, we keep the existing status
		// and rely on fullend gen to detect contract changes via the splice engine.
		// A full verify would parse the SSOT and recompute the hash.
	}

	return result
}
