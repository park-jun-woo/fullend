//ff:func feature=cli type=util control=iteration dimension=1
//ff:what --reset 플래그 파싱
package main

// parseResetFlag extracts --reset flag and returns (reset, remainingArgs).
func parseResetFlag(args []string) (bool, []string) {
	var remaining []string
	reset := false
	for _, a := range args {
		if a == "--reset" {
			reset = true
		} else {
			remaining = append(remaining, a)
		}
	}
	return reset, remaining
}
