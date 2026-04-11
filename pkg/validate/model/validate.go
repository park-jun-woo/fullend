//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — model/ 디렉토리 비어있음 검증 (M-1)
package model

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks that the model directory has at least one .go file.
func Validate(modelDir string) []validate.ValidationError {
	if modelDir == "" {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(modelDir, "*.go"))
	if len(matches) == 0 {
		return []validate.ValidationError{{
			Rule: "M-1", File: modelDir, Level: "ERROR",
			Message: "model/ directory is empty — at least package model declaration required", SeqIdx: -1,
		}}
	}
	for _, m := range matches {
		info, err := os.Stat(m)
		if err == nil && info.Size() > 0 {
			return nil
		}
	}
	return nil
}
