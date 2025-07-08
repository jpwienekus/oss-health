import {
  type DependencySortField,
  type SortDirection,
  useGetDependenciesQuery,
  type DependencyType,
} from '@/generated/graphql'
import { useEffect, useState } from 'react'
import { DataTable } from '@/components/admin/data-table'
import { columns } from '@/components/admin/columns'
import SwirlingEffectSpinner from '@/components/customized/spinner/spinner-06'
import { toast } from 'sonner'
import { useAuth } from '@/auth/AuthContext'
import { RequestLogin } from '@/components/request-login'

const sortColumnMap: Record<string, DependencySortField> = {
  name: 'NAME',
  ecosystem: 'ECOSYSTEM',
  scanStatus: 'SCAN_STATUS',
  scannedAt: 'SCANNED_AT',
  errorMessage: 'ERROR_MESSAGE',
}
const sortDirectionMap: Record<string, SortDirection> = {
  asc: 'ASC',
  desc: 'DESC',
}

export const Dependencies = () => {
  const { jwt } = useAuth()
  const [search, setSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('')
  const [statuses, setStatuses] = useState<string[]>([])
  const [sortColumn, setSortColumn] = useState<DependencySortField>('NAME')
  const [sortDirection, setSortDirection] = useState<SortDirection>('ASC')
  const [totals, setTotals] = useState<Record<string, number | undefined>>({
    completed: undefined,
    pending: undefined,
    failed: undefined,
  })

  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const { data, loading, error } = useGetDependenciesQuery({
    variables: {
      pagination: {
        page: page,
        pageSize: pageSize,
      },
      filter: {
        name: debouncedSearch,
        statuses: statuses,
      },
      sort: {
        field: sortColumn,
        direction: sortDirection,
      },
    },
    notifyOnNetworkStatusChange: true,
  })

  const dependencies: DependencyType[] = data?.dependencies.dependencies ?? []
  const totalPages = data?.dependencies?.totalPages ?? 0

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setDebouncedSearch(search)
    }, 500)

    return () => clearTimeout(timeoutId)
  }, [search])

  useEffect(() => {
    if (data && totals.completed === undefined) {
      setTotals({
        completed: data.dependencies.completed,
        pending: data.dependencies.pending,
        failed: data.dependencies.failed,
      })
    }
  }, [data, totals.completed])

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error('Could not fetch dependencies', {
      description: error.message,
    })
  }, [error])

  const handleSearch = (searchValue: string) => {
    setSearch(searchValue)
  }

  const handleStatusFilter = (selectedValues: string[]) => {
    setStatuses(selectedValues)
  }

  const handleSort = (sortColumnNew: string, sortDirectionNew: string) => {
    setSortColumn(sortColumnMap[sortColumnNew])
    setSortDirection(sortDirectionMap[sortDirectionNew])
  }

  const handlePageSize = (size: number) => {
    setPageSize(size)
  }

  const handleSetPage = (pageNew: number) => {
    setPage(pageNew)
  }

  return (
    <div>
      {!jwt && <RequestLogin />}

      {loading && jwt && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 dark:bg-slate-900/80 z-50">
          <SwirlingEffectSpinner />
        </div>
      )}

      {jwt && (
        <DataTable
          columns={columns}
          data={dependencies}
          handleSearch={handleSearch}
          handleStatusFilter={handleStatusFilter}
          statusTotals={totals}
          handleSort={handleSort}
          handlePageSize={handlePageSize}
          handleSetPage={handleSetPage}
          pageSize={pageSize}
          totalPages={totalPages}
          currentPage={page}
        />
      )}
    </div>
  )
}
