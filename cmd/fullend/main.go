//ff:func feature=cli type=command control=selection
//ff:what CLI 엔트리포인트 — 서브커맨드 분기
package main

import (
	"fmt"
	"os"

	"github.com/geul-org/fullend/internal/genmodel"
)

// Version is set at build time via -ldflags.
var Version = "dev"

const usage = `Usage: fullend <command> [arguments]

Commands:
  validate   [--skip kind,...] <specs-dir>                 Validate SSOT specs
  gen        [--skip kind,...] [--reset] <specs-dir> <artifacts-dir> Generate code from specs
  gen-model  <openapi-source> <output-dir>                 Generate Go model from external OpenAPI
  status     <specs-dir>                                   Show SSOT status summary
  contract   <specs-dir> <artifacts-dir>                   Show contract status
  chain      <operationId> <specs-dir>                     Trace feature chain for an operation
  map        [path]                                        Generate keyword map
  history    <file|dir> [--all] [--format yaml|json]       Show file change history
  version                                                  Show version

Skip kinds: openapi, ddl, ssac, model, stml, states, policy, scenario, func
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("fullend %s\n", Version)
	case "validate":
		skipKinds, args := parseSkipFlag(os.Args[2:])
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: fullend validate [--skip kind,...] <specs-dir>")
			os.Exit(2)
		}
		runValidate(args[0], skipKinds)
	case "gen":
		skipKinds, args := parseSkipFlag(os.Args[2:])
		reset, args := parseResetFlag(args)
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: fullend gen [--skip kind,...] [--reset] <specs-dir> <artifacts-dir>")
			os.Exit(2)
		}
		if reset && !confirmReset(args[1]) {
			os.Exit(0)
		}
		runGen(args[0], args[1], skipKinds, reset)
	case "gen-model":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: fullend gen-model <openapi-source> <output-dir>")
			os.Exit(2)
		}
		if err := genmodel.Generate(os.Args[2], os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "contract":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: fullend contract <specs-dir> <artifacts-dir>")
			os.Exit(2)
		}
		runContract(os.Args[2], os.Args[3])
	case "chain":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: fullend chain <operationId> <specs-dir>")
			os.Exit(2)
		}
		runChain(os.Args[2], os.Args[3])
	case "status":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: fullend status <specs-dir>")
			os.Exit(2)
		}
		runStatus(os.Args[2])
	case "map":
		runMap(os.Args[2:])
	case "history":
		runHistory(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Print(usage)
		os.Exit(2)
	}
}
