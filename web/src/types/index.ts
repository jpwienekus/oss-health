// export interface GitHubRepository {
//   id: number
//   githubId: number
//   name: string
//   description: string
//   stars: number
//   private: boolean
//   forks: number
//   watchers: number
//   score: number
//   vulnerabilities: number
//   dependencies: number
//   lastScannedAt: string
//   updatedAt: string
// }

export interface Dependency {
    id: number
    name: string
    ecosystem: string
    githubUrlResolved: boolean
    githubUrlResolveFailed: boolean
    githubUrlResolveFailedReason: string
}

export interface DependencyEdge {
  node: Dependency,
  cursor: number
}

export interface PageInfo {
    hasNextPage: boolean
    hasPreviousPage: boolean
    startCursor?: number
    endCursor?: number
}

export interface DependencyConnection {
  edges: DependencyEdge[]
  pageInfo: PageInfo
}

// export interface Pagination {
//   limit: number,
//   after?: number
// }
//
// export interface DependencyFilter {
//   name: string
//   ecosystem: string
//   githubUrlResolveFailed?: boolean
// }
//
// export interface DependencySort {
//
// }

    // field: DependencySortField = DependencySortField.ID
    // direction: SortDirection = SortDirection.ASC
/*
 
export const GET_DEPENDENCIES = gql`
  query Dependencies($pagination: DependencyPaginationInput, $filter: DependencyFilter, $sort: SortInput) {
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
 *
 */
