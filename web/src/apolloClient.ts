import {
  ApolloClient,
  ApolloLink,
  createHttpLink,
  InMemoryCache,
  Observable,
  from,
} from '@apollo/client'
import { setContext } from '@apollo/client/link/context'
import { loadErrorMessages, loadDevMessages } from '@apollo/client/dev'
loadDevMessages()
loadErrorMessages()
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

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

const delayLink = new ApolloLink((operation, forward) => {
  return new Observable((observer) => {
    delay(1500).then(() => {
      forward(operation).subscribe({
        next: observer.next.bind(observer),
        error: observer.error.bind(observer),
        complete: observer.complete.bind(observer),
      })
    })
  })
})

const simulateSlowLoad = false
const links = [authLink]

if (simulateSlowLoad) {
  links.push(delayLink)
}
links.push(httpLink)

export const client = new ApolloClient({
  link: from(links),
  cache: new InMemoryCache(),
})
