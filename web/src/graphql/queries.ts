import { gql } from 'graphql-request'

export const GET_REPOS = gql`
  query {
    repositories {
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

export const GET_USERNAME = gql`
  query {
    username
  }
`
