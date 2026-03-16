//ff:func feature=cli type=util control=sequence
//ff:what history YAML 한 줄을 파싱하여 FileHistory에 반영
package main

import (
	"strings"
	"time"

	"github.com/clari/whyso/pkg/history"
)

func parseHistoryYAMLLine(line string, h *history.FileHistory) {
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
