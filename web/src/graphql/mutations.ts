import { gql } from 'graphql-request'

export const SYNC_REPOS = gql`
  mutation SyncRepositories {
    syncRepositories {
      repositories {
        name
        description
        updatedAt
        url
        openIssues
        score
      }
      syncDate
    }
  }
`
export const SAVE_SELECTED_REPOSITORIES = gql`
  mutation SaveSelectedRepositories($selectedGithubRepositoryIds: [Int!]!) {
    saveSelectedRepositories(
      selectedGithubRepositoryIds: $selectedGithubRepositoryIds
    ) {
      name
      description
      githubId
      stars
      watchers
      updatedAt
      private
      forks
    }
  }
`
