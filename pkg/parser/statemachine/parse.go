//ff:func feature=statemachine type=parser control=iteration dimension=1 topic=states
//ff:what Mermaid stateDiagram 텍스트를 파싱하여 StateDiagram 구조체를 반환한다
package statemachine

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	// [*] --> stateName
	reInitial = regexp.MustCompile(`\[\*\]\s*-->\s*(\w+)`)
	// stateA --> stateB: EventName
	reTransition = regexp.MustCompile(`(\w+)\s*-->\s*(\w+)\s*:\s*(\w+)`)
)

// Parse parses Mermaid stateDiagram content with a given ID.
func Parse(id, content string) (*StateDiagram, error) {
	// Extract content inside ```mermaid ... ``` block.
	mermaidContent := extractMermaidBlock(content)
	if mermaidContent == "" {
		return nil, fmt.Errorf("no mermaid stateDiagram block found in %s", id)
	}

	d := &StateDiagram{ID: id}

	// Parse initial state.
	if m := reInitial.FindStringSubmatch(mermaidContent); len(m) > 1 {
		d.InitialState = m[1]
	}

	// Parse transitions.
	matches := reTransition.FindAllStringSubmatch(mermaidContent, -1)
	stateSet := make(map[string]bool)
	if d.InitialState != "" {
		stateSet[d.InitialState] = true
	}

	for _, m := range matches {
		from, to, event := m[1], m[2], m[3]
		d.Transitions = append(d.Transitions, Transition{
			From:  from,
			To:    to,
			Event: event,
		})
		stateSet[from] = true
		stateSet[to] = true
	}

	// 상태명 대소문자 일관성 검증: case-insensitive로 같은 이름이면 에러
	lowerMap := make(map[string]string) // lowercase → first seen form
	for s := range stateSet {
		low := strings.ToLower(s)
		if prev, exists := lowerMap[low]; exists && prev != s {
			return nil, fmt.Errorf("state name conflict in %s: %q and %q differ only in case", id, prev, s)
		}
		lowerMap[low] = s
	}

	for s := range stateSet {
		d.States = append(d.States, s)
	}
	sort.Strings(d.States)

	if len(d.Transitions) == 0 {
		return nil, fmt.Errorf("no transitions found in stateDiagram %s", id)
	}

	return d, nil
}
