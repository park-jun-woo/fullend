//ff:func feature=cli type=formatter control=iteration dimension=1
//ff:what history 맵의 각 항목을 stdout에 출력
package main

import (
	"fmt"
	"os"

	"github.com/clari/whyso/pkg/history"
)

// printHistoryEntries prints all history entries to stdout.
func printHistoryEntries(histories map[string]*history.FileHistory, format string) {
	for _, h := range histories {
		formatHistory(os.Stdout, h, format)
		fmt.Println("---")
	}
}
