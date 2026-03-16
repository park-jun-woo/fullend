//ff:func feature=reporter type=formatter control=iteration dimension=1
//ff:what 에러 목록과 제안을 들여쓰기하여 출력한다
package reporter

import (
	"fmt"
	"io"
)

// printErrors writes indented error messages with optional suggestions.
func printErrors(w io.Writer, errors []string, suggestions []string) {
	for i, e := range errors {
		fmt.Fprintf(w, "    %s\n", e)
		if i < len(suggestions) && suggestions[i] != "" {
			fmt.Fprintf(w, "      → 제안: %s\n", suggestions[i])
		}
	}
}
