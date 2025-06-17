import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuth } from '@/auth/AuthContext'
import { useState } from 'react'
import type { GitHubRepository } from '@/types'
import { getClient } from '@/graphql/client'
import { Calendar, Eye, GitFork, Import, Search, Star } from 'lucide-react'
import { GET_REPOSITORIES_FROM_GITHUB } from '@/graphql/queries'
import { Badge } from '../ui/badge'
import { ScrollArea } from '../ui/scroll-area'
import { Checkbox } from '../ui/checkbox'
import {formatDate} from '@/utils'

type ImportReposDialogParams = {
  alreadyTracked: number[]
  onConfirm: (selectedRepoIds: number[]) => void
}

export const ImportReposDialog = ({
  alreadyTracked,
  onConfirm,
}: ImportReposDialogParams) => {
  const { jwt } = useAuth()
  const [isOpen, setIsOpen] = useState<boolean>(false)
  const [searchTerm, setSearchTerm] = useState<string>('')
  const [isImporting, setIsImporting] = useState<boolean>(false)
  const [githubRepositories, setGitHubRepositories] = useState<
    GitHubRepository[]
  >([])
  const [selectedRepositories, setSelectedRepositories] = useState<number[]>([])

  const filteredRepos = githubRepositories
    .filter(
      (repo) =>
        repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        repo.description?.toLowerCase().includes(searchTerm.toLowerCase()),
    )
    .sort(
      (a: GitHubRepository, b: GitHubRepository) =>
        new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
    )

  const fetchGitHubRepos = async () => {
    if (!jwt) {
      return
    }

    setIsImporting(true)
    const client = getClient(jwt)
    const response = await client.request<{
      githubRepositories: GitHubRepository[]
    }>(GET_REPOSITORIES_FROM_GITHUB)
    setGitHubRepositories(
      response.githubRepositories.filter(
        (e) => !alreadyTracked.includes(e.githubId),
      ),
    )
    setIsImporting(false)
    setIsOpen(true)
  }

  const handleRepositoryToggle = (githubId: number) => {
    setSelectedRepositories((previous) =>
      previous.includes(githubId)
        ? previous.filter((id) => id !== githubId)
        : [...previous, githubId],
    )
  }

  const handleConfirm = () => {
    onConfirm(selectedRepositories)
    setIsOpen(false)
  }

  return (
    <div>
      <Button onClick={fetchGitHubRepos} disabled={isImporting}>
        <Import className="w-4 h-4" />
        <span>{isImporting ? 'Loading...' : 'Import from GitHub'}</span>
      </Button>
      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-hidden flex flex-col">
          <DialogHeader>
            <DialogTitle>Select Reporistories to import</DialogTitle>
          </DialogHeader>

          <div className="relative mb-1">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
            <Input
              placeholder="Search repositories..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10"
            />
          </div>

          <ScrollArea className="pr-4 h-72">
            <div className="space-y-3 py-2">
              {filteredRepos.map((repo, index) => (
                <div
                  key={index}
                  className={`flex items-start gap-3 p-4 rounded-lg border transition-colors cursor-pointer hover:bg-slate-50 ${selectedRepositories.includes(repo.githubId) ? 'bg-blue-50 border-blue-200' : 'bg-white'}`}
                  onClick={() => handleRepositoryToggle(repo.githubId)}
                >
                  <Checkbox
                    checked={selectedRepositories.includes(repo.githubId)}
                    onClick={(e) => {
                      e.stopPropagation()
                      handleRepositoryToggle(repo.githubId)
                    }}
                    className="mt-1"
                  />

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-semibold text-sm">{repo.name}</span>
                      {repo.private && (
                        <Badge variant="outline" className="text-xs">
                          Private
                        </Badge>
                      )}
                    </div>

                    <p className="text-sm text-slate-600 mb-2 line-clamp-2">
                      {repo.description ?? '-'}
                    </p>

                    <div className="flex flex-wrap gap-4 text-xs text-slate-500">
                      <div className="flex items-center gap-1">
                        <Star size={12} />
                        <span>{repo.stars}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <GitFork size={12} />
                        <span>{repo.forks}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <Eye size={12} />
                        <span>{repo.watchers}</span>
                      </div>
                      <div className="flex items-center gap-1">
                        <Calendar size={12} />
                        <span>{formatDate(repo.updatedAt)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>

          <DialogFooter className="flex gap-2">
            <Button variant="outline" onClick={() => setIsOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={handleConfirm}
              disabled={selectedRepositories.length === 0}
            >
              Confirm
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
