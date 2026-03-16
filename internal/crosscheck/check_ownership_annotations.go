//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what Rego resource_owner 참조에 대한 @ownership 어노테이션 존재 여부 검증
package crosscheck

import (
	"fmt"

	"github.com/geul-org/fullend/internal/policy"
)

func checkOwnershipAnnotations(ownerResources map[string]bool, allOwnerships []policy.OwnershipMapping) []CrossError {
	ownershipMap := make(map[string]bool)
	for _, om := range allOwnerships {
		ownershipMap[om.Resource] = true
	}

	var errs []CrossError
	for res := range ownerResources {
		if !ownershipMap[res] {
			errs = append(errs, CrossError{
				Rule:       "Policy ↔ SSaC",
				Context:    fmt.Sprintf("resource=%s", res),
				Message:    fmt.Sprintf("Rego references input.resource_owner for %q but no @ownership annotation found", res),
				Level:      "ERROR",
				Suggestion: fmt.Sprintf("Add # @ownership %s: <table>.<column> to policy/*.rego", res),
			})
		}
	}
	return errs
}
