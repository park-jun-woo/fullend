//ff:func feature=orchestrator type=util
//ff:what DefaultProfile returns the default Go backend + React frontend profile.

package orchestrator

import (
	ssacgenerator "github.com/geul-org/fullend/internal/ssac/generator"
	stmlgenerator "github.com/geul-org/fullend/internal/stml/generator"
)

// DefaultProfile returns Go backend + React frontend.
func DefaultProfile() *TargetProfile {
	return &TargetProfile{
		Backend:  ssacgenerator.DefaultTarget(),
		Frontend: stmlgenerator.DefaultTarget(),
	}
}
