//ff:func feature=gen-react type=generator control=sequence
//ff:what main.tsx 파일을 생성한다 (tanstack 유무에 따라 분기)

package react

import (
	"os"
	"path/filepath"
)

func writeMainTSX(srcDir string, stmlDeps map[string]string) error {
	_, useTanstack := stmlDeps["@tanstack/react-query"]
	src := mainTSXSource(useTanstack)
	return os.WriteFile(filepath.Join(srcDir, "main.tsx"), []byte(src), 0644)
}
