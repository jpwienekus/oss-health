import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  DateTime: { input: any; output: any; }
};

export type CronInfo = {
  __typename?: 'CronInfo';
  day: Scalars['String']['output'];
  hour: Scalars['String']['output'];
  total: Scalars['Int']['output'];
};

export type DependencyFilter = {
  name?: Scalars['String']['input'];
  statuses: Array<Scalars['String']['input']>;
};

export type DependencyPaginatedResponse = {
  __typename?: 'DependencyPaginatedResponse';
  completed: Scalars['Int']['output'];
  dependencies: Array<DependencyType>;
  failed: Scalars['Int']['output'];
  pending: Scalars['Int']['output'];
  totalPages: Scalars['Int']['output'];
};

export type DependencySortField =
  | 'ECOSYSTEM'
  | 'ERROR_MESSAGE'
  | 'ID'
  | 'NAME'
  | 'SCANNED_AT'
  | 'SCAN_STATUS';

export type DependencySortInput = {
  direction?: SortDirection;
  field?: DependencySortField;
};

export type DependencyType = {
  __typename?: 'DependencyType';
  ecosystem: Scalars['String']['output'];
  errorMessage?: Maybe<Scalars['String']['output']>;
  id: Scalars['Int']['output'];
  name: Scalars['String']['output'];
  repositoryUrl?: Maybe<Scalars['String']['output']>;
  scanStatus: Scalars['String']['output'];
  scannedAt?: Maybe<Scalars['DateTime']['output']>;
};

export type GitHubRepository = {
  __typename?: 'GitHubRepository';
  dependencies?: Maybe<Scalars['Int']['output']>;
  description?: Maybe<Scalars['String']['output']>;
  forks: Scalars['Int']['output'];
  githubId: Scalars['Int']['output'];
  id?: Maybe<Scalars['Int']['output']>;
  lastScannedAt?: Maybe<Scalars['DateTime']['output']>;
  name: Scalars['String']['output'];
  private: Scalars['Boolean']['output'];
  score?: Maybe<Scalars['Int']['output']>;
  stars: Scalars['Int']['output'];
  updatedAt?: Maybe<Scalars['DateTime']['output']>;
  url: Scalars['String']['output'];
  vulnerabilities?: Maybe<Scalars['Int']['output']>;
  watchers: Scalars['Int']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  saveSelectedRepositories: Array<GitHubRepository>;
};


export type MutationSaveSelectedRepositoriesArgs = {
  selectedGithubRepositoryIds: Array<Scalars['Int']['input']>;
};

export type PaginationInput = {
  page?: Scalars['Int']['input'];
  pageSize?: Scalars['Int']['input'];
};

export type Query = {
  __typename?: 'Query';
  dependencies: DependencyPaginatedResponse;
  getCronInfo: Array<CronInfo>;
  githubRepositories: Array<GitHubRepository>;
  repositories: Array<GitHubRepository>;
  username: Scalars['String']['output'];
};


export type QueryDependenciesArgs = {
  filter: DependencyFilter;
  pagination: PaginationInput;
  sort: DependencySortInput;
};

export type SortDirection =
  | 'ASC'
  | 'DESC';

export type SaveSelectedRepositoriesMutationVariables = Exact<{
  selectedGithubRepositoryIds: Array<Scalars['Int']['input']> | Scalars['Int']['input'];
}>;


export type SaveSelectedRepositoriesMutation = { __typename?: 'Mutation', saveSelectedRepositories: Array<{ __typename?: 'GitHubRepository', id?: number | null, name: string, description?: string | null, githubId: number, stars: number, watchers: number, forks: number, private: boolean, score?: number | null, vulnerabilities?: number | null, dependencies?: number | null, url: string, lastScannedAt?: any | null, updatedAt?: any | null }> };

export type GetCronInfoQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCronInfoQuery = { __typename?: 'Query', getCronInfo: Array<{ __typename?: 'CronInfo', day: string, hour: string, total: number }> };

export type GetDependenciesQueryVariables = Exact<{
  pagination: PaginationInput;
  filter: DependencyFilter;
  sort: DependencySortInput;
}>;


export type GetDependenciesQuery = { __typename?: 'Query', dependencies: { __typename?: 'DependencyPaginatedResponse', totalPages: number, completed: number, pending: number, failed: number, dependencies: Array<{ __typename?: 'DependencyType', id: number, name: string, ecosystem: string, scannedAt?: any | null, scanStatus: string, errorMessage?: string | null, repositoryUrl?: string | null }> } };

export type GithubRepositoriesQueryVariables = Exact<{ [key: string]: never; }>;


