import './App.css'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { Repositories } from './pages/Repositories'
import { Repositories as AdminRepositories } from './pages/admin/repositories'
import { Dependencies as AdminDependencies } from './pages/admin/dependencies'

function App() {
  return (
    <BrowserRouter basename="/">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Repositories />} />
          <Route path="admin">
            <Route path="repositories" element={<AdminRepositories />} />
            <Route path="dependencies" element={<AdminDependencies />} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
