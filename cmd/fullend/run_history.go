//ff:func feature=cli type=command control=sequence
//ff:what 파일 변경 히스토리 조회 및 출력
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/clari/whyso/pkg/history"
	"github.com/clari/whyso/pkg/parser"
)

func runHistory(args []string) {
	target, format, all, quiet := parseHistoryArgs(args)

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

	filter := buildHistoryFilter(absTarget, cwd, targetInfo, all)

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
		writeHistoryCache(histories, cacheDir, format)
	}

	// stdout
	if !quiet && !targetInfo.IsDir() {
		printHistoryStdout(histories, absTarget, cwd, cacheDir, format)
	}
}
