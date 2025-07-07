import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { AuthProvider } from './auth/AuthContext.tsx'
import { ApolloProvider } from '@apollo/client'
import { client } from '@/apolloClient.ts'
import { Toaster } from 'sonner'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ApolloProvider client={client}>
      <AuthProvider>
        <App />
        <Toaster position='top-center'/>
      </AuthProvider>
    </ApolloProvider>
  </StrictMode>
)
