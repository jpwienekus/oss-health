import { useGitHubPopupLogin } from "@/auth/useGitHubPopupLogin"
import { Button } from "@/components/ui/button"
import { IconBrandGithub } from "@tabler/icons-react"


export const RequestGithubLogin = () => {
  const { login } = useGitHubPopupLogin()

  return (
    <>
      <h1 className="mb-5">Log in with GitHub to view your repositories</h1>
      <Button onClick={login}>
        <IconBrandGithub />
        Log in with Github
      </Button>
    </>
  )

}
