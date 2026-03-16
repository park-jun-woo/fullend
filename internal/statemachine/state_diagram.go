//ff:type feature=statemachine type=model
//ff:what Mermaid stateDiagram 파싱 결과를 담는 구조체
package statemachine

// StateDiagram represents a parsed Mermaid stateDiagram.
type StateDiagram struct {
	ID           string       // derived from filename (e.g. "course")
	InitialState string       // state after [*] -->
	States       []string     // all unique state names
	Transitions  []Transition // all state transitions
}
