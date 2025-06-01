import { gql } from "graphql-request";

export const SYNC_REPOS = gql`
  mutation SyncRepositories {
    syncRepositories
  }
`
