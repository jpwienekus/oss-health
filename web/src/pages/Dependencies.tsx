import { useAuth } from "@/auth/AuthContext"
import { Button } from "@/components/ui/button"
import { getClient } from "@/graphql/client"
import { GET_DEPENDENCIES } from "@/graphql/queries"
import type { DependencyEdge, PageInfo } from "@/types"
import { useEffect, useState } from "react"

export const Dependencies = () => {
  const { jwt } = useAuth()
  const [afterCursor, setAfterCursor] = useState<number | null>(null)
  const [beforeCursor, setBeforeCursor] = useState<number | null>(null)
  const [search, setSearch] = useState('')
  const [dependencies, setDependencies] = useState<DependencyEdge[]>([])
  const [pageInfo, setPageInfo] = useState<PageInfo>({ hasNextPage: true, hasPreviousPage: false })
  const [loading, setLoading] = useState(false)

  const fetchDependencies = async (cursor: number | null = null, name: string = '', direction: 'next' | 'previous' = 'next') => {
    if (!jwt) {
      return
    }

    setLoading(false)

    try {
      // TODO: setup graphql generation
      const variables = {
        pagination: { limit: 1000, after: direction === 'next' ? cursor : null },
        filter: { name },
        sort: { field: 'ID', direction: 'ASC' }
      }

      const client = getClient(jwt)
      const data = await client.request(GET_DEPENDENCIES, variables)
      setDependencies(data.dependencies.edges)
      setAfterCursor(data.dependencies.pageInfo.endCursor)
      setBeforeCursor(data.dependencies.pageInfo.startCursor)
      setPageInfo(data.dependencies.pageInfo)
      console.log(33, data)
    } catch (error) {
      console.error("Error fetching dependencies:", error)
    }
  }

  // TODO: debounce
  const handleSearch = (value: string) => {
    fetchDependencies(null, value)
  }

  useEffect(() => {
    handleSearch(search);
  }, [jwt]);

  useEffect(() => {
    handleSearch(search)
  }, [search])

  const handleNext = () => {
    if (!pageInfo.hasNextPage) {
      return
    }

    fetchDependencies(afterCursor, search, 'next')
  }

  const handlePrevious = () => {
    if (!pageInfo.hasPreviousPage) {
      return
    }

    fetchDependencies(beforeCursor, search, 'previous')
  }

  return (
    <div>
      <span>hello world</span>

      <div className="flex justify-between space-x-2">
        <Button
          onClick={handlePrevious}
          disabled={!pageInfo.hasPreviousPage || loading}
          variant="outline"
        >
          Previous
        </Button>

        <Button
          onClick={handleNext}
          disabled={!pageInfo.hasNextPage || loading}
          variant="outline"
        >
          Next
        </Button>
      </div>
    </div>
  )
}
