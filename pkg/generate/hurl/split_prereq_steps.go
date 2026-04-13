//ff:func feature=gen-hurl type=util control=iteration dimension=1
//ff:what splits mid steps into prerequisite creates and remaining steps

package hurl

import "github.com/getkin/kin-openapi/openapi3"

// splitPrereqSteps separates prerequisite top-level creates (needed before auth) from remaining mid steps.
func splitPrereqSteps(midSteps []orderedStep, authSteps []scenarioStep, doc *openapi3.T) (prereq, remain []orderedStep) {
	authFKPrefixes := collectAuthFKResources(authSteps, doc)
	for _, ms := range midSteps {
		if ms.step.Method == "POST" && ms.order < 0 && matchFKPrefix(inferResource(ms.step.Path), authFKPrefixes) {
			prereq = append(prereq, ms)
			continue
		}
		remain = append(remain, ms)
	}
	return prereq, remain
}
