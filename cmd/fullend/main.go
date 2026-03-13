package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/clari/whyso/pkg/codemap"
	"github.com/clari/whyso/pkg/history"
	"github.com/clari/whyso/pkg/parser"
	"github.com/geul-org/fullend/internal/contract"
	"github.com/geul-org/fullend/internal/genmodel"
	"github.com/geul-org/fullend/internal/orchestrator"
	"github.com/geul-org/fullend/internal/reporter"
)

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

Skip kinds: openapi, ddl, ssac, model, stml, states, policy, scenario, func, terraform
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(2)
	}

	switch os.Args[1] {
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
		if reset {
			backendDir := filepath.Join(args[1], "backend")
			count := contract.CountPreserveFuncs(backendDir)
			if count > 0 {
				fmt.Fprintf(os.Stderr, "⚠ --reset: preserve 함수 %d개가 초기화됩니다.\n", count)
				fmt.Fprint(os.Stderr, "계속하시겠습니까? (Y/n): ")
				var answer string
				fmt.Scanln(&answer)
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer == "n" {
					fmt.Fprintln(os.Stderr, "취소됨")
					os.Exit(0)
				}
			}
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

// parseSkipFlag extracts --skip flag and returns (skipKinds, remainingArgs).
func parseSkipFlag(args []string) (map[orchestrator.SSOTKind]bool, []string) {
	skipKinds := make(map[orchestrator.SSOTKind]bool)
	var remaining []string

	for i := 0; i < len(args); i++ {
		if args[i] == "--skip" {
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "--skip requires a comma-separated list of SSOT kinds")
				os.Exit(2)
			}
			i++
			for _, s := range strings.Split(args[i], ",") {
				s = strings.TrimSpace(s)
				kind, ok := orchestrator.KindFromString(s)
				if !ok {
					fmt.Fprintf(os.Stderr, "unknown SSOT kind: %q\nvalid kinds: openapi, ddl, ssac, model, stml, states, policy, scenario, func, terraform\n", s)
					os.Exit(2)
				}
				skipKinds[kind] = true
			}
		} else {
			remaining = append(remaining, args[i])
		}
	}

	return skipKinds, remaining
}

