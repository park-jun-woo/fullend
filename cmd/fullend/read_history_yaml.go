//ff:func feature=cli type=util control=iteration dimension=1
//ff:what 캐시된 history YAML 파일 읽기
package main

import (
	"os"
	"strings"

	"github.com/clari/whyso/pkg/history"
)

func readHistoryYAML(path string) (*history.FileHistory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// simple YAML parsing — extract file and created fields
	lines := strings.Split(string(data), "\n")
	h := &history.FileHistory{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parseHistoryYAMLLine(line, h)
	}
	return h, nil
}