export type GithubRepositoriesQuery = { __typename?: 'Query', githubRepositories: Array<{ __typename?: 'GitHubRepository', id?: number | null, name: string, description?: string | null, githubId: number, stars: number, watchers: number, forks: number, private: boolean, score?: number | null, vulnerabilities?: number | null, dependencies?: number | null, url: string, lastScannedAt?: any | null, updatedAt?: any | null }> };

export type GetRepositoriesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetRepositoriesQuery = { __typename?: 'Query', repositories: Array<{ __typename?: 'GitHubRepository', id?: number | null, name: string, description?: string | null, githubId: number, private: boolean, score?: number | null, vulnerabilities?: number | null, dependencies?: number | null, lastScannedAt?: any | null, updatedAt?: any | null, stars: number, watchers: number, forks: number, url: string }> };

export type GetUserQueryVariables = Exact<{ [key: string]: never; }>;


export type GetUserQuery = { __typename?: 'Query', username: string };


export const SaveSelectedRepositoriesDocument = gql`
    mutation SaveSelectedRepositories($selectedGithubRepositoryIds: [Int!]!) {
  saveSelectedRepositories(
    selectedGithubRepositoryIds: $selectedGithubRepositoryIds
  ) {
    id
    name
    description
    githubId
    stars
    watchers
    forks
    private
    score
    vulnerabilities
    dependencies
    url
    lastScannedAt
    updatedAt
  }
}
    `;
export type SaveSelectedRepositoriesMutationFn = Apollo.MutationFunction<SaveSelectedRepositoriesMutation, SaveSelectedRepositoriesMutationVariables>;

/**
 * __useSaveSelectedRepositoriesMutation__
 *
 * To run a mutation, you first call `useSaveSelectedRepositoriesMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSaveSelectedRepositoriesMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [saveSelectedRepositoriesMutation, { data, loading, error }] = useSaveSelectedRepositoriesMutation({
 *   variables: {
 *      selectedGithubRepositoryIds: // value for 'selectedGithubRepositoryIds'
 *   },
 * });
 */
export function useSaveSelectedRepositoriesMutation(baseOptions?: Apollo.MutationHookOptions<SaveSelectedRepositoriesMutation, SaveSelectedRepositoriesMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SaveSelectedRepositoriesMutation, SaveSelectedRepositoriesMutationVariables>(SaveSelectedRepositoriesDocument, options);
      }
export type SaveSelectedRepositoriesMutationHookResult = ReturnType<typeof useSaveSelectedRepositoriesMutation>;
export type SaveSelectedRepositoriesMutationResult = Apollo.MutationResult<SaveSelectedRepositoriesMutation>;
export type SaveSelectedRepositoriesMutationOptions = Apollo.BaseMutationOptions<SaveSelectedRepositoriesMutation, SaveSelectedRepositoriesMutationVariables>;
export const GetCronInfoDocument = gql`
    query GetCronInfo {
  getCronInfo {
    day
    hour
    total
  }
}
    `;

/**
 * __useGetCronInfoQuery__
 *
 * To run a query within a React component, call `useGetCronInfoQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCronInfoQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCronInfoQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetCronInfoQuery(baseOptions?: Apollo.QueryHookOptions<GetCronInfoQuery, GetCronInfoQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCronInfoQuery, GetCronInfoQueryVariables>(GetCronInfoDocument, options);
      }
export function useGetCronInfoLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCronInfoQuery, GetCronInfoQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCronInfoQuery, GetCronInfoQueryVariables>(GetCronInfoDocument, options);
        }
export function useGetCronInfoSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCronInfoQuery, GetCronInfoQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCronInfoQuery, GetCronInfoQueryVariables>(GetCronInfoDocument, options);
        }
export type GetCronInfoQueryHookResult = ReturnType<typeof useGetCronInfoQuery>;
export type GetCronInfoLazyQueryHookResult = ReturnType<typeof useGetCronInfoLazyQuery>;
export type GetCronInfoSuspenseQueryHookResult = ReturnType<typeof useGetCronInfoSuspenseQuery>;
export type GetCronInfoQueryResult = Apollo.QueryResult<GetCronInfoQuery, GetCronInfoQueryVariables>;
export const GetDependenciesDocument = gql`
    query GetDependencies($pagination: PaginationInput!, $filter: DependencyFilter!, $sort: DependencySortInput!) {
  dependencies(pagination: $pagination, filter: $filter, sort: $sort) {
    dependencies {
      id
      name
      ecosystem
      scannedAt
      scanStatus
      errorMessage
      repositoryUrl
    }
    totalPages
    completed
    pending
    failed
  }
}
    `;

/**
 * __useGetDependenciesQuery__
 *
 * To run a query within a React component, call `useGetDependenciesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDependenciesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDependenciesQuery({
 *   variables: {
 *      pagination: // value for 'pagination'
 *      filter: // value for 'filter'
 *      sort: // value for 'sort'
 *   },
 * });
 */
