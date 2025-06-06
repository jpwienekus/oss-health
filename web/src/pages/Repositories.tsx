import type { GitHubRepository } from '@/types'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Activity, Calendar, Eye, GitFork, Package, Star } from 'lucide-react'
import { Input } from '@/components/ui/input'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/auth/AuthContext'
import { ImportReposDialog } from '@/components/repositories/ImportReposDialog'
import { getClient } from '@/graphql/client'
import { SAVE_SELECTED_REPOSITORIES } from '@/graphql/mutations'
import { GET_REPOSITORIES } from '@/graphql/queries'

export const Repositories = () => {
  const { jwt } = useAuth()
  const [data, setData] = useState<GitHubRepository[]>([])

  const [searchTerm, setSearchTerm] = useState('')
  const [sortBy, setSortBy] = useState<'health' | 'stars' | 'updated' | 'name'>(
    'health',
  )

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

      const client = getClient(jwt)
      const response = await client.request<{
        repositories: GitHubRepository[]
      }>(GET_REPOSITORIES)

      setData(response.repositories)
    }
    fetchRepositories()
  }, [jwt])

  const formatDate = (dateString: string | null) => {
    if (!dateString) {
      return '-'
    }

    const date = new Date(dateString)
    if (!date) {
      return '-'
    }

    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    })
  }

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
    console.log(33, response)
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
            <ImportReposDialog onConfirm={onDialogConfirm} />
          </div>
        </div>
      )}
      {jwt && data.length === 0 && (
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

      {jwt && data.length > 0 && (
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
              <Card key={index} className="hover:shadow-lg transition shadow">
                <CardHeader className="pb-4">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <CardTitle className="text-lg flex items-center gap-2">
                        <span>{repository.name}</span>
                        {repository.private && (
                          <Badge variant="secondary" className="text-xs">
                            Private
                          </Badge>
                        )}
                      </CardTitle>
                      <CardDescription className="mt-1">
                        {repository.description ?? '-'}
                      </CardDescription>
                    </div>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  <div className="space-y-3">
                    <div className="flex justify-between items-center text-sm">
                      <span className="flex items-center gap-1">
                        <Activity className="w-4 h-4" />
                        Maintenance
                      </span>
                    </div>

                    <div className="flex justify-between items-center text-sm">
                      <span className="flex items-center gap-1">
                        <Activity className="w-4 h-4" />
                        Community
                      </span>
                    </div>

                    <div className="flex justify-between items-center text-sm">
                      <span className="flex items-center gap-1">
                        <Activity className="w-4 h-4" />
                        Release Cadence
                      </span>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-4 pt-2 text-sm text-gray-600">
                    <div className="flex items-center gap-1">
                      <Star className="h-4 w-4" />
                      <span>{repository.stars}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <GitFork className="h-4 w-4" />
                      <span>{repository.forks}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Eye className="h-4 w-4" />
                      <span>{repository.watchers}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      <span>{formatDate(repository.updatedAt)}</span>
                    </div>
                  </div>

                  <div className="flex justify-between items-center">
                    <Badge className="">{333} vulnerabilities</Badge>
                    <Button variant="outline" size="sm" asChild>
                      <a href="" target="_blank" rel="noopener noreferrer">
                        View on GitHub
                      </a>
                    </Button>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
