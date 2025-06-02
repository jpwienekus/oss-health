import { useGitHubPopupLogin } from "@/auth/useGitHubPopupLogin"
import { RequestGithubLogin } from "@/components/repositories/RequestGithubLogin"
import type { ColumnDef } from "@tanstack/react-table"
import { getClient } from "@/graphql/client"
import { SYNC_REPOS } from "@/graphql/mutations"
import { GET_REPOS } from "@/graphql/queries"
import type { Repository } from "@/types"
import { DataTable } from "@/components/repositories/Datatable"
import { useEffect, useState } from "react"
import { Badge } from "@/components/ui/badge"
import { IconBrandGithub } from "@tabler/icons-react"

export const Repositories = () => {
  const { jwt } = useGitHubPopupLogin()
  const [data, setData] = useState<Repository[]>([])
  const [syncDate, setSyncDate] = useState<Date | undefined>(undefined)
  const columns: ColumnDef<Repository>[] = [
    {
      accessorKey: "name",
      header: "Name",
      cell: ({ row }) => <div className="w-[100px]">{row.getValue("name")}</div>,
    },
    {
      accessorKey: "description",
      header: "Description",
      cell: ({ row }) => {
        return (
          <div className="max-w-[500px] truncate">
            {row.getValue("description")}
          </div>
        )
      }
    },
    {
      accessorKey: "score",
      header: "Score",
      cell: ({ row }) => {
        const score: number = row.getValue("score")
        const variant = score < 7 ? 'red' : (score < 8 ? 'yellow' : 'green')

        return (
          <div>
            <Badge variant="secondary" className={`bg-${variant}-500 text-white dark:bg-${variant}-600`}>{score}</Badge>
          </div>
        )
      }
    },
    {
      accessorKey: "openIssues",
      header: "Open Issues",
      cell: ({ row }) => {
        return (
          <div>
            {row.getValue("openIssues")}
          </div>
        )
      }
    },
    {
      accessorKey: "updatedAt",
      header: "Updated At",
      cell: ({ row }) => {
        const dateFormatted = new Date(row.getValue('updatedAt')).toLocaleDateString('en-US', {
          year: 'numeric',
          month: 'short',
          day: 'numeric'
        })

        return <div>{dateFormatted}</div>
      },
    },
    {
      accessorKey: "url",
      header: "",
      cell: ({ row }) => {
        const url: string = row.getValue("url")
        return <div>
          <a href={url.replace("api.", "").replace("repos/", "")}
            target="_blank"
          >
            <IconBrandGithub size={15} />
          </a>
        </div>
      },
    },
  ]

  const fetchRepositories = async () => {
    if (!jwt) {
      return
    }

    const client = getClient(jwt)
    const response = await client.request<{ repositories: { repositories: Repository[], syncDate: string } }>(GET_REPOS)
    setData(response.repositories.repositories)
    setSyncDate(response.repositories.syncDate ? new Date(response.repositories.syncDate) : undefined)

  }

  const syncRepositories = async () => {
    if (!jwt) {
      return
    }

    const client = getClient(jwt)
    const response = await client.request<{ syncRepositories: { repositories: Repository[], syncDate: string } }>(SYNC_REPOS)
    setData(response.syncRepositories.repositories)
    setSyncDate(response.syncRepositories.syncDate ? new Date(response.syncRepositories.syncDate) : undefined)
  }

  useEffect(() => {
    fetchRepositories()
  }, []);

  return (
    <div className='flex min-h-svh flex-col items-center'>
      {!jwt ? (
        <RequestGithubLogin />
      ) : (
        <div className="container mx-auto">
          <DataTable columns={columns} data={data} syncDate={syncDate} sync={syncRepositories} />
        </div>
      )}
    </div>
  )
}
