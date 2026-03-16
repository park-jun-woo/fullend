//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what ChildNode 트리를 순회하여 모든 ActionBlock을 수집한다
package generator

import "github.com/geul-org/fullend/internal/stml/parser"

// collectAllActions walks the ChildNode tree and collects all ActionBlocks.
func collectAllActions(nodes []parser.ChildNode) []parser.ActionBlock {
	var actions []parser.ActionBlock
	for _, ch := range nodes {
		switch ch.Kind {
		case "action":
			actions = append(actions, *ch.Action)
		case "fetch":
			actions = append(actions, collectAllActions(ch.Fetch.Children)...)
		case "state":
			actions = append(actions, collectAllActions(ch.State.Children)...)
		case "static":
			actions = append(actions, collectAllActions(ch.Static.Children)...)
		case "each":
			actions = append(actions, collectAllActions(ch.Each.Children)...)
		}
	}
	return actions
}
