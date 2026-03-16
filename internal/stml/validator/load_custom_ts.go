//ff:func feature=stml-validate type=parser control=iteration dimension=1
//ff:what custom.ts 파일에서 export 함수명을 추출
package validator

import (
	"bufio"
	"os"
	"regexp"
)

// LoadCustomTS parses a custom.ts file and extracts exported function names.
func LoadCustomTS(path string) (*CustomSymbol, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &CustomSymbol{Functions: make(map[string]bool)}, nil
		}
		return nil, err
	}
	defer f.Close()

	cs := &CustomSymbol{Functions: make(map[string]bool)}
	re := regexp.MustCompile(`export\s+function\s+(\w+)`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if m := re.FindStringSubmatch(scanner.Text()); m != nil {
			cs.Functions[m[1]] = true
		}
	}
	return cs, scanner.Err()
}
