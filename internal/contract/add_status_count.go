//ff:func feature=contract type=util control=selection
//ff:what 상태 문자열에 따라 해당 카운터를 1 증가시킨다
package contract

// addStatusCount increments the appropriate counter based on status.
func addStatusCount(status string, gen, preserve, broken, orphan int) (int, int, int, int) {
	switch status {
	case "gen":
		gen++
	case "preserve":
		preserve++
	case "broken":
		broken++
	case "orphan":
		orphan++
	}
	return gen, preserve, broken, orphan
}
