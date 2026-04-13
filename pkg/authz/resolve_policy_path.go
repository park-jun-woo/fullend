//ff:func feature=pkg-authz type=loader control=iteration dimension=1
//ff:what resolvePolicyPath — OPA_POLICY_PATH env 우선, 미지정 시 기본 경로 후보 자동 탐색

package authz

import (
	"os"
	"path/filepath"
)

// defaultPolicyPathCandidates lists fallback locations searched when OPA_POLICY_PATH is unset.
// Executable cwd relative.
var defaultPolicyPathCandidates = []string{
	"./internal/authz",
	"./authz",
	"./policy",
}

// resolvePolicyPath returns the OPA policy path: env OPA_POLICY_PATH first,
// then the first existing fallback candidate.
func resolvePolicyPath() (string, bool) {
	if p := os.Getenv("OPA_POLICY_PATH"); p != "" {
		return p, true
	}
	for _, cand := range defaultPolicyPathCandidates {
		abs, err := filepath.Abs(cand)
		if err != nil {
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			return abs, true
		}
	}
	return "", false
}
