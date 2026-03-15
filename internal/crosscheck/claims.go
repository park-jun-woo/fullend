package crosscheck

import (
	"fmt"
	"strings"

	"github.com/geul-org/fullend/internal/policy"
	"github.com/geul-org/fullend/internal/projectconfig"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// CheckClaims validates that all currentUser field references in SSaC specs
// are defined in fullend.yaml backend.auth.claims.
func CheckClaims(serviceFuncs []ssacparser.ServiceFunc, claims map[string]projectconfig.ClaimDef) []CrossError {
	var errs []CrossError

	// Collect all currentUser field usages with locations.
	usedFields := collectCurrentUserFields(serviceFuncs)

	// Check if @auth is used (implies currentUser in scope).
	hasAuth := false
	for _, sf := range serviceFuncs {
		for _, seq := range sf.Sequences {
			if seq.Type == "auth" {
				hasAuth = true
				break
			}
		}
		if hasAuth {
			break
		}
	}

	// If currentUser is used but no claims config exists → ERROR.
	if (len(usedFields) > 0 || hasAuth) && claims == nil {
		errs = append(errs, CrossError{
			Rule:    "Claims ↔ SSaC",
			Context: "fullend.yaml",
			Message: "currentUser를 사용하지만 backend.auth.claims가 정의되지 않았습니다",
			Level:   "ERROR",
		})
		return errs
	}

	if claims == nil {
		return errs
	}

	// Check each used field exists in claims.
	for field, locations := range usedFields {
		if _, ok := claims[field]; !ok {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ SSaC",
				Context: strings.Join(locations, ", "),
				Message: fmt.Sprintf("currentUser.%s — claims에 미정의", field),
				Level:   "ERROR",
			})
		}
	}

	return errs
}

// CheckClaimsRego validates that Rego input.claims.xxx references match fullend.yaml claims values.
func CheckClaimsRego(policies []*policy.Policy, claims map[string]projectconfig.ClaimDef) []CrossError {
	if claims == nil {
		return nil
	}

	// Build set of claim keys (e.g., "user_id", "email", "role").
	claimValues := make(map[string]bool)
	for _, def := range claims {
		claimValues[def.Key] = true
	}

	// Collect all Rego claims refs across policies.
	regoRefs := make(map[string]string) // ref → file
	for _, p := range policies {
		for _, ref := range p.ClaimsRefs {
			if _, exists := regoRefs[ref]; !exists {
				regoRefs[ref] = p.File
			}
		}
	}

	var errs []CrossError

	// Forward: Rego claims ref → fullend.yaml claims value
	for ref, file := range regoRefs {
		if !claimValues[ref] {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ Rego",
				Context: file,
				Message: fmt.Sprintf("Rego input.claims.%s — fullend.yaml claims 값에 %q 없음", ref, ref),
				Level:   "ERROR",
			})
		}
	}

	// Reverse: fullend.yaml claim keys → Rego (WARNING if unused)
	for _, def := range claims {
		if _, used := regoRefs[def.Key]; !used {
			errs = append(errs, CrossError{
				Rule:    "Claims ↔ Rego",
				Context: "fullend.yaml",
				Message: fmt.Sprintf("claims 값 %q — Rego에서 미참조", def.Key),
				Level:   "WARNING",
			})
		}
	}

	return errs
}

// collectCurrentUserFields scans all SSaC sequences for currentUser.X references.
// Returns map[fieldName][]location.
func collectCurrentUserFields(funcs []ssacparser.ServiceFunc) map[string][]string {
	result := make(map[string][]string)

	for _, sf := range funcs {
		loc := sf.FileName + ":" + sf.Name
		for _, seq := range sf.Sequences {
			// Inputs: values starting with "currentUser."
			for _, val := range seq.Inputs {
				if strings.HasPrefix(val, "currentUser.") {
					field := strings.TrimPrefix(val, "currentUser.")
					result[field] = append(result[field], loc)
				}
			}
		}
	}

	return result
}
