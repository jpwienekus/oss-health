import { useAuth } from "@/auth/AuthContext"
import { Button } from "@/components/ui/button"
import { getClient } from "@/graphql/client"
import { GET_USERNAME } from "@/graphql/queries"
import { Shield, Github } from "lucide-react"
import { useEffect, useState } from "react"

export const Header = () => {
  const { jwt, loginWithGitHub } = useAuth()
  const [username, setUsername] = useState<string>('')


  useEffect(() => {
    const fetchUsername = async () => {
      if (!jwt) {
        return
      }
      const client = getClient(jwt)
      const response = await client.request<{ username: string }>(GET_USERNAME)
      setUsername(response.username)
    }
    fetchUsername()
  }, [jwt])


  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center py-4">
          <div className="flex items-center space-x-3">
            <Shield className="h-8 w-8 text-gray-900" />
            <div>
              <h1 className="text-2xl font-bold text-gray-900">OSS Health</h1>
              <p className="text-sm text-gray-500">Dependency Security & Health Monitoring</p>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            {jwt ? (
              <span>
                <Github className="w-4 h-4 inline mr-2" />
                {username}
              </span>
            ) : (
              <Button onClick={loginWithGitHub} className="px-4 py-2">
                <Github className="w-4 h-4 inline mr-2" />
                Log in with GitHub
              </Button>
            )
            }
          </div>
        </div>
      </div>

    </header>
  )
}
