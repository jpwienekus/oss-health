import { gql } from 'graphql-request'

export const GET_REPOSITORIES_FROM_GITHUB = gql`
  query {
    githubRepositories {
      name
      description
      githubId
      private
      lastScannedAt
      updatedAt
      stars
      watchers
      forks
    }
  }
`
export const GET_REPOSITORIES = gql`
  query {
    repositories {
      id
      name
      description
      githubId
      private
      score
      vulnerabilities
      dependencies
      lastScannedAt
      updatedAt
      stars
      watchers
      forks
    }
  }
`

export const GET_USERNAME = gql`
  query {
    username
  }
`

export const GET_DEPENDENCIES = gql`
  query Dependencies($pagination: PaginationInput!, $filter: DependencyFilter!, $sort: DependencySortInput!) {
    dependencies(pagination: $pagination, filter: $filter, sort: $sort) {
      edges {
        node {
          id
          name
          ecosystem
        }
        cursor
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
    }
  }
`
