import { useAuth } from "@/auth/AuthContext"
import { Button } from "@/components/ui/button"
import { type DependencySortField, type SortDirection, useGetDependenciesQuery, type DependencyType } from "@/generated/graphql";
import { useEffect, useRef, useState } from "react"
import { DataTable } from "./components/data-table";
import { columns } from "./components/columns";
import { statuses } from "./components/columns"

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
  const [pageSize, setPageSize] = useState(20)

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

  const handleNext = () => {
    setPage(page < totalPages ? page + 1 : totalPages)
  };

  const handlePrevious = () => {
    setPage(page > 1 ? page - 1 : 1)
  };

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={posts} handleSearch={handleSearch} handleStatusFilter={handleStatusFilter} statusTotals={totals} handleSort={handleSort}/>

      <div className="flex justify-between space-x-2">
        <Button
          onClick={handlePrevious}
          disabled={page == 1 || loading}
          variant="outline"
        >
          Previous
        </Button>

        <Button
          onClick={handleNext}
          disabled={page == totalPages || loading}
          variant="outline"
        >
          Next
        </Button>
      </div>
    </div>
  )
}
