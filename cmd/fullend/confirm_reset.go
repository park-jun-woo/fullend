//ff:func feature=cli type=util control=sequence
//ff:what --reset 플래그 사용 시 preserve 함수 초기화 확인
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/fullend/internal/contract"
)

// confirmReset checks for preserve funcs and prompts user. Returns false if user cancels.
func confirmReset(artifactsDir string) bool {
	backendDir := filepath.Join(artifactsDir, "backend")
	count := contract.CountPreserveFuncs(backendDir)
	if count == 0 {
		return true
	}
	fmt.Fprintf(os.Stderr, "⚠ --reset: preserve 함수 %d개가 초기화됩니다.\n", count)
	fmt.Fprint(os.Stderr, "계속하시겠습니까? (Y/n): ")
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "n" {
		fmt.Fprintln(os.Stderr, "취소됨")
		return false
	}
	return true
}
