//ff:func feature=gen-react type=generator control=sequence
//ff:what main.tsx 파일을 생성한다 (tanstack 유무에 따라 분기)

package react

import (
	"os"
	"path/filepath"
)

func writeMainTSX(srcDir string, stmlDeps map[string]string) error {
	var src string
	if _, ok := stmlDeps["@tanstack/react-query"]; ok {
		src = `import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App'

const queryClient = new QueryClient()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </QueryClientProvider>
  </React.StrictMode>,
)
`
	} else {
		src = `import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './App'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
)
`
	}
	return os.WriteFile(filepath.Join(srcDir, "main.tsx"), []byte(src), 0644)
}
