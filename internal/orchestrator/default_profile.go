//ff:func feature=orchestrator type=util control=sequence
//ff:what DefaultProfile returns the default Go backend + React frontend profile (pkg 경로).

package orchestrator

import (
	ssacgenerator "github.com/park-jun-woo/fullend/pkg/generate/gogin/ssac"
	stmlgenerator "github.com/park-jun-woo/fullend/pkg/generate/react/stml"
)

// DefaultProfile returns Go backend + React frontend.
func DefaultProfile() *TargetProfile {
	return &TargetProfile{
		Backend:  ssacgenerator.DefaultTarget(),
		Frontend: stmlgenerator.DefaultTarget(),
	}
}
