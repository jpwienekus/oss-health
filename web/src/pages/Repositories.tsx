import { useGitHubPopupLogin } from "@/auth/useGitHubPopupLogin"
import { RequestGithubLogin } from "@/components/repositories/RequestGithubLogin"
// import { getClient } from "@/graphql/client"
// import { SYNC_REPOS } from "@/graphql/mutations"
// import { GET_REPOS } from "@/graphql/queries"
// import type { Repository } from "@/types"

export const Repositories = () => {
  const { jwt } = useGitHubPopupLogin()
  
  // const handleSync = async() => {
  //   if (!jwt) {
  //     return
  //   }
  //
  //   const client = getClient(jwt)
  //   await client.request(SYNC_REPOS)
  //   const data = await client.request<{ repositories: Repository[]}>(GET_REPOS)
  //   console.log(3343, data)
  // }

  return (
    <div className='flex min-h-svh flex-col items-center justify-center'>
      {!jwt ? (
        <RequestGithubLogin />
      ) : (
        <span>logged in</span>
      )}
    </div>
  )
}
