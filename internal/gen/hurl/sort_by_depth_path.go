//ff:func feature=gen-hurl type=util
//ff:what Sorts scenario steps by path depth then path string.
package hurl

import "sort"

func sortByDepthPath(steps []scenarioStep) {
	sort.SliceStable(steps, func(i, j int) bool {
		if steps[i].PathDepth != steps[j].PathDepth {
			return steps[i].PathDepth < steps[j].PathDepth
		}
		return steps[i].Path < steps[j].Path
	})
}
