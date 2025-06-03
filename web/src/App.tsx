import './App.css'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/layout/Layout'
import { Repositories } from './pages/Repositories'

function App() {
  return (
    <BrowserRouter basename="/">
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Repositories />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
