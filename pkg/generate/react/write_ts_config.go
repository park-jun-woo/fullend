//ff:func feature=gen-react type=generator control=sequence
//ff:what tsconfig.json 파일을 생성한다

package react

import (
	"os"
	"path/filepath"
)

func writeTSConfig(dir string) error {
	src := `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "isolatedModules": true,
    "moduleDetection": "force",
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "noUncheckedSideEffectImports": true
  },
  "include": ["src"]
}
`
	return os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(src), 0644)
}
