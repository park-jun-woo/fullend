//ff:func feature=cli type=util control=iteration dimension=1
//ff:what history 서브커맨드 인자 파싱
package main

// parseHistoryArgs parses history subcommand arguments and returns (target, format, all, quiet).
func parseHistoryArgs(args []string) (string, string, bool, bool) {
	var target string
	format := "yaml"
	var all, quiet bool
	skip := false

	for i, a := range args {
		if skip {
			skip = false
			continue
		}
		if a == "--format" && i+1 < len(args) {
			format = args[i+1]
			skip = true
			continue
		}
		if a == "--all" {
			all = true
			continue
		}
		if a == "-q" || a == "--quiet" {
			quiet = true
			continue
		}
		if target == "" {
			target = a
		}
	}

	return target, format, all, quiet
}