func runValidate(specsDir string, skipKinds map[orchestrator.SSOTKind]bool) {
	detected, err := orchestrator.DetectSSOTs(specsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	report := orchestrator.Validate(specsDir, detected, skipKinds)
	reporter.Print(os.Stdout, report)

	if report.HasFailure() {
		os.Exit(1)
	}
}

func runContract(specsDir, artifactsDir string) {
	funcs, err := contract.ScanDir(artifactsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	funcs = contract.Verify(specsDir, funcs)
	reporter.PrintContract(os.Stdout, funcs)

	_, _, broken, orphan := contract.Summary(funcs)
	if broken > 0 || orphan > 0 {
		os.Exit(1)
	}
}

func runStatus(specsDir string) {
	detected, err := orchestrator.DetectSSOTs(specsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	lines := orchestrator.Status(specsDir, detected)
	orchestrator.PrintStatus(os.Stdout, lines)
}

func runChain(operationID, specsDir string) {
	links, err := orchestrator.Chain(specsDir, operationID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	rlinks := make([]reporter.ChainLink, len(links))
	for i, l := range links {
		rlinks[i] = reporter.ChainLink{Kind: l.Kind, File: l.File, Line: l.Line, Summary: l.Summary, Ownership: l.Ownership}
	}
	reporter.PrintChain(os.Stdout, operationID, rlinks)
}

func runGen(specsDir, artifactsDir string, skipKinds map[orchestrator.SSOTKind]bool, reset bool) {
	report, ok := orchestrator.Gen(specsDir, artifactsDir, skipKinds, reset)
	reporter.Print(os.Stdout, report)

	if !ok {
		os.Exit(1)
	}
}

func runMap(args []string) {
	var target, outputFile string
	var force bool

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-f", "--force":
			force = true
		default:
			if target == "" {
				target = args[i]
			}
		}
	}

	if target == "" {
		target = "."
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if outputFile == "" {
		cwd, _ := os.Getwd()
		dir := filepath.Join(cwd, ".whyso")
		os.MkdirAll(dir, 0755)
		outputFile = filepath.Join(dir, "_map.md")
	}

	if !force && !codemap.NeedsUpdate(absTarget, outputFile) {
		fmt.Fprintln(os.Stderr, "up to date")
		return
	}

	sections, err := codemap.BuildMap(absTarget)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if len(sections) == 0 {
		fmt.Println("No keywords found.")
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	codemap.FormatMap(f, sections)
	codemap.FormatMap(os.Stdout, sections)
}

func runHistory(args []string) {
	var target, format string
	var all, quiet bool
	format = "yaml"

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case "--all":
			all = true
		case "-q", "--quiet":
			quiet = true
		default:
			if target == "" {
				target = args[i]
			}
		}
	}

	if target == "" {
		fmt.Fprintln(os.Stderr, "Usage: fullend history <file|dir> [--all] [--format yaml|json]")
		os.Exit(2)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	sessionsDir, err := parser.DetectSessionsDir(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	targetInfo, err := os.Stat(absTarget)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	filter := func(relPath string) bool {
		if targetInfo.IsDir() {
			if !all {
				return false
			}
			targetRel, err := filepath.Rel(cwd, absTarget)
			if err != nil {
				return false
			}
			if targetRel == "." {
				return true
			}
			return strings.HasPrefix(relPath, targetRel+"/") || relPath == targetRel
		}
		targetRel, err := filepath.Rel(cwd, absTarget)
		if err != nil {
			return false
		}
		return relPath == targetRel
	}

	cacheDir := filepath.Join(cwd, ".whyso")
	os.MkdirAll(cacheDir, 0755)

	// incremental: check cache mtime
	var since time.Time
	if !targetInfo.IsDir() {
		targetRel, _ := filepath.Rel(cwd, absTarget)
		cachedPath := filepath.Join(cacheDir, targetRel+"."+format)
		if info, err := os.Stat(cachedPath); err == nil {
			since = info.ModTime()
		}
	}

	var histories map[string]*history.FileHistory
	if since.IsZero() {
		histories, err = history.BuildHistories(sessionsDir, cwd, filter)
	} else {
		histories, err = history.BuildHistoriesIncremental(sessionsDir, cwd, since, filter)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// write cache
	if len(histories) > 0 {
		for relPath, h := range histories {
			outPath := filepath.Join(cacheDir, relPath+"."+format)
			os.MkdirAll(filepath.Dir(outPath), 0755)
			// merge with existing
			if existing, err := readHistoryYAML(outPath); err == nil {
				h = history.Merge(existing, h)
			}
			f, err := os.Create(outPath)
			if err != nil {
				continue
			}
			formatHistory(f, h, format)
			f.Close()
		}
	}

	// stdout
	if !quiet && !targetInfo.IsDir() {
		if len(histories) > 0 {
			for _, h := range histories {
				formatHistory(os.Stdout, h, format)
				fmt.Println("---")
			}
		} else {
			// read from cache
			targetRel, _ := filepath.Rel(cwd, absTarget)
			cachedPath := filepath.Join(cacheDir, targetRel+"."+format)
			if cached, err := readHistoryYAML(cachedPath); err == nil {
				formatHistory(os.Stdout, cached, format)
				fmt.Println("---")
			}
		}
	}
}

func formatHistory(w io.Writer, h *history.FileHistory, format string) {
	if format == "json" {
		fmt.Fprintf(w, "{\n")
		fmt.Fprintf(w, "  \"apiVersion\": \"whyso/v1\",\n")
		fmt.Fprintf(w, "  \"file\": %q,\n", h.File)
		fmt.Fprintf(w, "  \"created\": %q,\n", h.Created.Format(time.RFC3339))
		fmt.Fprintf(w, "  \"history\": []\n") // simplified
		fmt.Fprintf(w, "}\n")
		return
	}
	fmt.Fprintf(w, "apiVersion: whyso/v1\n")
	fmt.Fprintf(w, "file: %s\n", h.File)
	fmt.Fprintf(w, "created: %s\n", h.Created.Format(time.RFC3339))
	fmt.Fprintf(w, "history:\n")
	for _, e := range h.History {
		fmt.Fprintf(w, "  - timestamp: %s\n", e.Timestamp.Format(time.RFC3339))
		fmt.Fprintf(w, "    session: %s\n", e.Session)
		fmt.Fprintf(w, "    user_request: %q\n", e.UserRequest)
		if e.Answer != "" {
			fmt.Fprintf(w, "    answer: %q\n", e.Answer)
		}
		fmt.Fprintf(w, "    tool: %s\n", e.Tool)
		if e.Subagent {
			fmt.Fprintf(w, "    subagent: true\n")
		}
		if len(e.Sources) == 1 {
			fmt.Fprintf(w, "    source: %s:%d\n", e.Sources[0].File, e.Sources[0].Line)
		} else if len(e.Sources) > 1 {
			fmt.Fprintf(w, "    sources:\n")
			for _, s := range e.Sources {
				fmt.Fprintf(w, "      - %s:%d\n", s.File, s.Line)
			}
		}
	}
}

func readHistoryYAML(path string) (*history.FileHistory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// simple YAML parsing — extract file and created fields, then history entries
	lines := strings.Split(string(data), "\n")
	h := &history.FileHistory{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "file: ") {
			h.File = strings.TrimPrefix(line, "file: ")
		}
		if strings.HasPrefix(line, "created: ") {
			t, err := time.Parse(time.RFC3339, strings.TrimPrefix(line, "created: "))
			if err == nil {
				h.Created = t
			}
		}
	}
	return h, nil
}
