import { ApolloClient, createHttpLink, InMemoryCache } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'

const httpLink = createHttpLink({
  uri: `${import.meta.env.VITE_API_URL}/graphql`, 
})

const authLink = setContext((_, { headers }) => {
  const jwt = localStorage.getItem('jwt')

  return {
    headers: {
      ...headers,
      Authorization: jwt ? `Bearer ${jwt}` : '',
    },
  }

})

export const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
})
