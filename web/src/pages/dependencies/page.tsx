import { type DependencySortField, type SortDirection, useGetDependenciesQuery, type DependencyType } from "@/generated/graphql";
import { useEffect, useState } from "react"
import { DataTable } from "./components/data-table";
import { columns } from "./components/columns";
import SwirlingEffectSpinner from "@/components/customized/spinner/spinner-06";

const sortColumnMap: Record<string, DependencySortField> = {
  name: "NAME",
  ecosystem: "ECOSYSTEM",
  status: "STATUS",
  repositoryUrlCheckedAt: 'CHECKED_AT',
  repositoryUrlResolveFailedReason: "FAILED_REASON"
}
const sortDirectionMap: Record<string, SortDirection> = {
  asc: "ASC",
  desc: "DESC"
}

export const Dependencies = () => {
  const [search, setSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('');
  const [statuses, setStatuses] = useState<string[]>([])
  const [sortColumn, setSortColumn] = useState<DependencySortField>("NAME")
  const [sortDirection, setSortDirection] = useState<SortDirection>("ASC")
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
        direction: sortDirection
      }
    },
    notifyOnNetworkStatusChange: true,
  })

  const posts: DependencyType[] = (data?.dependencies.dependencies ?? [])
  const totalPages = data?.dependencies?.totalPages ?? 0

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setDebouncedSearch(search);
    }, 500);

    return () => clearTimeout(timeoutId);
  }, [search]);

  useEffect(() => {
    if (data && totals.completed === undefined) {
      setTotals({
        completed: data.dependencies.completed,
        pending: data.dependencies.pending,
        failed: data.dependencies.failed,
      })
    }
  }, [data, totals.completed])

  const handleSearch = (searchValue: string) => {
    setSearch(searchValue)
  }

  const handleStatusFilter = (selectedValues: string[]) => {
    setStatuses(selectedValues)
  }

  const handleSort = (sortColumn: string, sortDirection: string) => {
    setSortColumn(sortColumnMap[sortColumn])
    setSortDirection(sortDirectionMap[sortDirection])
  }

  const handlePageSize = (size: number) => {
    setPageSize(size)
  }

  const handleSetPage = (page: number) => {
    setPage(page)
  }

  return (
    <div className="container mx-auto py-10">
      {loading && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 z-50">
          <SwirlingEffectSpinner />
        </div>
      )}

      {!loading && (
        <DataTable columns={columns} data={posts} handleSearch={handleSearch} handleStatusFilter={handleStatusFilter} statusTotals={totals} handleSort={handleSort} handlePageSize={handlePageSize} handleSetPage={handleSetPage} pageSize={pageSize} totalPages={totalPages} currentPage={page} />
      )}
    </div>
  )
}
