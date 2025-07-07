import './App.css'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { Repositories } from './pages/Repositories'
import { Repositories as AdminRepositories } from './pages/admin/repositories'
import { Dependencies as AdminDependencies } from './pages/admin/dependencies'
import { ThemeProvider } from './components/theme-provider'

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
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
    </ThemeProvider>
  )
}

export default App
