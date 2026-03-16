//ff:func feature=cli type=command control=sequence
//ff:what 키워드 맵 생성 및 출력
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clari/whyso/pkg/codemap"
)

func runMap(args []string) {
	target, outputFile, force := parseMapArgs(args)

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
