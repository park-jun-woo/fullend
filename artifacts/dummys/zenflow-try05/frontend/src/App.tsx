import { Routes, Route } from 'react-router-dom'
import Login from './pages/login'
import Workflows from './pages/workflows'
import WorkflowDetail from './pages/workflow-detail'

export default function App() {
  return (
    <Routes>
      <Route path="/auth/login" element={<Login />} />
      <Route path="/workflows" element={<Workflows />} />
      <Route path="/workflows/:id" element={<WorkflowDetail />} />
    </Routes>
  )
}
