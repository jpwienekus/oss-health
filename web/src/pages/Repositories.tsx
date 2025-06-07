import type { GitHubRepository } from '@/types'
import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  AlertTriangle,
  CheckCircle,
  Package,
  Star,
  XCircle,
} from 'lucide-react'
import { Input } from '@/components/ui/input'
import {
  Card,
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

  const getHealthIcon = (score: number) => {
    if (score >= 80) {
      return <CheckCircle className="w-5 h-5 text-green-600" />
    } else if (score >= 60) {
      return <AlertTriangle className="w-5 h-5 text-yello-600" />
    } else {
      return <XCircle className="w-5 h-5 text-red-600" />
    }
  }

  const getHealthColor = (score: number) => {
    if (score >= 80) {
      return 'text-green-600'
    } else if (score >= 60) {
      return 'text-yellow-600'
    } else {
      return 'text-red-600'
    }
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
            <ImportReposDialog
              onConfirm={onDialogConfirm}
              alreadyTracked={data.map((e) => e.githubId)}
            />
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
                        {getHealthIcon(repository.score)}
                        <span
                          className={`text-sm font-medium ${getHealthColor(repository.score)}`}
                        >
                          {repository.score}/100
                        </span>
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

                  <div className="flex flex-wrap gap-4 text-xs text-slate-500">
                    <div className="flex items-center gap-1">
                      <Package size={12} />
                      <span>{repository.dependencies} dependencies</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <AlertTriangle size={12} />
                      <span
                        className={
                          repository.vulnerabilities > 0 ? 'text-red-600' : ''
                        }
                      >
                        {repository.vulnerabilities} vulnerabilities
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Star size={12} />
                      <span>{repository.stars}</span>
                    </div>
                  </div>
                </CardHeader>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
