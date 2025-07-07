import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { LogIn, Package } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/auth/AuthContext'
import { ImportReposDialog } from '@/components/repositories/ImportReposDialog'
import { RepositoryOverview } from '@/components/repositories/RepositoryOverview'
import { useGetRepositoriesQuery, useSaveSelectedRepositoriesMutation, type GitHubRepository } from '@/generated/graphql'
import SwirlingEffectSpinner from '@/components/customized/spinner/spinner-06'
import { toast } from 'sonner'

export const Repositories = () => {
  const { jwt } = useAuth()
  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState<'health' | 'stars' | 'updated' | 'name'>(
    'health',
  )
  const [repositories, setRepositories] = useState<GitHubRepository[]>([])

  const { data, loading, error } = useGetRepositoriesQuery()

  const filteredAndSortedRepositories = repositories
    .filter(
      (repo) =>
        repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        repo.description?.toLowerCase().includes(searchTerm.toLowerCase()),
    )
    .sort((a, b) => {
      switch (sortBy) {
        case 'health':
          return (b.score ?? 0) - (a.score ?? 0)
        case 'stars':
          return b.stars - a.stars
        case 'updated':
          return (
            new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime()
          )
        case 'name':
          return a.name.localeCompare(b.name)
        default:
          return 0
      }
    })


  const [saveRepositories] = useSaveSelectedRepositoriesMutation()

  useEffect(() => {
    if (data?.repositories) {
      setRepositories(data.repositories)
    }
  }, [data])

  useEffect(() => {
    if (!error) {
      return
    }

    if (error.message.includes('token')) {
      return
    }

    toast.error("Could not fetch repositories", {
      description: error.message,
    })

  }, [error])

  const onDialogConfirm = async (selectedRepositoryIds: number[]) => {
    const { data } = await saveRepositories({
      variables: {
        selectedGithubRepositoryIds: selectedRepositoryIds
      }
    })

    if (data) {
      setRepositories(data.saveSelectedRepositories)
    }
  }

  return (
    <div>
      {!jwt && (
        <div className="text-center py-12">
          <LogIn className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            You're not logged in
          </h3>
          <p className="text-gray-500">
            Log in with your GitHub account to view your repositories
          </p>
        </div>
      )}

      {jwt && (
        <div className="mb-6">
          <div className="flex justify-between items-center mb-8">
            <div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Repositories
              </h2>
              <p className="text-gray-600">
                Manage and monitor your imported repositories
              </p>
            </div>
            {!loading && (
              <ImportReposDialog
                onConfirm={onDialogConfirm}
                alreadyTracked={repositories.map((e) => e.githubId)}
              />
            )}
          </div>
        </div>
      )}

      {jwt && loading && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 z-50">
          <SwirlingEffectSpinner />
        </div>
      )}

      {jwt && !loading && repositories.length === 0 && (
        <div className="text-center py-12">
          <Package className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No repositories found
          </h3>
          <p className="text-gray-500">
            Import your GitHub repositories to start monitoring dependency
            health
          </p>
        </div>
      )}

      {jwt && !loading && repositories.length > 0 && (
        <div>
          <div className="flex flex-col sm:flex-row gap-4 mb-8">
            <div className="flex-1">
              <Input
                placeholder="Search repositories..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full"
              />
            </div>
            <div className="flex gap-2">
              <Button
                variant={sortBy === 'health' ? 'default' : 'outline'}
                onClick={() => setSortBy('health')}
                size="sm"
              >
                Health Score
              </Button>
              <Button
                variant={sortBy === 'stars' ? 'default' : 'outline'}
                onClick={() => setSortBy('stars')}
                size="sm"
              >
                Stars
              </Button>
              <Button
                variant={sortBy === 'updated' ? 'default' : 'outline'}
                onClick={() => setSortBy('updated')}
                size="sm"
              >
                Updated
              </Button>
              <Button
                variant={sortBy === 'name' ? 'default' : 'outline'}
                onClick={() => setSortBy('name')}
                size="sm"
              >
                Name
              </Button>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {filteredAndSortedRepositories.map((repository, index) => (
              <RepositoryOverview key={index} repository={repository} />
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
