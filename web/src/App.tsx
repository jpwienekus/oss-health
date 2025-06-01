import './App.css'
import { Button } from './components/ui/button'
import { useGitHubPopupLogin } from './auth/useGitHubPopupLogin'

function App() {
  const { jwt, login } = useGitHubPopupLogin()
  return (
    <div className='flex min-h-svh flex-col items-center justify-center'>
      {!false ? (
        <Button onClick={login}>Click Me</Button>
      ) : (
        <>
          <h1>You're logged in</h1>
        </>
      )}
    </div>
  )
}

export default App
