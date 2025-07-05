import { useAuth } from "@/auth/AuthContext"
import { Button } from "@/components/ui/button"
import { DependencySortField, SortDirection, useGetDependenciesQuery, type DependencyType } from "@/generated/graphql";
import { useEffect, useRef, useState } from "react"
import { DataTable } from "./components/data-table";
import { columns } from "./components/columns";

export const Dependencies = () => {
  const [search, setSearch] = useState('')
  const [debouncedSearch, setDebouncedSearch] = useState('');

  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)

  const { data, loading, error } = useGetDependenciesQuery({
    variables: {
      pagination: {
        page: page,
        pageSize: pageSize,
      },
      filter: {
        name: debouncedSearch
      },
      sort: {
        field: DependencySortField.Name,
        direction: SortDirection.Asc
      }
    },
    notifyOnNetworkStatusChange: true,
  })

  const posts: DependencyType[] = (data?.dependencies.dependencies ?? [])
  const totalPages = data?.dependencies?.total ?? 0


  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setDebouncedSearch(search);
    }, 500);

    return () => clearTimeout(timeoutId);
  }, [search]);

  const handleSearch = (searchValue: string) => {
    setSearch(searchValue)
  }

  const handleNext = () => {
    setPage(page < totalPages ? page + 1 : totalPages)
  };

  const handlePrevious = () => {
    setPage(page > 1 ? page - 1 : 1)
  };

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={posts} handleSearch={handleSearch} />

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
