import { createContext, useContext, type ReactNode } from 'react'
import { useGitHubPopupLogin } from './useGitHubPopupLogin'

const AuthContext = createContext<ReturnType<typeof useGitHubPopupLogin> | null>(null)

type AuthProviderProps = {
  children: ReactNode
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const auth = useGitHubPopupLogin()
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}
