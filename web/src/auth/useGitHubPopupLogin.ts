import { useEffect, useState } from "react";
import { generateCodeChallenge, generateCodeVerifier } from "./pkce";
import { jwtDecode } from "jwt-decode"

interface JwtPayload {
  exp: number
  sub: string
}

export function useGitHubPopupLogin() {
  const [jwt, setJwt] = useState<string | null>(localStorage.getItem("jwt"))

  useEffect(() => {
    if (!jwt) {
      return
    }

    const { exp } = jwtDecode<JwtPayload>(jwt)
    const expiryTime = exp * 1000
    const now = Date.now()
    const timeout = expiryTime - now

    if (timeout <= 0) {
      logout()
      return
    }

    const timer = setTimeout(() => {
      logout()
    }, timeout)

    return () => clearTimeout(timer)
  }, [jwt])

  const logout = () => {
    localStorage.removeItem("jwt")
    setJwt(null)
  }

  const login = async () => {
    const codeVerifier = generateCodeVerifier()
    const codeChallenge = await generateCodeChallenge(codeVerifier)
    localStorage.setItem("pkce_verifier", codeVerifier)

    const url = "http://localhost:8000"
    const popup = window.open(`${url}/auth/github/login?code_challenge=${codeChallenge}`, "github-oauth", "width=600,height=700")
    const check = setInterval(() => {
      if (!popup || popup.closed) {
        clearInterval(check)
      }
    }, 500)

    const messageHandler = async (event: MessageEvent) => {
      if (process.env.NODE_ENV !== 'development' && event.origin !== window.location.origin) {
        return
      }

      const { type, code } = event.data

      if (type === "github-oauth-code" && code) {
        const verifier = localStorage.getItem("pkce_verifier")
        const response = await fetch(`${url}/auth/github/token`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ code, code_verifier: verifier })
        })
        const data = await response.json()
        setJwt(data.access_token)
        localStorage.setItem("jwt", data.access_token)
      }
    }

    window.addEventListener("message", messageHandler, { once: true })
  }

  return { jwt, login }
}
