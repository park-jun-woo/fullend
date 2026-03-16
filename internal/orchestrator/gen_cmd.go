//ff:func feature=orchestrator type=command control=sequence
//ff:what Gen runs validate first, then generates code from all detected SSOTs.

package orchestrator

import (
	"github.com/geul-org/fullend/internal/reporter"
)

// Gen runs validate first, then generates code from all detected SSOTs.
// Returns the validate report (with gen steps appended) and whether gen succeeded.
func Gen(specsDir, artifactsDir string, skipKinds map[SSOTKind]bool, reset ...bool) (*reporter.Report, bool) {
	r := len(reset) > 0 && reset[0]
	return GenWith(DefaultProfile(), specsDir, artifactsDir, skipKinds, r)
}