export function useGetDependenciesQuery(baseOptions: Apollo.QueryHookOptions<GetDependenciesQuery, GetDependenciesQueryVariables> & ({ variables: GetDependenciesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDependenciesQuery, GetDependenciesQueryVariables>(GetDependenciesDocument, options);
      }
export function useGetDependenciesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDependenciesQuery, GetDependenciesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDependenciesQuery, GetDependenciesQueryVariables>(GetDependenciesDocument, options);
        }
export function useGetDependenciesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDependenciesQuery, GetDependenciesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetDependenciesQuery, GetDependenciesQueryVariables>(GetDependenciesDocument, options);
        }
export type GetDependenciesQueryHookResult = ReturnType<typeof useGetDependenciesQuery>;
export type GetDependenciesLazyQueryHookResult = ReturnType<typeof useGetDependenciesLazyQuery>;
export type GetDependenciesSuspenseQueryHookResult = ReturnType<typeof useGetDependenciesSuspenseQuery>;
export type GetDependenciesQueryResult = Apollo.QueryResult<GetDependenciesQuery, GetDependenciesQueryVariables>;
export const GithubRepositoriesDocument = gql`
    query GithubRepositories {
  githubRepositories {
    id
    name
    description
    githubId
    stars
    watchers
    forks
    private
    score
    vulnerabilities
    dependencies
    url
    lastScannedAt
    updatedAt
  }
}
    `;

/**
 * __useGithubRepositoriesQuery__
 *
 * To run a query within a React component, call `useGithubRepositoriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGithubRepositoriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGithubRepositoriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGithubRepositoriesQuery(baseOptions?: Apollo.QueryHookOptions<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>(GithubRepositoriesDocument, options);
      }
export function useGithubRepositoriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>(GithubRepositoriesDocument, options);
        }
export function useGithubRepositoriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>(GithubRepositoriesDocument, options);
        }
export type GithubRepositoriesQueryHookResult = ReturnType<typeof useGithubRepositoriesQuery>;
export type GithubRepositoriesLazyQueryHookResult = ReturnType<typeof useGithubRepositoriesLazyQuery>;
export type GithubRepositoriesSuspenseQueryHookResult = ReturnType<typeof useGithubRepositoriesSuspenseQuery>;
export type GithubRepositoriesQueryResult = Apollo.QueryResult<GithubRepositoriesQuery, GithubRepositoriesQueryVariables>;
export const GetRepositoriesDocument = gql`
    query GetRepositories {
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
    url
  }
}
    `;

/**
 * __useGetRepositoriesQuery__
 *
 * To run a query within a React component, call `useGetRepositoriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRepositoriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRepositoriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetRepositoriesQuery(baseOptions?: Apollo.QueryHookOptions<GetRepositoriesQuery, GetRepositoriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetRepositoriesQuery, GetRepositoriesQueryVariables>(GetRepositoriesDocument, options);
      }
export function useGetRepositoriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetRepositoriesQuery, GetRepositoriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetRepositoriesQuery, GetRepositoriesQueryVariables>(GetRepositoriesDocument, options);
        }
export function useGetRepositoriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetRepositoriesQuery, GetRepositoriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetRepositoriesQuery, GetRepositoriesQueryVariables>(GetRepositoriesDocument, options);
        }
export type GetRepositoriesQueryHookResult = ReturnType<typeof useGetRepositoriesQuery>;
export type GetRepositoriesLazyQueryHookResult = ReturnType<typeof useGetRepositoriesLazyQuery>;
export type GetRepositoriesSuspenseQueryHookResult = ReturnType<typeof useGetRepositoriesSuspenseQuery>;
export type GetRepositoriesQueryResult = Apollo.QueryResult<GetRepositoriesQuery, GetRepositoriesQueryVariables>;
export const GetUserDocument = gql`
    query GetUser {
  username
}
    `;

/**
 * __useGetUserQuery__
 *
 * To run a query within a React component, call `useGetUserQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetUserQuery(baseOptions?: Apollo.QueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
      }
export function useGetUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export function useGetUserSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export type GetUserQueryHookResult = ReturnType<typeof useGetUserQuery>;
export type GetUserLazyQueryHookResult = ReturnType<typeof useGetUserLazyQuery>;
export type GetUserSuspenseQueryHookResult = ReturnType<typeof useGetUserSuspenseQuery>;
export type GetUserQueryResult = Apollo.QueryResult<GetUserQuery, GetUserQueryVariables>;