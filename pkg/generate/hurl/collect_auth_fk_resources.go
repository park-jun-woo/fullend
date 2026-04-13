//ff:func feature=gen-hurl type=util control=iteration dimension=2
//ff:what Detects FK-like fields in auth operation request bodies for pre-auth resource creation.
package hurl

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// collectAuthFKResources detects FK-like fields in auth operation request bodies
// and returns the set of plural resource names that must be created before auth.
// e.g. Register body has "org_id" -> returns {"organizations": true}
func collectAuthFKResources(authSteps []scenarioStep, doc *openapi3.T) []string {
	var fkPrefixes []string
	for _, s := range authSteps {
		reqSchema := getRequestSchema(s.Operation)
		if reqSchema == nil {
			continue
		}
		for name := range reqSchema.Properties {
			if strings.HasSuffix(name, "_id") {
				// org_id -> "org"
				fkPrefixes = append(fkPrefixes, strings.TrimSuffix(name, "_id"))
			}
		}
	}
	return fkPrefixes
}
