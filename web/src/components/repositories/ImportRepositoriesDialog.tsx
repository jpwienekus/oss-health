import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useEffect, useState } from 'react'
import { Calendar, Eye, GitFork, Import, Search, Star } from 'lucide-react'
import { Badge } from '../ui/badge'
import { ScrollArea } from '../ui/scroll-area'
import { Checkbox } from '../ui/checkbox'
import { formatDate } from '@/utils'
import {
  useGithubRepositoriesLazyQuery,
  type GitHubRepository,
} from '@/generated/graphql'
import { toast } from 'sonner'

type ImportReposDialogParams = {
  alreadyTracked: number[]
  onConfirm: (selectedRepoIds: number[]) => void
}

export const ImportReposDialog = ({
  alreadyTracked,
  onConfirm,
}: ImportReposDialogParams) => {
  // TODO: Fetch with account info
  const MAX_REPOSITORIES = 2
  const [isOpen, setIsOpen] = useState<boolean>(false)
  const [searchTerm, setSearchTerm] = useState<string>('')

  const [selectedRepositories, setSelectedRepositories] = useState<number[]>([])
  const [getRepositories, { data, error, loading }] =
    useGithubRepositoriesLazyQuery()

  const filteredRepos = (data?.githubRepositories ?? [])
    .filter(
      (repo) =>
        (repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          repo.description?.toLowerCase().includes(searchTerm.toLowerCase())) &&
        !alreadyTracked.includes(repo.githubId),
    )
    .sort(
      (a: GitHubRepository, b: GitHubRepository) =>
        new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime(),
    )

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error('Could not fetch repositories', {
      description: error.message,
    })
  }, [error])

  const fetchGitHubRepos = async () => {
    await getRepositories()
    setSelectedRepositories([])
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
      {alreadyTracked.length < MAX_REPOSITORIES && (
        <Button onClick={fetchGitHubRepos} disabled={loading}>
          <Import className="w-4 h-4" />
          <span>{loading ? 'Loading...' : 'Import from GitHub'}</span>
        </Button>
      )}
      {alreadyTracked.length >= MAX_REPOSITORIES && (
        <div className="text-sm text-slate-500 dark:text-slate-400 p-2 rounded-md">
          Maximum number of repositories ({MAX_REPOSITORIES}) reached
        </div>
      )}
      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent
          className="max-w-2xl max-h-[80vh] overflow-hidden flex flex-col"
          aria-describedby="import-dialog-description"
        >
          <DialogHeader>
            <DialogTitle>Select Reporistories to import</DialogTitle>
          </DialogHeader>

          <DialogDescription></DialogDescription>

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
                  className={`flex items-start gap-3 p-4 rounded-lg border transition-colors cursor-pointer hover:bg-slate-50 dark:hover:bg-slate-800 ${
                    selectedRepositories.includes(repo.githubId)
                      ? 'bg-blue-50 border-blue-200 dark:bg-blue-900 dark:border-blue-700'
                      : ''
                  }}`}
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

                    <p className="text-sm text-slate-500 dark:text-slate-400 mb-2 line-clamp-2">
                      {repo.description ?? '-'}
                    </p>

                    <div className="flex flex-wrap gap-4 text-xs text-slate-500 dark:text-slate-400">
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
