//ff:func feature=cli type=util control=iteration dimension=1
//ff:what map 서브커맨드 인자 파싱
package main

// parseMapArgs parses map subcommand arguments and returns (target, outputFile, force).
func parseMapArgs(args []string) (string, string, bool) {
	var target, outputFile string
	var force bool
	skip := false

	for i, a := range args {
		if skip {
			skip = false
			continue
		}
		if a == "-o" && i+1 < len(args) {
			outputFile = args[i+1]
			skip = true
			continue
		}
		if a == "-f" || a == "--force" {
			force = true
			continue
		}
		if target == "" {
			target = a
		}
	}

	return target, outputFile, force
}
