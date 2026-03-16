//ff:func feature=gen-hurl type=generator control=iteration dimension=2
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
			if ms, ok := classifyMidStep(s, stateOps, branchSkip, transitionOrder, resourceFirstTransition); ok {
				midSteps = append(midSteps, ms)
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

	prereqSteps, remainMidSteps := splitPrereqSteps(midSteps, authSteps, doc)

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
