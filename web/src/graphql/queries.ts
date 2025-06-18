import { gql } from 'graphql-request'

export const GET_REPOSITORIES_FROM_GITHUB = gql`
  query {
    githubRepositories {
      name
      description
      githubId
      private
      scannedDate
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
      scannedDate
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

export const MANUAL_SCAN_DEBUG = gql`
  query ManualScanDebug($repositoryId: Int!) {
    manualScanDebug(repositoryId: $repositoryId) {
      id
      name
      description
      githubId
      private
      score
      vulnerabilities
      dependencies
      scannedDate
      updatedAt
      stars
      watchers
      forks
    }
  }
`
