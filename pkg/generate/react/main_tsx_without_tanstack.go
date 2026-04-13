//ff:func feature=gen-react type=util control=sequence
//ff:what tanstack 미포함 main.tsx 템플릿 상수

package react

const mainTSXWithoutTanstack = `import React from 'react'
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
