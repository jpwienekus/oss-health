import {
  AlertTriangle,
  CheckCircle,
  Package,
  XCircle,
  CalendarClock,
} from 'lucide-react'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { formatDate } from '@/utils'
import type { GitHubRepository } from '@/generated/graphql'

type RepositoryOverviewParams = {
  repository: GitHubRepository
}

export const RepositoryOverview = ({
  repository,
}: RepositoryOverviewParams) => {
  const getHealthIcon = (score: number) => {
    if (score >= 80) {
      return <CheckCircle className="w-4 h-4 text-green-600" />
    } else if (score >= 60) {
      return <AlertTriangle className="w-4 h-4 text-yellow-600" />
    } else {
      return <XCircle className="w-4 h-4 text-red-600" />
    }
  }

  const getHealthColor = (score: number | null) => {
    if (!score) {
      return ''
    } else if (score >= 80) {
      return 'text-green-600'
    } else if (score >= 60) {
      return 'text-yellow-600'
    } else {
      return 'text-red-600'
    }
  }

  return (
    <Card className="hover:shadow-lg transition shadow">
      <CardHeader className="pb-4">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <CardTitle className="text-lg flex items-center gap-2">
              <span className="text-md font-medium">{repository.name}</span>
              {repository.lastScannedAt !== null &&
              repository.lastScannedAt !== undefined ? (
                <>
                  {getHealthIcon(repository.score ?? 0)}
                  <span
                    className={`text-xs font-medium ${getHealthColor(repository.score ?? null)}`}
                  >
                    {repository.score}/100
                  </span>
                </>
              ) : (
                <span className="text-xs font-medium text-muted-foreground">
                  Not scanned yet
                </span>
              )}
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

        <div className="flex flex-wrap gap-4 text-xs text-slate-500 dark:text-slate-400">
          <div className="flex items-center gap-1">
            <Package size={12} />
            <span>
              {repository.lastScannedAt ? repository.dependencies : '-'}{' '}
              dependencies
            </span>
          </div>
          <div className="flex items-center gap-1">
            <AlertTriangle size={12} />
            <span
              className={
                repository.lastScannedAt &&
                repository.vulnerabilities &&
                repository.vulnerabilities > 0
                  ? 'text-red-600'
                  : ''
              }
            >
              {repository.lastScannedAt ? repository.vulnerabilities : '-'}{' '}
              vulnerabilities
            </span>
          </div>
          <div className="flex items-center gap-1">
            <CalendarClock size={12} />
            <span>
              {repository.lastScannedAt
                ? formatDate(repository.lastScannedAt)
                : '-'}
            </span>
          </div>
        </div>
      </CardHeader>
    </Card>
  )
}
