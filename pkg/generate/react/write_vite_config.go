//ff:func feature=gen-react type=generator control=sequence
//ff:what vite.config.ts 파일을 생성한다

package react

import (
	"os"
	"path/filepath"
)

func writeViteConfig(dir string) error {
	src := `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
`
	return os.WriteFile(filepath.Join(dir, "vite.config.ts"), []byte(src), 0644)
}
