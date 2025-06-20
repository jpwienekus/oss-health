import { gql } from 'graphql-request'

export const SAVE_SELECTED_REPOSITORIES = gql`
  mutation SaveSelectedRepositories($selectedGithubRepositoryIds: [Int!]!) {
    saveSelectedRepositories(
      selectedGithubRepositoryIds: $selectedGithubRepositoryIds
    ) {
      name
      description
      githubId
      updatedAt
      private
      stars
      watchers
      forks
    }
  }
`
