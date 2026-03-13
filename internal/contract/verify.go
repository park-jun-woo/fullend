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

// Summary returns counts by status.
func Summary(funcs []FuncStatus) (gen, preserve, broken, orphan int) {
	for _, f := range funcs {
		switch f.Status {
		case "gen":
			gen++
		case "preserve":
			preserve++
		case "broken":
			broken++
		case "orphan":
			orphan++
		}
	}
	return
}
