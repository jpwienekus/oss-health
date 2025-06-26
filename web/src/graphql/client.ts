import { GraphQLClient } from 'graphql-request'

export const getClient = (jwt: string) => {
  return new GraphQLClient(`${import.meta.env.VITE_API_URL}/graphql`, {
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
  })
}
