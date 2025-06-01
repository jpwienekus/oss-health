import { gql } from "graphql-request";

export const GET_REPOS = gql`
  query {
    repositories {
      name
      description
      updatedAt
    }
  }
`
