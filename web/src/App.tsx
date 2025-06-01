import './App.css'
import { Button } from './components/ui/button'
import { useGitHubPopupLogin } from './auth/useGitHubPopupLogin'
import { getClient } from './graphql/client'
import { SYNC_REPOS } from './graphql/mutations'
import { GET_REPOS } from './graphql/queries'
import type { Repository } from './types'

function App() {
  const { jwt, login } = useGitHubPopupLogin()
  
  const handleSync = async() => {
    if (!jwt) {
      return
    }

    const client = getClient(jwt)
    await client.request(SYNC_REPOS)
    const data = await client.request<{ repositories: Repository[]}>(GET_REPOS)
    console.log(3343, data)
  }


  return (
    <div className='flex min-h-svh flex-col items-center justify-center'>
      {!jwt ? (
        <Button onClick={login}>Click Me</Button>
      ) : (
        <>
          <h1>You're logged in</h1>
          <Button onClick={handleSync}>Sync</Button>
        </>
      )}
    </div>
  )
}

export default App
