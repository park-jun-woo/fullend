//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=scenario-check
//ff:what .hurl 파일에서 요청/응답 쌍 추출
package scenario

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var (
	reHurlRequest  = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+(?:\{\{host\}\}|https?://[^/]*)(/.+)`)
	reHurlResponse = regexp.MustCompile(`^HTTP\s+(\d+)`)
)

// parseHurlFile extracts request/response pairs from a .hurl file.
func parseHurlFile(path string) []HurlEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var entries []HurlEntry
	var current *HurlEntry

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		current, entries = processHurlLine(line, lineNum, path, current, entries)
	}

	if current != nil {
		entries = append(entries, *current)
	}

	return entries
}
