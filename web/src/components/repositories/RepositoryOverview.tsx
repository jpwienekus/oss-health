import type { GitHubRepository } from '@/types'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  AlertTriangle,
  CheckCircle,
  Package,
  XCircle,
  RefreshCw,
  CalendarClock,
} from 'lucide-react'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/auth/AuthContext'
import { getClient } from '@/graphql/client'
import { MANUAL_SCAN_DEBUG } from '@/graphql/queries'
import { formatDate } from '@/utils'

type RepositoryOverviewParams = {
  repository: GitHubRepository
  onUpdate: (result: GitHubRepository[]) => void
}

export const RepositoryOverview = ({
  repository,
  onUpdate,
}: RepositoryOverviewParams) => {
  const { jwt } = useAuth()
  const [loading, setLoading] = useState(false)

  const manualScanForDebug = async (repositoryId: number) => {
    if (!jwt) {
      return
    }
    const client = getClient(jwt)
    setLoading(true)
    const response = await client.request<{
      manualScanDebug: GitHubRepository[]
    }>(MANUAL_SCAN_DEBUG, {
      repositoryId,
    })
    setLoading(false)
    onUpdate(response.manualScanDebug)
  }

  const getHealthIcon = (score: number) => {
    if (score >= 80) {
      return <CheckCircle className="w-4 h-4 text-green-600" />
    } else if (score >= 60) {
      return <AlertTriangle className="w-4 h-4 text-yello-600" />
    } else {
      return <XCircle className="w-4 h-4 text-red-600" />
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
    <Card className="hover:shadow-lg transition shadow">
      <CardHeader className="pb-4">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <CardTitle className="text-lg flex items-center gap-2">
              <span className="text-md font-medium">{repository.name}</span>
              {repository.scannedDate !== null &&
              repository.scannedDate !== undefined ? (
                <>
                  {getHealthIcon(repository.score)}
                  <span
                    className={`text-xs font-medium ${getHealthColor(repository.score)}`}
                  >
                    {repository.score}/100
                  </span>
                </>
              ) : (
                <span className="text-xs font-medium text-muted-foreground">
                  Not scanned yet
                </span>
              )}
              <Button
                variant="ghost"
                size="icon"
                onClick={() => manualScanForDebug(repository.id)}
              >
                {loading ? (
                  <RefreshCw className="w-4 h-4 animate-spin" />
                ) : (
                  <RefreshCw className="w-2 h-2" />
                )}
              </Button>
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
            <span>
              {repository.scannedDate ? repository.dependencies : '-'}{' '}
              dependencies
            </span>
          </div>
          <div className="flex items-center gap-1">
            <AlertTriangle size={12} />
            <span
              className={
                repository.scannedDate && repository.vulnerabilities > 0
                  ? 'text-red-600'
                  : ''
              }
            >
              {repository.scannedDate ? repository.vulnerabilities : '-'}{' '}
              vulnerabilities
            </span>
          </div>
          <div className="flex items-center gap-1">
            <CalendarClock size={12} />
            <span>
              {repository.scannedDate
                ? formatDate(repository.scannedDate)
                : '-'}
            </span>
          </div>
        </div>
      </CardHeader>
    </Card>
  )
}
