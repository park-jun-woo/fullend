//ff:type feature=orchestrator type=model
//ff:what TargetProfile defines the backend + frontend code generation targets.

package orchestrator

import (
	ssacgenerator "github.com/geul-org/fullend/internal/ssac/generator"
	stmlgenerator "github.com/geul-org/fullend/internal/stml/generator"
)

// TargetProfile defines the backend + frontend code generation targets.
type TargetProfile struct {
	Backend  ssacgenerator.Target
	Frontend stmlgenerator.Target
}
