//ff:func feature=orchestrator type=command control=sequence
//ff:what RunExec runs an external command with a timeout, skipping if tool not installed.

package orchestrator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	// Ensure ~/go/bin is in PATH so that go-installed tools are found.
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	goBin := filepath.Join(home, "go", "bin")
	path := os.Getenv("PATH")
	if !strings.Contains(path, goBin) {
		os.Setenv("PATH", goBin+string(os.PathListSeparator)+path)
	}
}

const execTimeout = 30 * time.Second

// RunExec runs an external command with a timeout.
// If the command is not found in PATH, it returns Skipped=true.
func RunExec(name string, args ...string) ExecResult {
	_, err := exec.LookPath(name)
	if err != nil {
		return ExecResult{Skipped: true}
	}

	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return ExecResult{
			Err:    fmt.Errorf("%s failed: %w", name, err),
			Stderr: stderr.String(),
		}
	}

	return ExecResult{}
}
