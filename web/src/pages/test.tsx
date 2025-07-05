
import { useAuth } from "@/auth/AuthContext"
import { Button } from "@/components/ui/button"
import { DependencySortField, SortDirection, useGetDependenciesQuery, type DependencyEdge, type PageInfo } from "@/generated/graphql"
import { useEffect, useState } from "react"

export const Dependencies = () => {
  const { jwt } = useAuth()
  const [afterCursor, setAfterCursor] = useState<number | null>(null)
  const [beforeCursor, setBeforeCursor] = useState<number | null>(null)
  const [search, setSearch] = useState('')
  const [dependencies, setDependencies] = useState<DependencyEdge[]>([])
  const [pageInfo, setPageInfo] = useState<PageInfo>({ hasNextPage: true, hasPreviousPage: false })
  // const [loading, setLoading] = useState(false)
  //
  const [direction, setDirection] = useState<'next' | 'previous'>('next')
  const [cursor, setCursor] = useState<number | null>(null)


  const { data, loading, error } = useGetDependenciesQuery({
    variables: {
      pagination: {
        limit: 100,
        after: direction === 'next' ? cursor : null,
      },
      filter: {
        name: search
      },
      sort: {
        field: DependencySortField.Id,
        direction: SortDirection.Asc
      }
    },
    skip: !jwt,
    notifyOnNetworkStatusChange: true
  })

  useEffect(() => {
    if (!data?.dependencies?.edges) {
      return
    }

    setDependencies(data.dependencies.edges)
    setPageInfo(data.dependencies.pageInfo)
    setAfterCursor(data.dependencies.pageInfo.endCursor ?? null)
    setBeforeCursor(data.dependencies.pageInfo.startCursor ?? null)
  }, [data])


  const handleNext = () => {
    if (!pageInfo.hasNextPage) {
      return
    }

    setCursor(afterCursor)
    setDirection('next')
  }

  const handlePrevious = () => {
    if (!pageInfo.hasPreviousPage) {
      return
    }


    setCursor(beforeCursor)
    setDirection('previous')
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
