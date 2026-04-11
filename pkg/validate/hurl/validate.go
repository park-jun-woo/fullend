//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what Validate — .feature 파일 deprecated 검증 (H-1)
package hurl

import (
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/validate"
)

// Validate checks for deprecated .feature files.
func Validate(testsDir string) []validate.ValidationError {
	if testsDir == "" {
		return nil
	}
	matches, _ := filepath.Glob(filepath.Join(testsDir, "*.feature"))
	var errs []validate.ValidationError
	for _, m := range matches {
		errs = append(errs, validate.ValidationError{
			Rule: "H-1", File: strings.TrimPrefix(m, testsDir+"/"), Level: "ERROR",
			Message: ".feature files are deprecated — use .hurl format", SeqIdx: -1,
		})
	}
	return errs
}
