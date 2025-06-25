import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from 'react'
import { generateCodeChallenge, generateCodeVerifier } from './pkce'
import { jwtDecode } from 'jwt-decode'

const AuthContext = createContext<{
  jwt: string | null
  loginWithGitHub: () => void
}>({ jwt: null, loginWithGitHub: () => {} })

type AuthProviderProps = {
  children: ReactNode
}

interface JwtPayload {
  exp: number
  sub: string
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const API_URL = import.meta.env.VITE_API_URL
  const [jwt, setJwt] = useState<string | null>(null)

  useEffect(() => {
    const token = localStorage.getItem('jwt')

    if (token) {
      const { exp } = jwtDecode<JwtPayload>(token)
      const expiryTime = exp * 1000
      const now = Date.now()
      const timeout = expiryTime - now

      if (timeout <= 0) {
        logout()
        return
      }

      setJwt(token)

      const timer = setTimeout(() => {
        logout()
      }, timeout)

      return () => clearTimeout(timer)
    }
  }, [])

  useEffect(() => {
    const handler = async (event: MessageEvent) => {
      if (
        process.env.NODE_ENV !== 'development' &&
        event.origin !== window.location.origin
      ) {
        return
      }

      const { type, code } = event.data

      if (type === 'github-oauth-code' && code) {
        const verifier = localStorage.getItem('pkce_verifier')
        const res = await fetch(`${API_URL}/auth/github/token`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ code, code_verifier: verifier }),
        })
        const data = await res.json()

        if (data.access_token) {
          localStorage.setItem('jwt', data.access_token)
          setJwt(data.access_token)
        }
      }
    }

    window.addEventListener('message', handler)
    return () => window.removeEventListener('message', handler)
  }, [])

  // GitHub login popup logic
  const loginWithGitHub = async () => {
    const verifier = generateCodeVerifier()
    const challenge = await generateCodeChallenge(verifier)
    localStorage.setItem('pkce_verifier', verifier)

    const authUrl = `${API_URL}/auth/github/login?code_challenge=${challenge}`

    window.open(authUrl, '_blank', 'popup,width=500,height=600')
  }

  const logout = () => {
    localStorage.removeItem('jwt')
    setJwt(null)
  }

  return (
    <AuthContext.Provider value={{ jwt, loginWithGitHub }}>
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}
