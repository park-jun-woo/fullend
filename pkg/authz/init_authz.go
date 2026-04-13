//ff:func feature=pkg-authz type=loader control=sequence
//ff:what 글로벌 인가 상태를 초기화 — OPA_POLICY_PATH env 또는 기본 경로 탐색

package authz

import (
	"database/sql"
	"fmt"
	"os"
)

var globalPolicy string
var globalDB *sql.DB
var globalOwnerships []OwnershipMapping

// Init initializes the global authz state.
// OPA policy source resolution:
//  1. OPA_POLICY_PATH env — file or directory (directory loads all *.rego)
//  2. fallback: ./internal/authz, ./authz, ./policy (first existing directory)
// DISABLE_AUTHZ=1 전체 skip.
func Init(db *sql.DB, ownerships []OwnershipMapping) error {
	globalDB = db
	globalOwnerships = ownerships

	if os.Getenv("DISABLE_AUTHZ") == "1" {
		return nil
	}

	policyPath, ok := resolvePolicyPath()
	if !ok {
		return fmt.Errorf("OPA_POLICY_PATH env not set and no fallback path exists " +
			"(tried ./internal/authz, ./authz, ./policy — set DISABLE_AUTHZ=1 to skip)")
	}

	data, err := loadPolicyFromPath(policyPath)
	if err != nil {
		return err
	}
	globalPolicy = data
	return nil
}
