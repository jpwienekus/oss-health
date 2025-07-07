import { useAuth } from '@/auth/AuthContext'
import { Button } from '@/components/ui/button'
import { useGetUserQuery } from '@/generated/graphql'
import { Shield, Github, BarChart3, Users, Database, Settings } from 'lucide-react'
import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { toast } from 'sonner'

import { Navbar } from './nav-bar'

export const Header = () => {
  const { jwt, loginWithGitHub } = useAuth()

  const { data, loading, error } = useGetUserQuery({
    skip: !jwt,
    notifyOnNetworkStatusChange: true
  })

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error("Could not log in", {
      description: error.message,
    })

  }, [error])

  const isAdminPath = location.pathname.startsWith('/admin')
  const isAuthenticated = !loading && jwt

  return (
    <Navbar />
  )
}
