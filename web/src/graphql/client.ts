import { GraphQLClient } from 'graphql-request'

export const getClient = (jwt: string) => {
  const url = 'http://localhost:8000'
  return new GraphQLClient(`${url}/graphql`, {
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  })
}
