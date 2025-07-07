import { type DependencySortField, type SortDirection, useGetDependenciesQuery, type DependencyType } from "@/generated/graphql";
import { useEffect, useState } from "react"
import { DataTable } from "./components/data-table";
import { columns } from "./components/columns";
import SwirlingEffectSpinner from "@/components/customized/spinner/spinner-06";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/auth/AuthContext";
import { LogIn } from "lucide-react";

const sortColumnMap: Record<string, DependencySortField> = {
  name: "NAME",
  ecosystem: "ECOSYSTEM",
  scanStatus: "SCAN_STATUS",
  scannedAt: 'SCANNED_AT',
  errorMessage: "ERROR_MESSAGE"
}
const sortDirectionMap: Record<string, SortDirection> = {
  asc: "ASC",
  desc: "DESC"
}

export const Dependencies = () => {
  const { jwt } = useAuth()
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

  useEffect(() => {
    if (!error) {
      return
    }

    toast.error("Could not fetch dependencies", {
      description: error.message,
    })

  }, [error])

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
      {!jwt && (
        <div className="text-center py-12">
          <LogIn className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            You're not logged in
          </h3>
          <p className="text-gray-500">
            Log in with your GitHub account to view your repositories
          </p>
        </div>
      )}

      {loading && jwt && (
        <div className="fixed inset-0 flex items-center justify-center bg-white/80 z-50">
          <SwirlingEffectSpinner />
        </div>
      )}

      { jwt && (
        <DataTable columns={columns} data={posts} handleSearch={handleSearch} handleStatusFilter={handleStatusFilter} statusTotals={totals} handleSort={handleSort} handlePageSize={handlePageSize} handleSetPage={handleSetPage} pageSize={pageSize} totalPages={totalPages} currentPage={page} />
      )}
    </div>
  )
}
