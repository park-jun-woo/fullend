//ff:type feature=orchestrator type=model
//ff:what TargetProfile defines the backend + frontend code generation targets (pkg 경로).

package orchestrator

import (
	ssacgenerator "github.com/park-jun-woo/fullend/pkg/generate/gogin/ssac"
	stmlgenerator "github.com/park-jun-woo/fullend/pkg/generate/react/stml"
)

// TargetProfile defines the backend + frontend code generation targets.
type TargetProfile struct {
	Backend  ssacgenerator.Target
	Frontend stmlgenerator.Target
}
