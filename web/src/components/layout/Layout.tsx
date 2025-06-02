import { useEffect, useState } from "react"
import { Header } from "./Header"
import { Outlet } from "react-router-dom"

export const Layout = () => {
  const THEME_KEY = 'theme'
  const [isDarkMode, setIsDarkMode] = useState(() => {
    if (typeof window === 'undefined') {
      return false
    }

    const savedTheme = localStorage.getItem(THEME_KEY)
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches

    return savedTheme === 'dark' || (!savedTheme && prefersDark)
  })

  const toggleDarkMode = () => {
    setIsDarkMode(!isDarkMode)
  }

  useEffect(() => {
    const root = window.document.documentElement

    if (isDarkMode) {
      root.classList.add('dark')
      localStorage.setItem(THEME_KEY, 'dark')
    } else {
      root.classList.remove('dark')
      localStorage.setItem(THEME_KEY, 'light')
    }
  }, [isDarkMode])

  return (
    <div className="flex flex-col min-h-screen">
      <Header isDarkMode={isDarkMode} toggleDarkMode={toggleDarkMode} />

      <main className="container mx-auto px-4 py-8 flex-grow">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-8 xl:col-span-9">
            <Outlet />
          </div>

          {/*<div className="lg:col-span-4 xl:col-span-3">
            <Sidebar />
          </div> */}
        </div>

      </main>

      {/*<Footer /> */}
    </div>
  )
}
