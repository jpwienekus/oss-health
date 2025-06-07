import { gql } from 'graphql-request'

export const GET_REPOSITORIES_FROM_GITHUB = gql`
  query {
    githubRepositories {
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
export const GET_REPOSITORIES = gql`
  query {
    repositories {
      name
      description
      githubId
      stars
      watchers
      updatedAt
      private
      forks
      score
      vulnerabilities
      dependencies
    }
  }
`

export const GET_USERNAME = gql`
  query {
    username
  }
`
