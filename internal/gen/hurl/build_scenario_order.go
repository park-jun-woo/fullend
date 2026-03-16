//ff:func feature=gen-hurl type=generator control=iteration
//ff:what Scenario ordering logic — dependency-aware sorting of endpoints for smoke tests.
package hurl

import (
	"sort"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// buildScenarioOrder sorts endpoints into a dependency-aware execution order.
// Strategy:
//  1. Auth (Register, Login)
//  2. Interleaved creates + state transitions:
//     - Top-level creates first (depth <= 2)
//     - State transitions in stateDiagram BFS order
//     - Nested creates (depth > 2) after parent resource's first transition
//  3. Updates (PUT without @state)
//  4. Read (GET) last
//  5. DELETE in FK dependency order
func buildScenarioOrder(doc *openapi3.T, specsDir string, diagrams []*statemachine.StateDiagram, serviceFuncs []ssacparser.ServiceFunc) []scenarioStep {
	var all []scenarioStep

	for path, pi := range doc.Paths.Map() {
		depth := countPathSegments(path)
		for method, op := range pi.Operations() {
			if op == nil {
				continue
			}
			isAuth := isAuthOperation(op.OperationID)
			all = append(all, scenarioStep{
				OperationID: op.OperationID,
				Method:      method,
				Path:        path,
				Operation:   op,
				PathDepth:   depth,
				IsAuth:      isAuth,
			})
		}
	}

	// Build set of operationIDs that have @state in SSaC.
	stateOps := buildStateOpsSet(serviceFuncs)

	// Build transition order from stateDiagrams: event -> order index.
	transitionOrder := buildTransitionOrder(diagrams)

	// Build set of branching events to skip (from same state, keep only first).
	branchSkip := buildBranchSkipSet(diagrams, transitionOrder)

	// Build resource -> first transition order (for nested create placement).
	resourceFirstTransition := buildResourceFirstTransition(diagrams, transitionOrder)

	// Assign a sort key to each non-auth, non-read, non-delete step.
	type orderedStep struct {
		step  scenarioStep
		order float64
	}

	var authSteps []scenarioStep
	var midSteps []orderedStep
	var readSteps, deleteSteps []scenarioStep

	for _, s := range all {
		if s.IsAuth {
			authSteps = append(authSteps, s)
			continue
		}
		switch s.Method {
		case "GET":
			readSteps = append(readSteps, s)
		case "DELETE":
			deleteSteps = append(deleteSteps, s)
		default:
			if stateOps[s.OperationID] {
				// Skip branching transitions (only first path in smoke).
				if branchSkip[s.OperationID] {
					continue
				}
				// State transition: use diagram BFS order.
				midSteps = append(midSteps, orderedStep{s, float64(transitionOrder[s.OperationID])})
			} else if s.Method == "POST" {
				parentResource := findParentResource(s.Path)
				if parentResource == "" {
					// Top-level create: before all transitions.
					midSteps = append(midSteps, orderedStep{s, -1.0})
				} else if firstOrd, ok := resourceFirstTransition[parentResource]; ok {
					// Nested create: after parent's first transition.
					midSteps = append(midSteps, orderedStep{s, float64(firstOrd) + 0.5})
				} else {
					midSteps = append(midSteps, orderedStep{s, -0.5})
				}
			} else {
				// PUT/PATCH without @state: after transitions.
				midSteps = append(midSteps, orderedStep{s, 900.0})
			}
		}
	}

	// Sort auth: Register before Login.
	sort.SliceStable(authSteps, func(i, j int) bool {
		return authOrder(authSteps[i].OperationID) < authOrder(authSteps[j].OperationID)
	})

	// Sort mid steps by order, then depth, then path for stability.
	sort.SliceStable(midSteps, func(i, j int) bool {
		if midSteps[i].order != midSteps[j].order {
			return midSteps[i].order < midSteps[j].order
		}
		if midSteps[i].step.PathDepth != midSteps[j].step.PathDepth {
			return midSteps[i].step.PathDepth < midSteps[j].step.PathDepth
		}
		return midSteps[i].step.Path < midSteps[j].step.Path
	})

	// Sort reads by depth, then path.
	sortByDepthPath(readSteps)
	// Sort deletes by FK dependency.
	deleteSteps = sortDeletesByFK(deleteSteps, specsDir)

	// Detect auth FK prerequisites: if Register request body has _id fields,
	// move the corresponding top-level create endpoints before auth.
	var prereqSteps []orderedStep
	var remainMidSteps []orderedStep
	authFKPrefixes := collectAuthFKResources(authSteps, doc)
	for _, ms := range midSteps {
		if ms.step.Method == "POST" && ms.order < 0 {
			resource := inferResource(ms.step.Path)
			if matchFKPrefix(resource, authFKPrefixes) {
				prereqSteps = append(prereqSteps, ms)
				continue
			}
		}
		remainMidSteps = append(remainMidSteps, ms)
	}

	// Final order: prereq creates -> auth -> mid -> read -> delete
	var result []scenarioStep
	for _, ps := range prereqSteps {
		result = append(result, ps.step)
	}
	result = append(result, authSteps...)
	for _, ms := range remainMidSteps {
		result = append(result, ms.step)
	}
	result = append(result, readSteps...)
	result = append(result, deleteSteps...)

	return result
}
