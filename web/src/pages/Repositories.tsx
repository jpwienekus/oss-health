import type { GitHubRepository } from '@/types'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Package } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/auth/AuthContext'
import { ImportReposDialog } from '@/components/repositories/ImportReposDialog'
import { getClient } from '@/graphql/client'
import { SAVE_SELECTED_REPOSITORIES } from '@/graphql/mutations'
import { GET_REPOSITORIES } from '@/graphql/queries'
import { RepositoryOverview } from '@/components/repositories/RepositoryOverview'

export const Repositories = () => {
  const { jwt } = useAuth()
  const [data, setData] = useState<GitHubRepository[]>([])

  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState<'health' | 'stars' | 'updated' | 'name'>(
    'health',
  )
  const [isLoading, setIsLoading] = useState(false)

  const filteredAndSortedRepositories = data
    .filter(
      (repo) =>
        repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        repo.description?.toLowerCase().includes(searchTerm.toLowerCase()),
    )
    .sort((a: GitHubRepository, b: GitHubRepository) => {
      switch (sortBy) {
        case 'health':
          return b.score - a.score
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

  useEffect(() => {
    const fetchRepositories = async () => {
      if (!jwt) {
        return
      }

      setIsLoading(true)
      const client = getClient(jwt)
      const response = await client.request<{
        repositories: GitHubRepository[]
      }>(GET_REPOSITORIES)
      setData(response.repositories)
      setIsLoading(false)
    }
    fetchRepositories()
  }, [jwt])

  const onDialogConfirm = async (selectedRepositoryIds: number[]) => {
    if (!jwt) {
      return
    }

    const client = getClient(jwt)
    const response = await client.request<{
      saveSelectedRepositories: GitHubRepository[]
    }>(SAVE_SELECTED_REPOSITORIES, {
      selectedGithubRepositoryIds: selectedRepositoryIds,
    })
    setData(response.saveSelectedRepositories)
  }

  return (
    <div>
      {!jwt && (
        <div className="text-center py-12">
          <Package className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No repositories found
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
            {!isLoading && (
              <ImportReposDialog
                onConfirm={onDialogConfirm}
                alreadyTracked={data.map((e) => e.githubId)}
              />
            )}
          </div>
        </div>
      )}

      {jwt && isLoading && (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto mb-4"></div>
          <p className="text-gray-500">Loading repositories...</p>
        </div>
      )}

      {jwt && !isLoading && data.length === 0 && (
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

      {jwt && !isLoading && data.length > 0 && (
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
